package telemetry

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	otelprometheus "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

const appName = "storeit-backend"

type ShutdownFunc = func(ctx context.Context) error

// NewPropagator creates a new TextMapPropagator for distributed tracing
func NewPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

// NewTracerProvider creates and configures a new TracerProvider
func NewTracerProvider() (*trace.TracerProvider, ShutdownFunc, error) {
	appResource, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(appName),
		),
	)
	if err != nil {
		return nil, nil, err
	}

	traceExporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		return nil, nil, err
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithResource(appResource),
		trace.WithBatcher(traceExporter, trace.WithBatchTimeout(time.Second)),
	)

	return tracerProvider, func(ctx context.Context) error {
		_ = tracerProvider.ForceFlush(ctx)
		return tracerProvider.Shutdown(ctx)
	}, nil
}

// NewMeterProvider creates and configures a new MeterProvider
func NewMeterProvider() (metric.MeterProvider, ShutdownFunc, error) {
	appResource, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(appName),
		),
	)
	if err != nil {
		return nil, nil, err
	}

	prometheusExporter, err := otelprometheus.New(
		otelprometheus.WithRegisterer(prometheus.DefaultRegisterer),
	)
	if err != nil {
		return nil, nil, err
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(appResource),
		sdkmetric.WithReader(prometheusExporter),
	)

	return meterProvider, func(ctx context.Context) error {
		_ = meterProvider.ForceFlush(ctx)
		return meterProvider.Shutdown(ctx)
	}, nil
}
