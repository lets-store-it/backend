package telemetry

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type tracedOperation[T any] struct {
	ctx    context.Context
	span   trace.Span
	result T
	err    error
}

func WithTrace[T any](ctx context.Context, tracer trace.Tracer, name string, fn func(context.Context, trace.Span) (T, error)) (T, error) {
	ctx, span := tracer.Start(ctx, name)
	defer span.End()

	op := &tracedOperation[T]{
		ctx:  ctx,
		span: span,
	}

	op.result, op.err = fn(ctx, span)
	if op.err != nil {
		span.RecordError(op.err)
		span.SetStatus(codes.Error, op.err.Error())
	} else {
		span.SetStatus(codes.Ok, fmt.Sprintf("%s completed successfully", name))
	}

	return op.result, op.err
}

func withTraceVoid(ctx context.Context, tracer trace.Tracer, name string, fn func(context.Context, trace.Span) error) error {
	_, err := WithTrace[struct{}](ctx, tracer, name, func(ctx context.Context, span trace.Span) (struct{}, error) {
		return struct{}{}, fn(ctx, span)
	})
	return err
}

func addAttributes(span trace.Span, attrs ...attribute.KeyValue) {
	span.SetAttributes(attrs...)
}
