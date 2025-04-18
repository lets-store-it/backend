package telemetry

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	otelprometheus "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

const appName = "storeit-backend"

type ShutdownFunc = func(ctx context.Context) error

var (
	tracerProvider *trace.TracerProvider
	meterProvider  *sdkmetric.MeterProvider
	shutdownFuncs  []ShutdownFunc
)

// InitTelemetry initializes global telemetry providers
func InitTelemetry(ctx context.Context) error {
	// Log all environment variables related to OTLP configuration
	log.Printf("OTLP Configuration:")
	log.Printf("OTEL_EXPORTER_OTLP_ENDPOINT=%s", os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	log.Printf("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT=%s", os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT"))
	log.Printf("OTEL_EXPORTER_OTLP_HEADERS=%s", os.Getenv("OTEL_EXPORTER_OTLP_HEADERS"))
	log.Printf("OTEL_EXPORTER_OTLP_TRACES_HEADERS=%s", os.Getenv("OTEL_EXPORTER_OTLP_TRACES_HEADERS"))
	log.Printf("OTEL_EXPORTER_OTLP_TIMEOUT=%s", os.Getenv("OTEL_EXPORTER_OTLP_TIMEOUT"))
	log.Printf("OTEL_EXPORTER_OTLP_TRACES_TIMEOUT=%s", os.Getenv("OTEL_EXPORTER_OTLP_TRACES_TIMEOUT"))

	// Create resource
	log.Printf("Creating telemetry resource...")
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

	// Initialize tracer provider with OTLP HTTP exporter
	log.Printf("Creating OTLP trace exporter...")
	traceExporter, err := otlptracehttp.New(ctx)
	if err != nil {
		return fmt.Errorf("failed to create trace exporter: %w", err)
	}

	log.Printf("Creating trace provider...")
	tp := trace.NewTracerProvider(
		trace.WithResource(appResource),
		trace.WithBatcher(traceExporter,
			trace.WithBatchTimeout(5*time.Second), // More frequent exports for debugging
			trace.WithMaxExportBatchSize(1),       // Export immediately for debugging
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

	// Create a test span to verify tracing is working
	log.Printf("Creating test span...")
	tr := tp.Tracer("test")
	ctx, span := tr.Start(ctx, "TestSpan")
	span.SetAttributes(
		attribute.String("test.key", "test.value"),
		attribute.String("service.name", appName),
	)
	span.End()
	log.Printf("Test span created and ended")

	// Force flush to ensure the test span is sent
	log.Printf("Force flushing test span...")
	if err := tp.ForceFlush(ctx); err != nil {
		log.Printf("Warning: Failed to flush test span: %v", err)
	} else {
		log.Printf("Test span flushed successfully")
	}

	// Initialize meter provider with Prometheus
	log.Printf("Creating Prometheus exporter...")
	prometheusExporter, err := otelprometheus.New(
		otelprometheus.WithRegisterer(prometheus.DefaultRegisterer),
	)
	if err != nil {
		return fmt.Errorf("failed to create prometheus exporter: %w", err)
	}

	log.Printf("Creating meter provider...")
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(appResource),
		sdkmetric.WithReader(prometheusExporter),
	)
	meterProvider = mp
	otel.SetMeterProvider(mp)
	shutdownFuncs = append(shutdownFuncs, mp.Shutdown)

	// Set global propagator
	log.Printf("Setting global propagator...")
	otel.SetTextMapPropagator(NewPropagator())

	log.Printf("Telemetry initialization completed successfully")
	return nil
}

// Shutdown gracefully shuts down all telemetry providers
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

// GetTracerProvider returns the global tracer provider
func GetTracerProvider() *trace.TracerProvider {
	return tracerProvider
}

// GetMeterProvider returns the global meter provider
func GetMeterProvider() *sdkmetric.MeterProvider {
	return meterProvider
}

// NewPropagator creates a new TextMapPropagator for distributed tracing
func NewPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}
