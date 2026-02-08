package observability

import (
	"context"
	"os"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	// "github.com/prometheus/client_golang/prometheus"
	// "github.com/prometheus/client_golang/prometheus/promhttp"
	
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

func SetupObservability(ctx context.Context) (func(context.Context) error, error) {
	tpShutdown, err := initTracer()
	if err != nil {
		return nil, err
	}

	initLogger()
	
	return func(ctx context.Context) error {
		var err1 error
		if tpShutdown != nil {
			err1 = tpShutdown(ctx)
		}
		return err1
	}, nil
}
