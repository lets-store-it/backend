package auth

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/services/audit"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type AuthService struct {
	queries *sqlc.Queries
	pgxPool *pgxpool.Pool
	tracer  trace.Tracer

	auditService *audit.AuditService
}

type AuthServiceConfig struct {
	Queries *sqlc.Queries
	PGXPool *pgxpool.Pool

	AuditService *audit.AuditService
}

func New(cfg AuthServiceConfig) *AuthService {
	if cfg.Queries == nil || cfg.PGXPool == nil || cfg.AuditService == nil {
		panic("AuthServiceConfig is invalid")
	}

	return &AuthService{
		queries:      cfg.Queries,
		pgxPool:      cfg.PGXPool,
		tracer:       otel.GetTracerProvider().Tracer("auth-service"),
		auditService: cfg.AuditService,
	}
}
