package server

import (
    "github.com/go-chi/chi/v5"
    // chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/traceylum1/observability-api/internal/handlers"
	"github.com/traceylum1/observability-api/internal/server/middleware"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()
	// r.Use(chiMiddleware.Logger)
	r.Use(
        middleware.RequestID, 
        middleware.Tracing, 
		middleware.Logging,
		middleware.Metrics,
    )

	r.Get("/", handlers.Hello)
	
	r.Get("/items", handlers.GetUser)
	r.Get("/items/{id}", handlers.GetUserInfo)

	r.Get("/healthz", handlers.Live)
	r.Get("/readyz", handlers.Ready)

	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	return r
}
