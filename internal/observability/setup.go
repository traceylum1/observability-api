package observability

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	// "github.com/prometheus/client_golang/prometheus"
	// "github.com/prometheus/client_golang/prometheus/promhttp"
	
)

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

func setupObservability(ctx context.Context) (func(context.Context) error, error) {
	tpShutdown, err := initTracer()
	if err != nil {
		return nil, err
	}
	
}


// type metrics struct {
// 	cpuTemp  prometheus.Gauge
// 	hdFailures *prometheus.CounterVec
// }

// func NewMetrics(reg prometheus.Registerer) *metrics {
// 	m := &metrics{
// 		cpuTemp: prometheus.NewGauge(
// 			prometheus.GaugeOpts{
// 				Name: "cpu_temperature_celsius",
// 				Help: "Current temperature of the CPU.",
// 		}),
// 		hdFailures: prometheus.NewCounterVec(
// 			prometheus.CounterOpts{
// 				Name: "hd_errors_total",
// 				Help: "Number of hard-disk errors.",
// 			},
// 			[]string{"device"},
// 		),
// 	}
// 	reg.MustRegister(m.cpuTemp)
// 	reg.MustRegister(m.hdFailures)
// 	return m
// }

	// // Create a non-global registry.
	// reg := prometheus.NewRegistry()

	// // Create new metrics and register them using the custom registry.
	// m := NewMetrics(reg)
	// // Set values for the new created metrics.
	// m.cpuTemp.Set(65.3)
	// m.hdFailures.With(prometheus.Labels{"device":"/dev/sda"}).Inc()

	// http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))