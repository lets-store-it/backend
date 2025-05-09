package storage

import (
	"errors"
	"regexp"
	"strings"

	database "github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/services"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
// ErrStorageGroupNotFound = errors.New("storage group not found")
// ErrCellsGroupNotFound   = errors.New("cells group not found")
// ErrCellNotFound         = errors.New("cell not found")

// ErrInvalidStorageGroup = errors.New("invalid storage group")
// ErrInvalidCellsGroup   = errors.New("invalid cells group")
// ErrInvalidCell         = errors.New("invalid cell")
)

const (
	maxNameLength  = 100
	maxAliasLength = 100
	aliasPattern   = "^[\\w-]+$"
)

type StorageService struct {
	queries *database.Queries
	tracer  trace.Tracer
}

type StorageServiceConfig struct {
	Queries *database.Queries
}

func New(cfg *StorageServiceConfig) (*StorageService, error) {
	if cfg == nil || cfg.Queries == nil {
		return nil, errors.New("invalid configuration")
	}

	return &StorageService{
		queries: cfg.Queries,
		tracer:  otel.GetTracerProvider().Tracer("storage-service"),
	}, nil
}

func (s *StorageService) validateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.Join(services.ErrValidationError, errors.New("name cannot be empty"))
	}
	if len(name) > maxNameLength {
		return errors.Join(services.ErrValidationError, errors.New("name is too long (max 100 characters)"))
	}
	return nil
}

func (s *StorageService) validateAlias(alias string) error {
	if strings.TrimSpace(alias) == "" {
		return errors.Join(services.ErrValidationError, errors.New("alias cannot be empty"))
	}
	if len(alias) > maxAliasLength {
		return errors.Join(services.ErrValidationError, errors.New("alias is too long (max 100 characters)"))
	}
	matched, _ := regexp.MatchString(aliasPattern, alias)
	if !matched {
		return errors.Join(services.ErrValidationError, errors.New("alias can only contain letters, numbers, and hyphens (no spaces)"))
	}
	return nil
}
