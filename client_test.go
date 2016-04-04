package consul

import (
	"testing"

	consul "github.com/hashicorp/consul/api"
	"github.com/nbio/st"
	"gopkg.in/h2non/gock.v0"
)

const consulResponse = `
[
  {
    "Node":{
      "Node":"foo",
      "Address":"127.0.0.1",
      "TaggedAddresses":{
        "wan":"127.0.0.1"
      },
      "CreateIndex":7,
      "ModifyIndex":375588
    },
    "Service":{
      "ID":"web",
      "Service":"web",
      "Tags":null,
      "Address":"",
      "Port":80,
      "EnableTagOverride":false,
      "CreateIndex":13,
      "ModifyIndex":13
    }
  }
]`

func TestClient(t *testing.T) {
	defer gock.Off()

	gock.New("http://demo.consul.io").
		Get("/v1/health/service/web").
		Reply(200).
		Type("json").
		BodyString(consulResponse)

	gock.New("http://127.0.0.1:80").
		Get("/").
		Reply(200).
		BodyString("hello world")

	config := consul.DefaultConfig()
	config.Address = "demo.consul.io"
	gock.InterceptClient(config.HttpClient)

	client := NewClient(config)
	entries, _, err := client.Health("web", "", nil)
	st.Expect(t, err, nil)
	st.Expect(t, len(entries), 1)

	entry := entries[0]
	st.Expect(t, entry.Node.Node, "foo")
	st.Expect(t, entry.Node.Address, "127.0.0.1")
	st.Expect(t, entry.Service.Service, "web")
	st.Expect(t, entry.Service.Port, 80)
}
