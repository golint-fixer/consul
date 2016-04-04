package consul

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/nbio/st"
	"gopkg.in/h2non/gock.v0"
)

func TestConsulSimpleClient(t *testing.T) {
	defer gock.Off()

	config := NewConfig("web", "http://consul.io")
	consul := New(config)
	gock.InterceptClient(config.Instances[0].HttpClient)

	gock.New("http://consul.io").
		Get("/v1/health/service/web").
		Reply(200).
		Type("json").
		BodyString(consulResponse)

	req := &http.Request{URL: &url.URL{}}

	var called bool
	consul.HandleHTTP(nil, req, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	st.Expect(t, called, true)
	st.Expect(t, req.Host, "127.0.0.1:80")
	st.Expect(t, req.URL.Host, "127.0.0.1:80")
}
