package main

import (
	"flag"
	"log"
	"net/http"

	"go_poject/api-gateway/config"
	"go_poject/api-gateway/loadbalancer"
	"go_poject/api-gateway/metrics"
	"go_poject/api-gateway/middleware"
	"go_poject/api-gateway/router"
)

func main() {
	cfgPath := flag.String("config", "api-gateway/config/config.yaml", "config file path")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	m := metrics.New()
	lb := loadbalancer.NewRoundRobin()
	var handler http.Handler = router.NewGatewayRouter(cfg.Routes, lb, m)

	if cfg.Circuit.Enabled {
		cb := middleware.NewCircuitBreaker(
			cfg.Circuit.FailureThreshold,
			cfg.Circuit.SuccessThreshold,
			cfg.Circuit.OpenTimeoutSeconds,
		)
		handler = cb.Middleware(handler)
	}

	if cfg.RateLimit.Enabled {
		rl := middleware.NewRateLimiter(cfg.RateLimit.RPS, cfg.RateLimit.Burst)
		handler = rl.Middleware(handler)
	}

	handler = middleware.Logger(handler)

	log.Printf("gateway listening on %s", cfg.Server.Address)
	if err := http.ListenAndServe(cfg.Server.Address, handler); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
