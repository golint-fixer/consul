package consul

import (
	"net/http"
	"net/url"
	"testing"

	consul "github.com/hashicorp/consul/api"
	"github.com/nbio/st"
	"gopkg.in/h2non/gock.v0"
	"gopkg.in/vinxi/utils.v0"
)

func TestConsulSimpleClient(t *testing.T) {
	defer gock.Off()

	gock.New("http://consul.io").
		Get("/v1/health/service/web").
		Reply(200).
		Type("json").
		BodyString(consulResponse)

	config := NewConfig("web", "http://demo.consul.io")
	gock.InterceptClient(config.Instances[0].HttpClient)
	consul := New(config)

	w := utils.NewWriterStub()
	req := &http.Request{URL: &url.URL{}}

	var called bool
	consul.HandleHTTP(w, req, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	st.Expect(t, called, true)
	st.Expect(t, req.Host, "127.0.0.1:80")
	st.Expect(t, req.URL.Host, "127.0.0.1:80")
	st.Expect(t, w.Code, 200)
	st.Expect(t, string(w.Body), "")
}

func TestConsulCustomMapper(t *testing.T) {
	defer gock.Off()

	gock.New("http://consul.io").
		Get("/v1/health/service/web").
		Reply(200).
		Type("json").
		BodyString(consulResponse)

	config := NewConfig("web", "http://demo.consul.io")
	config.Mapper = func(list []*consul.ServiceEntry) []string {
		return MapConsulEntries(list)
	}
	gock.InterceptClient(config.Instances[0].HttpClient)
	consul := New(config)

	w := utils.NewWriterStub()
	req := &http.Request{URL: &url.URL{}}

	var called bool
	consul.HandleHTTP(w, req, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	st.Expect(t, called, true)
	st.Expect(t, req.Host, "127.0.0.1:80")
	st.Expect(t, req.URL.Host, "127.0.0.1:80")
	st.Expect(t, w.Code, 200)
	st.Expect(t, string(w.Body), "")
}
