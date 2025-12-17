package main

import (
    "net/http"
    "log/slog"
    "os"

    "github.com/traceylum1/observability-api/internal/server"
)

func initLogger() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func main() {
    r := server.NewRouter()

    initLogger()

    slog.Info("server started", "port", 3000)
    http.ListenAndServe(":3000", r)
}