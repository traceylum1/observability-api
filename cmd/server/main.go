package main

import (
    "net/http"
    "log/slog"
    "context"
    "log"
	"time"
	"os"
	"os/signal"
	"syscall"


    "github.com/traceylum1/observability-api/internal/server"
	"github.com/traceylum1/observability-api/internal/observability"
)



func main() {
    r := server.NewRouter()

    ctx := context.Background()

	// ---- Observability setup ----
	shutdownObs, err := observability.SetupObservability(ctx)
	if err != nil {
		log.Fatal(err)
	}

    slog.Info("server started", "port", 3000)
    server := &http.Server{
		Addr: ":3000",
		Handler: r,
	}

	// ---- Start server ----
	go func() {
		slog.Info("HTTP server listening on :3000")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// ---- Wait for shutdown signal ----
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	slog.Info("shutdown signal received")

	// ---- Graceful shutdown ----
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Stop accepting new requests & wait for in-flight
	if err = server.Shutdown(shutdownCtx); err != nil {
		slog.Info("server shutdown failed", "error", err)
	}

	// 2. Flush observability (traces, metrics)
	if err = shutdownObs(shutdownCtx); err != nil {
		slog.Info("observability shutdown error", "error", err)
	}

	slog.Info("server exited cleanly")
}