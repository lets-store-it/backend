package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type TxFn[T any] func(ctx context.Context, tx pgx.Tx) (T, error)
type VoidTxFn func(ctx context.Context, tx pgx.Tx) error

func WithTransaction[T any](ctx context.Context, pool *pgxpool.Pool, tracer trace.Tracer, f TxFn[T]) (result T, err error) {
	ctx, span := tracer.Start(ctx, "database.WithTransaction")
	defer span.End()

	tx, err := pool.Begin(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to begin transaction")
		return result, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			span.RecordError(fmt.Errorf("panic in transaction: %v", p))
			span.SetStatus(codes.Error, "panic in transaction")

			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				span.RecordError(rollbackErr)
				err = fmt.Errorf("panic: %v, rollback failed: %w", p, rollbackErr)
				return
			}
			panic(p)
		}

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "transaction failed")

			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				span.RecordError(rollbackErr)
				err = fmt.Errorf("original error: %v, rollback failed: %w", err, rollbackErr)
			}
			return
		}

		if commitErr := tx.Commit(ctx); commitErr != nil {
			span.RecordError(commitErr)
			span.SetStatus(codes.Error, "failed to commit transaction")
			err = fmt.Errorf("failed to commit transaction: %w", commitErr)
			return
		}

		span.SetStatus(codes.Ok, "transaction completed successfully")
	}()

	result, err = f(ctx, tx)
	if err != nil {
		span.SetAttributes(attribute.Bool("transaction.success", false))
		return result, err
	}

	span.SetAttributes(attribute.Bool("transaction.success", true))
	return result, nil
}

func WithVoidTransaction(ctx context.Context, pool *pgxpool.Pool, tracer trace.Tracer, f VoidTxFn) error {
	_, err := WithTransaction(ctx, pool, tracer, func(ctx context.Context, tx pgx.Tx) (struct{}, error) {
		return struct{}{}, f(ctx, tx)
	})
	return err
}
