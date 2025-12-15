package server

import (
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
	"github.com/traceylum1/observability-api/internal/handlers"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", handlers.Hello)
	r.Get("/user", handlers.GetUser)
	r.Get("/user/:user_id", handlers.GetUserInfo)

	return r
}
