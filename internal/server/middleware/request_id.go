package middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/traceylum1/observability-api/internal/observability"
)


func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		start := time.Now()

		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.NewString()
		}

		ctx = observability.WithRequestID(ctx, reqID)
		ctx = observability.WithRequestStart(ctx, start)


		w.Header().Set("X-Request-ID", reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}