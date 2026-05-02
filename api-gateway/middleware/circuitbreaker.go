package middleware

import (
	"net/http"
	"sync"
	"time"
)

type State string

const (
	Closed   State = "closed"
	Open     State = "open"
	HalfOpen State = "half_open"
)

type CircuitBreaker struct {
	mu               sync.Mutex
	state            State
	failureCount     int
	successCount     int
	failureThreshold int
	successThreshold int
	openTimeout      time.Duration
	openedAt         time.Time
}

func NewCircuitBreaker(failureThreshold, successThreshold, openTimeoutSec int) *CircuitBreaker {
	return &CircuitBreaker{
		state:            Closed,
		failureThreshold: failureThreshold,
		successThreshold: successThreshold,
		openTimeout:      time.Duration(openTimeoutSec) * time.Second,
	}
}

func (c *CircuitBreaker) allow() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.state == Open {
		if time.Since(c.openedAt) > c.openTimeout {
			c.state = HalfOpen
			c.successCount = 0
			return true
		}
		return false
	}
	return true
}

func (c *CircuitBreaker) onSuccess() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.state == HalfOpen {
		c.successCount++
		if c.successCount >= c.successThreshold {
			c.state = Closed
			c.failureCount = 0
			c.successCount = 0
		}
		return
	}
	c.failureCount = 0
}

func (c *CircuitBreaker) onFailure() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.failureCount++
	if c.failureCount >= c.failureThreshold {
		c.state = Open
		c.openedAt = time.Now()
	}
}

func (c *CircuitBreaker) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if !c.allow() {
			http.Error(w, "service temporarily unavailable (circuit open)", http.StatusServiceUnavailable)
			return
		}

		rw := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, req)
		if rw.status >= 500 {
			c.onFailure()
			return
		}
		c.onSuccess()
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}
