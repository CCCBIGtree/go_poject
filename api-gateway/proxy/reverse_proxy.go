package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewReverseProxy(target string) (*httputil.ReverseProxy, error) {
	u, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	p := httputil.NewSingleHostReverseProxy(u)
	originalDirector := p.Director
	p.Director = func(req *http.Request) {
		originalDirector(req)
		req.Header.Set("X-Gateway", "go-api-gateway")
	}
	return p, nil
}
