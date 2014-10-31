package proxy

import "net/http"

type Rewriter interface {
	Rewrite(req *http.Request)
}

type RewriterFunc func(req *http.Request)

// Rewrite calls f(req).
func (f RewriterFunc) Rewrite(req *http.Request) {
	f(req)
}
