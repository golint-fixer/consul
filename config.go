package consul

import (
	"net/http"
	"net/url"
	"time"

	consul "github.com/hashicorp/consul/api"
)

var (
	// DefaultRefreshTime defines the default refresh inverval to be used.
	DefaultRefreshTime = 5 * time.Minute
)

// Config is used to configure Consul clients and servers.
type Config struct {
	// Service stores the Consul service name used for discovery.
	Service string

	// Tag stores the optional Consul service tag.
	Tag string

	// QueryOptions stores the Consul client addition query options.
	QueryOptions *consul.QueryOptions

	// Instances stores the consul.Config objects per Consul server.
	Instances []*consul.Config

	// Datacenter to use. If not provided, the default agent datacenter is used.
	Datacenter string

	// Token is used to provide a per-request ACL token
	// which overrides the agent's default token.
	Token string

	// HttpClient is the client to use. Default will be
	// used if not provided.
	HttpClient *http.Client

	// HttpAuth is the auth info to use for http access.
	HttpAuth *consul.HttpBasicAuth

	// WaitTime limits how long a Watch will block. If not provided,
	// the agent default values will be used.
	WaitTime time.Duration

	// RefreshTime defines the Consul server refresh how long a Watch will block. If not provided,
	// the agent default values will be used.
	RefreshTime time.Duration
}

// NewConfig creates a new Consul client config preconfigured for
// the given Consul service and a list of Consul servers.
func NewConfig(service string, servers ...string) *Config {
	c := &Config{Service: service, RefreshTime: DefaultRefreshTime}
	c.SetServer(servers...)
	return c
}

// SetServer sets one or multiple Consul servers, creating the default config
func (c *Config) SetServer(server ...string) error {
	servers := []*consul.Config{}
	for _, uri := range server {
		u, err := url.Parse(uri)
		if err != nil {
			return err
		}
		servers = append(servers, c.newConsulConfig(u))
	}
	c.Instances = servers
	return nil
}

// newConsulConfig creates a new Consul config for the given server URL
// based on default, env or parent params.
func (c *Config) newConsulConfig(u *url.URL) *consul.Config {
	config := consul.DefaultConfig()
	config.Address = u.Host
	config.Scheme = u.Scheme

	// Apply defaults fields based on parent config
	if c.Datacenter != "" {
		config.Datacenter = c.Datacenter
	}
	if c.HttpClient != nil {
		config.HttpClient = c.HttpClient
	}
	if c.HttpAuth != nil {
		config.HttpAuth = c.HttpAuth
	}
	if c.WaitTime != 0 {
		config.WaitTime = c.WaitTime
	}
	if c.Token != "" {
		config.Token = c.Token
	}

	return config
}
