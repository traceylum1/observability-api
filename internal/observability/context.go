package observability

import (
	"context"
	"time"
)

type requestIDKey struct{}
type requestStartKey struct{}

// WithRequestID stores the requestID in context.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

// RequestIDFromContext returns the requestID.
func RequestIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(requestIDKey{}).(string)
	return id, ok
}

// WithRequestStart stores the canonical request start time in context.
func WithRequestStart(ctx context.Context, t time.Time) context.Context {
	return context.WithValue(ctx, requestStartKey{}, t)
}

// RequestStartFromContext returns the request start time if present.
func RequestStartFromContext(ctx context.Context) (time.Time, bool) {
	t, ok := ctx.Value(requestStartKey{}).(time.Time)
	return t, ok
}

// RequestDuration returns time since request start, or since now if missing.
func RequestDuration(ctx context.Context) time.Duration {
	if t, ok := RequestStartFromContext(ctx); ok {
		return time.Since(t)
	}
	return 0
}