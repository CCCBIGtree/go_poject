package router

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"go_poject/api-gateway/config"
	"go_poject/api-gateway/loadbalancer"
	"go_poject/api-gateway/metrics"
	"go_poject/api-gateway/proxy"
	"go_poject/api-gateway/web"
)

type GatewayRouter struct {
	routes  []config.RouteConfig
	lb      *loadbalancer.RoundRobin
	metrics *metrics.Collector
}

func NewGatewayRouter(routes []config.RouteConfig, lb *loadbalancer.RoundRobin, m *metrics.Collector) *GatewayRouter {
	return &GatewayRouter{routes: routes, lb: lb, metrics: m}
}

func (g *GatewayRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(web.DashboardHTML))
		return
	}

	if r.URL.Path == "/metrics" {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(g.metrics.Snapshot())
		return
	}

	start := time.Now()
	for _, route := range g.routes {
		if strings.HasPrefix(r.URL.Path, route.PathPrefix) {
			target := g.lb.Next(route.PathPrefix, route.Upstreams)
			if target == "" {
				http.Error(w, "no upstream available", http.StatusBadGateway)
				g.metrics.Observe(start, http.StatusBadGateway)
				return
			}

			p, err := proxy.NewReverseProxy(target)
			if err != nil {
				http.Error(w, "invalid upstream", http.StatusBadGateway)
				g.metrics.Observe(start, http.StatusBadGateway)
				return
			}
			rw := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			p.ServeHTTP(rw, r)
			g.metrics.Observe(start, rw.status)
			return
		}
	}

	http.NotFound(w, r)
	g.metrics.Observe(start, http.StatusNotFound)
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}
