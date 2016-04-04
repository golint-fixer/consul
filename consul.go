// Package consul implements a high-level wrapper Consul HTTP client for easy service discovery.
// Provides additional features, such as time based lookups and retry policy.
package consul

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	consul "github.com/hashicorp/consul/api"
	"gopkg.in/vinxi/balancer.v0"
)

var (
	// DefaultBalancer stores the roundrobin balancer used by default.
	DefaultBalancer = balancer.DefaultBalancer

	// DefaultWaitTimeout defines the wait interval for servers become available.
	DefaultWaitInterval = 100 * time.Millisecond
)

var (
	// ErrDiscoveryTimeout is used in case that discovery timeout exceeded.
	ErrDiscoveryTimeout = errors.New("consul: cannot discover servers due to timeout")
)

// Consul is a wrapper around the Consul API.
type Consul struct {
	// Mutex provides a struct mutex to prevent data races.
	sync.Mutex

	// quit is used internally to open/close the Consul servers update goroutine.
	quit chan bool

	// nodes is used to cached server nodes URLs provided by Consul servers for the specific service.
	nodes []string

	// Config stores the Consul client vinxi config options used for discovery.
	Config *Config

	// Retrier stores the retry strategy to be used if Consul discovery process fails.
	Retrier Retrier

	// Balancer stores the balancer to be used to distribute traffic
	// load across multiple servers provided by Consul.
	Balancer balancer.Balancer
}

// NewClient returns an implementation of the Client interface expecting a fully
// setup Consul Client.
func New(config *Config) *Consul {
	c := &Consul{Config: config, Retrier: DefaultRetrier, Balancer: DefaultBalancer}
	c.Start() // starts the server discovery background job
	return c
}

// nextConsulServer returns the next available server based on the current iteration index.
func (c *Consul) nextConsulServer(index int) (*consul.Config, bool) {
	servers := c.Config.Instances
	if l := len(servers); index < l {
		return servers[index], index != (l - 1)
	}
	return servers[0], false
}

// UpdateNodes is used to update a list of server nodes for the current discovery service.
func (c *Consul) UpdateNodes() ([]string, error) {
	var retries int
	var entries []*consul.ServiceEntry

	err := NewRetrier(c.Retrier).Run(func() error {
		var err error
		config, more := c.nextConsulServer(retries)
		if !more {
			retries = 0
		}
		client := NewClient(config)

		entries, _, err = client.Health(c.Config.Service, c.Config.Tag, c.Config.QueryOptions)
		return err
	})

	return mapNodes(entries), err
}

// Stop stops the Consul servers update interval goroutine.
func (c *Consul) Stop() {
	close(c.quit)
}

// Start starts the Consul servers update interval goroutine.
func (c *Consul) Start() {
	go c.updateInterval(c.Config.RefreshTime)
}

// updateInterval recursively ask to Consul servers to update the list of available server nodes.
func (c *Consul) updateInterval(interval time.Duration) {
	for {
		select {
		case <-c.quit:
			return
		default:
			nodes, err := c.UpdateNodes()
			if err != nil || nodes == nil {
				// TODO: handle error
			}

			c.Lock()
			c.nodes = nodes
			c.Unlock()
			time.Sleep(interval)
		}
	}
}

// GetNodes returns a list of server nodes hostnames for the configured service.
func (c *Consul) GetNodes() ([]string, error) {
	// Wait until nodes are available.
	// TODO: consider using a custom channel or WaitGroup
	loops := 0
	for c.nodes == nil {
		if loops > 100 { // Stop after 5 seconds. This should be configurable
			return nil, ErrDiscoveryTimeout
		}
		loops++
		time.Sleep(DefaultWaitInterval)
	}

	c.Lock()
	defer c.Unlock()
	return c.nodes, nil
}

func (c *Consul) getTargetHost(nodes []string) (string, error) {
	if c.Balancer == nil {
		return nodes[0], nil
	}
	return c.Balancer.Balance(nodes)
}

// HandleHTTP returns the list of healthy entries for a given service filtered by tag.
func (c *Consul) HandleHTTP(w http.ResponseWriter, r *http.Request, h http.Handler) {
	if len(c.Config.Instances) == 0 {
		h.ServeHTTP(w, r)
		return
	}

	// Retrieve latest service server nodes
	nodes, err := c.GetNodes()
	if err != nil || len(nodes) == 0 {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(err.Error()))
		return
	}

	// Balance traffic using the configured balancer
	target, err := c.getTargetHost(nodes)
	if err != nil {
		h.ServeHTTP(w, r)
		return
	}

	// Define the URL to forward the request
	r.Host = target
	r.URL.Host = target

	h.ServeHTTP(w, r)
}

// mapNodes maps the Consul specific service entry into a string hostname.
func mapNodes(entries []*consul.ServiceEntry) []string {
	instances := make([]string, len(entries))

	for i, entry := range entries {
		addr := entry.Node.Address

		if entry.Service.Address != "" {
			addr = entry.Service.Address
		}

		instances[i] = fmt.Sprintf("%s:%d", addr, entry.Service.Port)
	}

	return instances
}
