package telemetry

import (
	"context"
	"fmt"
	"runtime/debug"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type TraceFn[T any] func(context.Context, trace.Span) (T, error)
type VoidTraceFn func(context.Context, trace.Span) error

func WithTrace[T any](
	ctx context.Context,
	tracer trace.Tracer,
	name string,
	fn TraceFn[T],
	attrs ...attribute.KeyValue,
) (result T, err error) {
	ctx, span := tracer.Start(ctx, name, trace.WithAttributes(attrs...))
	defer span.End()

	// Recover from panics and record them as errors
	defer func() {
		if r := recover(); r != nil {
			stack := string(debug.Stack())
			err = fmt.Errorf("panic in traced operation: %v\nstack:\n%s", r, stack)
			span.RecordError(err, trace.WithAttributes(
				attribute.String("panic.stack", stack),
				attribute.String("panic.value", fmt.Sprintf("%v", r)),
			))
			span.SetStatus(codes.Error, "panic in traced operation")
			panic(r) // Re-panic after recording
		}
	}()

	result, err = fn(ctx, span)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(
			attribute.Bool("error", true),
			attribute.String("error.type", fmt.Sprintf("%T", err)),
		)
	} else {
		span.SetStatus(codes.Ok, fmt.Sprintf("%s completed successfully", name))
		span.SetAttributes(attribute.Bool("error", false))
	}

	return result, err
}

func WithVoidTrace(
	ctx context.Context,
	tracer trace.Tracer,
	name string,
	fn VoidTraceFn,
	attrs ...attribute.KeyValue,
) error {
	_, err := WithTrace(ctx, tracer, name, func(ctx context.Context, span trace.Span) (struct{}, error) {
		if err := fn(ctx, span); err != nil {
			return struct{}{}, err
		}
		return struct{}{}, nil
	}, attrs...)
	return err
}
