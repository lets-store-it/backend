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

func safeEndSpan(span trace.Span) {
	defer func() {
		recover() // Silently recover from any panic during span.End()
	}()
	if span != nil {
		span.End()
	}
}

func safeRecordError(span trace.Span, err error) {
	defer func() {
		recover() // Silently recover from any panic during error recording
	}()
	if span != nil && err != nil {
		span.RecordError(err)
	}
}

func safeSetStatus(span trace.Span, code codes.Code, description string) {
	defer func() {
		recover() // Silently recover from any panic during status setting
	}()
	if span != nil {
		span.SetStatus(code, description)
	}
}

func safeSetAttributes(span trace.Span, attrs ...attribute.KeyValue) {
	defer func() {
		recover() // Silently recover from any panic during attribute setting
	}()
	if span != nil {
		span.SetAttributes(attrs...)
	}
}

func WithTrace[T any](
	ctx context.Context,
	tracer trace.Tracer,
	name string,
	fn TraceFn[T],
	attrs ...attribute.KeyValue,
) (result T, err error) {
	if tracer == nil {
		return fn(ctx, trace.SpanFromContext(ctx))
	}

	var span trace.Span
	defer func() {
		if r := recover(); r != nil {
			stack := string(debug.Stack())
			err = fmt.Errorf("panic in traced operation: %v\nstack:\n%s", r, stack)

			if span != nil {
				safeRecordError(span, err)
				safeSetStatus(span, codes.Error, "panic in traced operation")
				safeSetAttributes(span,
					attribute.String("panic.stack", stack),
					attribute.String("panic.value", fmt.Sprintf("%v", r)),
				)
				safeEndSpan(span)
			}
			panic(r) // Re-panic after recording
		}
	}()

	ctx, span = tracer.Start(ctx, name, trace.WithAttributes(attrs...))
	defer func() {
		if span != nil {
			if err != nil {
				safeRecordError(span, err)
				safeSetStatus(span, codes.Error, err.Error())
				safeSetAttributes(span,
					attribute.Bool("error", true),
					attribute.String("error.type", fmt.Sprintf("%T", err)),
				)
			} else {
				safeSetStatus(span, codes.Ok, fmt.Sprintf("%s completed successfully", name))
				safeSetAttributes(span, attribute.Bool("error", false))
			}
			safeEndSpan(span)
		}
	}()

	return fn(ctx, span)
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
