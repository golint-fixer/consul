// Package consul implements a high-level wrapper Consul HTTP client for easy service discovery.
// Provides additional features, such as time based lookups and retry policy.
package consul

import consul "github.com/hashicorp/consul/api"

// Client is a wrapper around the Consul API.
type Client interface {
	Health(service string, tag string, queryOpts *consul.QueryOptions) ([]*consul.ServiceEntry, *consul.QueryMeta, error)
}

type client struct {
	consul *consul.Client
}

// NewClient returns an implementation of the Client interface expecting a fully
// setup Consul Client.
func NewClient(c *consul.Config) Client {
	cli, _ := consul.NewClient(&*c)
	return &client{consul: cli}
}

// Health returns the list of healthy entries for a given service filtered by tag.
func (c *client) Health(service, tag string, opts *consul.QueryOptions) ([]*consul.ServiceEntry, *consul.QueryMeta, error) {
	return c.consul.Health().Service(service, tag, true, opts)
}
