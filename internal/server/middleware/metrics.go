package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
)

func routePattern(r *http.Request) string {
	if rctx := chi.RouteContext(r.Context()); rctx != nil {
		if pattern := rctx.RoutePattern(); pattern != "" {
			return pattern
		}
	}
	return "unknown"
}

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "http",
			Subsystem: "server",
			Name:      "requests_total",
			Help:      "Total number of HTTP requests.",
		},
		[]string{"method", "route", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "http",
			Subsystem: "server",
			Name:      "request_duration_seconds",
			Help:      "HTTP request latency.",
			// Good default buckets for APIs
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "route"},
	)
)

func init() {
	// Register metrics once at startup
	prometheus.MustRegister(
		httpRequestsTotal,
		httpRequestDuration,
	)
}

// Metrics instruments incoming HTTP requests.
func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Capture status code
		rec := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(rec, r)

		elapsed := time.Since(start).Seconds()

		route := routePattern(r)
		method := r.Method
		status := strconv.Itoa(rec.status)

		httpRequestsTotal.WithLabelValues(
			method,
			route,
			status,
		).Inc()

		httpRequestDuration.WithLabelValues(
			method,
			route,
		).Observe(elapsed)
	})
}