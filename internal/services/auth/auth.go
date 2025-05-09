package auth

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/let-store-it/backend/generated/sqlc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type AuthService struct {
	queries *sqlc.Queries
	pgxPool *pgxpool.Pool
	tracer  trace.Tracer
}

func New(queries *sqlc.Queries, pgxPool *pgxpool.Pool) *AuthService {
	return &AuthService{
		queries: queries,
		pgxPool: pgxPool,
		tracer:  otel.GetTracerProvider().Tracer("auth-service"),
	}
}
