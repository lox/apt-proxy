package proxy

import (
	"fmt"
	"net/http"
	"regexp"
	"time"
)

// New creates a new pattern, which is a regular expression and a duration to
// override the Cache-Control max-age of any responses that match
func NewCachePattern(pattern string, d time.Duration) *cachePattern {
	return &cachePattern{Regexp: regexp.MustCompile(pattern), Duration: d}
}

type cachePattern struct {
	*regexp.Regexp
	Duration time.Duration
}

type CachePatternSlice []*cachePattern

// MatchString tries to match a given string across all patterns
func (r CachePatternSlice) MatchString(subject string) (bool, *cachePattern) {
	for _, p := range r {
		if p.MatchString(subject) {
			return true, p
		}
	}

	return false, nil
}

type cachePatternTransport struct {
	patterns CachePatternSlice
	upstream http.RoundTripper
}

func (t cachePatternTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	resp, err := t.upstream.RoundTrip(r)
	if match, pattern := t.patterns.MatchString(r.URL.Path); match {
		resp.Header.Add("Cache-Control", fmt.Sprintf("max-age=%.f", pattern.Duration.Seconds()))
	}
	return resp, err
}
