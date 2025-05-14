package storage

import (
	"errors"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	database "github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/common"
	"github.com/let-store-it/backend/internal/services/audit"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const (
	maxNameLength  = 100
	maxAliasLength = 100
	aliasPattern   = "^[\\w-]+$"
)

type StorageService struct {
	queries *database.Queries
	pgxPool *pgxpool.Pool
	tracer  trace.Tracer
	audit   *audit.AuditService
}

type StorageServiceConfig struct {
	Queries *database.Queries
	PGXPool *pgxpool.Pool
	Audit   *audit.AuditService
}

func New(cfg *StorageServiceConfig) *StorageService {
	if cfg == nil || cfg.Queries == nil || cfg.PGXPool == nil || cfg.Audit == nil {
		panic("invalid configuration")
	}

	return &StorageService{
		queries: cfg.Queries,
		pgxPool: cfg.PGXPool,
		tracer:  otel.GetTracerProvider().Tracer("storage-service"),
		audit:   cfg.Audit,
	}
}

func (s *StorageService) validateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.Join(common.ErrValidationError, errors.New("name cannot be empty"))
	}
	if len(name) > maxNameLength {
		return errors.Join(common.ErrValidationError, errors.New("name is too long (max 100 characters)"))
	}
	return nil
}

func (s *StorageService) validateAlias(alias string) error {
	if strings.TrimSpace(alias) == "" {
		return errors.Join(common.ErrValidationError, errors.New("alias cannot be empty"))
	}
	if len(alias) > maxAliasLength {
		return errors.Join(common.ErrValidationError, errors.New("alias is too long (max 100 characters)"))
	}
	matched, _ := regexp.MatchString(aliasPattern, alias)
	if !matched {
		return errors.Join(common.ErrValidationError, errors.New("alias can only contain letters, numbers, and hyphens (no spaces)"))
	}
	return nil
}
