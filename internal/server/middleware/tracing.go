package middleware

import (
	"net/http"
	"time"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

var tracer = otel.Tracer("observability-api/http")

func Tracing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Start a span
		ctx, span := tracer.Start(
			r.Context(),
			r.Method+" "+r.URL.Path,
			trace.WithSpanKind(trace.SpanKindServer),
		)

		
		defer span.End()

		// Add request attributes
		span.SetAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.path", r.URL.Path),
		)

		// Wrap response writer to capture status code
		rec := &statusRecorder{
			ResponseWriter: w,
			status: http.StatusOK,
		}

		fmt.Println("Tracing middleware span valid:",
			span.SpanContext().IsValid(),
			span.SpanContext().TraceID().String(),
		)

		// Call next middleware / handler
		next.ServeHTTP(rec, r.WithContext(ctx))

		// Record response info
		span.SetAttributes(
			attribute.Int("http.status_code", rec.status),
			attribute.Int64("http.duration_ms", time.Since(start).Milliseconds()),
		)

		// Mark errors
		if rec.status >= 500 {
			span.SetStatus(codes.Error, "server error")
		}
	})
}
