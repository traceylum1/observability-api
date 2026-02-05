package middleware

import (
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"github.com/traceylum1/observability-api/internal/observability"
)


var tracer = otel.Tracer("observability-api/http")

func Tracing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start, _ := observability.RequestStartFromContext(r.Context())

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

		// Call next middleware / handler
		next.ServeHTTP(w, r.WithContext(ctx))

		status := 0
		if rec, ok := w.(*statusRecorder); ok {
			status = rec.status
		}
		
		// Record response info
		span.SetAttributes(
			attribute.Int("http.status_code", status),
			attribute.Int64("http.duration_ms", time.Since(start).Milliseconds()),
		)

		// Mark errors
		if status >= 500 {
			span.SetStatus(codes.Error, "server error")
		}
	})
}
