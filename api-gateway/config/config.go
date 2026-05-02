package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server     ServerConfig             `yaml:"server"`
	RateLimit  RateLimitConfig          `yaml:"rate_limit"`
	Circuit    CircuitBreakerConfig     `yaml:"circuit_breaker"`
	Routes     []RouteConfig            `yaml:"routes"`
	Middleware []string                 `yaml:"middleware"`
	Metadata   map[string]string        `yaml:"metadata"`
}

type ServerConfig struct {
	Address string `yaml:"address"`
}

type RateLimitConfig struct {
	Enabled bool    `yaml:"enabled"`
	RPS     float64 `yaml:"rps"`
	Burst   int     `yaml:"burst"`
}

type CircuitBreakerConfig struct {
	Enabled             bool `yaml:"enabled"`
	FailureThreshold    int  `yaml:"failure_threshold"`
	SuccessThreshold    int  `yaml:"success_threshold"`
	OpenTimeoutSeconds  int  `yaml:"open_timeout_seconds"`
	SlidingWindowSize   int  `yaml:"sliding_window_size"`
}

type RouteConfig struct {
	PathPrefix string   `yaml:"path_prefix"`
	Upstreams  []string `yaml:"upstreams"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.Server.Address == "" {
		cfg.Server.Address = ":8080"
	}
	if cfg.RateLimit.Burst == 0 {
		cfg.RateLimit.Burst = 100
	}
	if cfg.RateLimit.RPS == 0 {
		cfg.RateLimit.RPS = 100
	}
	if cfg.Circuit.FailureThreshold == 0 {
		cfg.Circuit.FailureThreshold = 5
	}
	if cfg.Circuit.SuccessThreshold == 0 {
		cfg.Circuit.SuccessThreshold = 2
	}
	if cfg.Circuit.OpenTimeoutSeconds == 0 {
		cfg.Circuit.OpenTimeoutSeconds = 10
	}
	if cfg.Circuit.SlidingWindowSize == 0 {
		cfg.Circuit.SlidingWindowSize = 20
	}

	return &cfg, nil
}
