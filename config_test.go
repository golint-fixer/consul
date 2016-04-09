package consul

import (
	"testing"

	consul "github.com/hashicorp/consul/api"
	"github.com/nbio/st"
)

func TestConfig(t *testing.T) {
	config := NewConfig("web", "http://foo.com", "http://bar.com")
	config.Token = "foo"
	config.Datacenter = "foo"

	st.Expect(t, config.Service, "web")
	st.Expect(t, config.RefreshTime, DefaultRefreshTime)
	st.Expect(t, config.Tag, "")
	st.Expect(t, config.Token, "foo")
	st.Expect(t, config.Datacenter, "foo")
	st.Expect(t, len(config.Instances), 2)
	st.Expect(t, config.Instances[0].Address, "foo.com")
	st.Expect(t, config.Instances[1].Address, "bar.com")

	// Check Consul defaults in child instances
	foo := config.Instances[0]
	st.Expect(t, foo.Token, "")
	st.Expect(t, foo.Datacenter, "")

	bar := config.Instances[0]
	st.Expect(t, bar.Token, "")
	st.Expect(t, bar.Datacenter, "")
}

func TestMapConsulEntries(t *testing.T) {
	service := &consul.ServiceEntry{
		Node:    &consul.Node{Node: "foo", Address: "127.0.0.1"},
		Service: &consul.AgentService{Service: "foo", Port: 80},
	}
	list := []*consul.ServiceEntry{service}
	st.Expect(t, MapConsulEntries(list), []string{"127.0.0.1:80"})
}
