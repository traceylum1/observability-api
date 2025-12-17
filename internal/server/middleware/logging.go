package middleware

import (
	"net/http"
	"time"
	"fmt"
	"log/slog"

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

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		recorder := &statusRecorder{
			ResponseWriter: w,
			status: http.StatusOK,
		}

		next.ServeHTTP(recorder, r)

		duration := time.Since(start)

		logger := slog.Default()

		// Attach request ID if present
		if reqID, ok := r.Context().Value(requestIDKey).(string); ok {
			logger = logger.With("request_id", reqID)
		}

		span := trace.SpanFromContext(r.Context())
		sc := span.SpanContext()

		fmt.Println("=== LOGGING MIDDLEWARE DEBUG ===")
		fmt.Println("Span type:", fmt.Sprintf("%T", span))
		fmt.Println("SpanContext valid:", sc.IsValid())
		fmt.Println("TraceID:", sc.TraceID().String())
		fmt.Println("SpanID:", sc.SpanID().String())
		fmt.Println("IsRecording:", span.IsRecording())
		fmt.Println("================================")

		// Attach trace info if present
		if span := trace.SpanFromContext(r.Context()); span.SpanContext().IsValid() {
			sc := span.SpanContext()
			logger = logger.With(
				"trace_id", sc.TraceID().String(),
				"span_id", sc.SpanID().String(),
			)
		}

		logger.Info("http request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", recorder.status,
			"duration_ms", duration.Milliseconds(),
		)
	})
}
