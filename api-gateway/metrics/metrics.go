package metrics

import (
	"sync/atomic"
	"time"
)

type Collector struct {
	requests uint64
	errors   uint64
	totalNS  uint64
}

func New() *Collector {
	return &Collector{}
}

func (c *Collector) Observe(start time.Time, status int) {
	atomic.AddUint64(&c.requests, 1)
	atomic.AddUint64(&c.totalNS, uint64(time.Since(start).Nanoseconds()))
	if status >= 500 {
		atomic.AddUint64(&c.errors, 1)
	}
}

func (c *Collector) Snapshot() map[string]uint64 {
	req := atomic.LoadUint64(&c.requests)
	errCount := atomic.LoadUint64(&c.errors)
	total := atomic.LoadUint64(&c.totalNS)
	avg := uint64(0)
	if req > 0 {
		avg = total / req
	}

	return map[string]uint64{
		"requests_total":    req,
		"errors_total":      errCount,
		"avg_latency_ns":    avg,
	}
}
