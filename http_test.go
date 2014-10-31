package main_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/lox/apt-proxy/proxy"
)

type testFixture struct {
	proxy, backend *httptest.Server
}

func newTestFixture(handler http.HandlerFunc, ap *proxy.AptProxy) *testFixture {
	backend := httptest.NewServer(http.HandlerFunc(handler))
	backendURL, err := url.Parse(backend.URL)
	if err != nil {
		panic(err)
	}

	ap.Rewriters = append(ap.Rewriters, proxy.RewriterFunc(func(req *http.Request) {
		req.URL.Host = backendURL.Host
	}))

	return &testFixture{httptest.NewServer(ap.Handler()), backend}
}

// client returns an http client configured to use the provided proxy
func (f *testFixture) client() *http.Client {
	proxyURL, err := url.Parse(f.proxy.URL)
	if err != nil {
		panic(err)
	}

	return &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
}

func (f *testFixture) close() {
	f.proxy.Close()
	f.backend.Close()
}

func assertHeader(t *testing.T, r *http.Response, header string, expected string) {
	if r.Header.Get(header) != expected {
		t.Fatalf("Expected header %s=%s, but got '%s'",
			header, expected, r.Header.Get(header))
	}
}

func TestProxyAddsCacheControlHeader(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Llamas rock"))
	}
	ap := proxy.NewAptProxy()
	ap.CachePatterns = append(ap.CachePatterns, proxy.NewCachePattern(".", time.Hour*100))
	fixture := newTestFixture(handler, ap)
	defer fixture.close()

	resp1, err := fixture.client().Get(fixture.backend.URL)
	if err != nil {
		t.Fatal(err)
	}

	assertHeader(t, resp1, "Cache-Control", "max-age=360000")
}

func TestRewritesApply(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Llamas rock"))
	}
	fixture := newTestFixture(handler, proxy.NewAptProxy())
	defer fixture.close()

	resp, err := fixture.client().Get(
		"http://archive.ubuntu.com/ubuntu/pool/main/b/bind9/dnsutils_9.9.5.dfsg-3_amd64.deb")
	if err != nil {
		t.Fatal(err)
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(contents, []byte("Llamas rock")) != 0 {
		t.Fatalf("Response content was incorrect, rewrites not applying?")
	}
}
