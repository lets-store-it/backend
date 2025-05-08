package telemetry

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	otelprometheus "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type ShutdownFunc = func(ctx context.Context) error

var (
	tracerProvider *trace.TracerProvider
	meterProvider  *sdkmetric.MeterProvider
	shutdownFuncs  []ShutdownFunc
)

func InitTelemetry(ctx context.Context, appName string) error {
	appResource, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(appName),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	// Tracing
	traceExporter, err := otlptracehttp.New(ctx)
	if err != nil {
		return fmt.Errorf("failed to create trace exporter: %w", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithResource(appResource),
		trace.WithBatcher(traceExporter,
			trace.WithBatchTimeout(5*time.Second),
			// trace.WithMaxExportBatchSize(1),
		),
	)
	tracerProvider = tp
	otel.SetTracerProvider(tp)

	shutdownFuncs = append(shutdownFuncs, func(ctx context.Context) error {
		log.Printf("Flushing and shutting down trace provider...")
		if err := tp.ForceFlush(ctx); err != nil {
			log.Printf("Error flushing traces: %v", err)
		}
		return tp.Shutdown(ctx)
	})

	// Metrics
	prometheusExporter, err := otelprometheus.New(
		otelprometheus.WithRegisterer(prometheus.DefaultRegisterer),
	)
	if err != nil {
		return fmt.Errorf("failed to create prometheus exporter: %w", err)
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(appResource),
		sdkmetric.WithReader(prometheusExporter),
	)
	meterProvider = mp
	otel.SetMeterProvider(mp)
	shutdownFuncs = append(shutdownFuncs, mp.Shutdown)

	otel.SetTextMapPropagator(newPropagator())

	return nil
}

func Shutdown(ctx context.Context) error {
	log.Printf("Shutting down telemetry providers")
	var lastErr error
	for _, shutdown := range shutdownFuncs {
		if err := shutdown(ctx); err != nil {
			lastErr = err
			log.Printf("Error during telemetry shutdown: %v", err)
		}
	}
	return lastErr
}

func GetTracerProvider() *trace.TracerProvider {
	return tracerProvider
}

func GetMeterProvider() *sdkmetric.MeterProvider {
	return meterProvider
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}
