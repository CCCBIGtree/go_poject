package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		log.Printf("method=%s path=%s status=%d latency=%s remote=%s", r.Method, r.URL.Path, rw.status, time.Since(start), r.RemoteAddr)
	})
}
