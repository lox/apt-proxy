package proxy

import (
	"net/http"
	"net/http/httputil"

	"github.com/lox/apt-proxy/ubuntu"
)

type AptProxy struct {
	Transport     http.RoundTripper
	Rewriters     []Rewriter
	CachePatterns CachePatternSlice
}

func NewAptProxy() *AptProxy {
	return &AptProxy{
		Transport:     http.DefaultTransport,
		Rewriters:     []Rewriter{ubuntu.NewRewriter()},
		CachePatterns: CachePatternSlice{},
	}
}

func (ap *AptProxy) Handler() http.Handler {
	return &handler{ap: ap, Handler: &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			for _, rewrite := range ap.Rewriters {
				rewrite.Rewrite(r)
			}

			r.Host = r.URL.Host
		},
		Transport: &cachePatternTransport{ap.CachePatterns, ap.Transport},
	}}
}

type handler struct {
	http.Handler
	ap *AptProxy
}
