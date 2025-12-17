package main

import (
    "net/http"
    "log/slog"
    "os"
    "context"
    "log"

    "github.com/traceylum1/observability-api/internal/server"
    "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func initLogger() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func initTracer() (func(context.Context) error, error) {
	exp, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("observability-api"),
		)),
	)

	otel.SetTracerProvider(tp)

	return tp.Shutdown, nil
}

func main() {
    r := server.NewRouter()

    initLogger()

    ctx := context.Background()

	shutdown, err := initTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer shutdown(ctx)

    slog.Info("server started", "port", 3000)
    http.ListenAndServe(":3000", r)
}