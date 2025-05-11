package organization

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/database"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services"
	"github.com/let-store-it/backend/internal/telemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	maxNameLength      = 100
	maxSubdomainLength = 63
	maxAliasLength     = 100
)

func validateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%w: name cannot be empty", services.ErrValidationError)
	}
	if len(name) > maxNameLength {
		return fmt.Errorf("%w: name is too long", services.ErrValidationError)
	}
	return nil
}

func validateSubdomain(subdomain string) error {
	if strings.TrimSpace(subdomain) == "" {
		return fmt.Errorf("%w: subdomain cannot be empty", services.ErrValidationError)
	}
	if len(subdomain) > maxSubdomainLength {
		return fmt.Errorf("%w: subdomain is too long", services.ErrValidationError)
	}
	matched, _ := regexp.MatchString("^[a-z0-9-]+$", subdomain)
	if !matched {
		return fmt.Errorf("%w: subdomain can only contain lowercase letters, numbers, and hyphens", services.ErrValidationError)
	}
	return nil
}

func validateAlias(alias string) error {
	if strings.TrimSpace(alias) == "" {
		return fmt.Errorf("%w: alias cannot be empty", services.ErrValidationError)
	}
	if len(alias) > maxAliasLength {
		return fmt.Errorf("%w: alias is too long", services.ErrValidationError)
	}
	matched, _ := regexp.MatchString("^[\\w-]+$", alias)
	if !matched {
		return fmt.Errorf("%w: alias can only contain letters, numbers, and hyphens (no spaces)", services.ErrValidationError)
	}
	return nil
}

type OrganizationService struct {
	queries *sqlc.Queries
	pgxPool *pgxpool.Pool
	tracer  trace.Tracer
}

type OrganizationServiceConfig struct {
	Queries *sqlc.Queries
	PGXPool *pgxpool.Pool
}

func New(cfg OrganizationServiceConfig) *OrganizationService {
	return &OrganizationService{
		queries: cfg.Queries,
		pgxPool: cfg.PGXPool,
		tracer:  otel.GetTracerProvider().Tracer("organization-service"),
	}
}

func (s *OrganizationService) CreateOrganization(ctx context.Context, name string, subdomain string) (*models.Organization, error) {
	return telemetry.WithTrace(ctx, s.tracer, "Create", func(ctx context.Context, span trace.Span) (*models.Organization, error) {
		span.SetAttributes(
			attribute.String("org.name", name),
			attribute.String("org.subdomain", subdomain),
		)

		if err := validateName(name); err != nil {
			return nil, err
		}
		if err := validateSubdomain(subdomain); err != nil {
			return nil, err
		}

		org, err := s.queries.CreateOrganization(ctx, sqlc.CreateOrganizationParams{
			Name:      name,
			Subdomain: subdomain,
		})
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		span.SetAttributes(attribute.String("org.id", org.ID.String()))
		return toOrganizationModel(org), nil
	})
}

func (s *OrganizationService) GetUsersOrganization(ctx context.Context, userID uuid.UUID) ([]*models.Organization, error) {
	return telemetry.WithTrace(ctx, s.tracer, "GetUsersOrgs", func(ctx context.Context, span trace.Span) ([]*models.Organization, error) {
		span.SetAttributes(
			attribute.String("user.id", userID.String()),
		)

		res, err := s.queries.GetUserOrgs(ctx, database.PgUUID(userID))
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		orgs := make([]*models.Organization, len(res))
		for i, org := range res {
			orgs[i] = toOrganizationModel(org)
		}

		return orgs, nil
	})
}

func (s *OrganizationService) GetOrganizationByID(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	return telemetry.WithTrace(ctx, s.tracer, "GetByID", func(ctx context.Context, span trace.Span) (*models.Organization, error) {
		span.SetAttributes(
			attribute.String("org.id", id.String()),
		)

		org, err := s.queries.GetOrganization(ctx, database.PgUUID(id))
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		return toOrganizationModel(org), nil
	})
}

func (s *OrganizationService) DeleteOrganization(ctx context.Context, id uuid.UUID) error {
	return telemetry.WithVoidTrace(ctx, s.tracer, "Delete", func(ctx context.Context, span trace.Span) error {
		span.SetAttributes(
			attribute.String("org.id", id.String()),
		)

		err := s.queries.DeleteOrganization(ctx, database.PgUUID(id))
		if err != nil {
			return services.MapDbErrorToService(err)
		}

		return nil
	})
}

func (s *OrganizationService) Update(ctx context.Context, org *models.Organization) (*models.Organization, error) {
	return telemetry.WithTrace(ctx, s.tracer, "Update", func(ctx context.Context, span trace.Span) (*models.Organization, error) {
		span.SetAttributes(
			attribute.String("org.id", org.ID.String()),
		)

		if err := validateName(org.Name); err != nil {
			return nil, err
		}

		updatedOrg, err := s.queries.UpdateOrganization(ctx, sqlc.UpdateOrganizationParams{
			ID:   database.PgUUID(org.ID),
			Name: org.Name,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to update organization: %w", err)
		}

		return toOrganizationModel(updatedOrg), nil
	})
}

// Organization Unit methods
func (s *OrganizationService) CreateUnit(ctx context.Context, orgID uuid.UUID, name string, alias string, address string) (*models.OrganizationUnit, error) {
	return telemetry.WithTrace(ctx, s.tracer, "CreateUnit", func(ctx context.Context, span trace.Span) (*models.OrganizationUnit, error) {
		span.SetAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("org_unit.name", name),
			attribute.String("org_unit.alias", alias),
			attribute.String("org_unit.address", address),
		)

		if err := validateName(name); err != nil {
			return nil, err
		}
		if err := validateAlias(alias); err != nil {
			return nil, err
		}

		unit, err := s.queries.CreateOrgUnit(ctx, sqlc.CreateOrgUnitParams{
			OrgID:   database.PgUUID(orgID),
			Name:    name,
			Alias:   alias,
			Address: database.PgText(address),
		})
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		return toOrganizationUnitModel(unit), nil
	})
}

func (s *OrganizationService) GetAllUnits(ctx context.Context, orgID uuid.UUID) ([]*models.OrganizationUnit, error) {
	return telemetry.WithTrace(ctx, s.tracer, "GetAllUnits", func(ctx context.Context, span trace.Span) ([]*models.OrganizationUnit, error) {
		span.SetAttributes(
			attribute.String("org.id", orgID.String()),
		)

		units, err := s.queries.GetOrgUnits(ctx, database.PgUUID(orgID))
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		result := make([]*models.OrganizationUnit, len(units))
		for i, unit := range units {
			result[i] = toOrganizationUnitModel(unit)
		}

		return result, nil
	})
}

func (s *OrganizationService) GetUnitByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*models.OrganizationUnit, error) {
	return telemetry.WithTrace(ctx, s.tracer, "GetUnitByID", func(ctx context.Context, span trace.Span) (*models.OrganizationUnit, error) {
		span.SetAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("unit.id", id.String()),
		)

		unit, err := s.queries.GetOrgUnitById(ctx, sqlc.GetOrgUnitByIdParams{
			OrgID: database.PgUUID(orgID),
			ID:    database.PgUUID(id),
		})
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		return toOrganizationUnitModel(unit), nil
	})
}

func (s *OrganizationService) DeleteUnit(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	return telemetry.WithVoidTrace(ctx, s.tracer, "DeleteUnit", func(ctx context.Context, span trace.Span) error {
		span.SetAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("unit.id", id.String()),
		)

		err := s.queries.DeleteOrgUnit(ctx, sqlc.DeleteOrgUnitParams{
			OrgID: database.PgUUID(orgID),
			ID:    database.PgUUID(id),
		})
		if err != nil {
			return services.MapDbErrorToService(err)
		}

		return nil
	})
}

func (s *OrganizationService) UpdateUnit(ctx context.Context, unit *models.OrganizationUnit) (*models.OrganizationUnit, error) {
	return telemetry.WithTrace(ctx, s.tracer, "UpdateUnit", func(ctx context.Context, span trace.Span) (*models.OrganizationUnit, error) {
		span.SetAttributes(
			attribute.String("org.id", unit.OrgID.String()),
			attribute.String("unit.id", unit.ID.String()),
		)

		if err := validateName(unit.Name); err != nil {
			return nil, err
		}
		if err := validateAlias(unit.Alias); err != nil {
			return nil, err
		}

		var address string
		if unit.Address != nil {
			address = *unit.Address
		}

		updatedUnit, err := s.queries.UpdateOrgUnit(ctx, sqlc.UpdateOrgUnitParams{
			ID:      database.PgUUID(unit.ID),
			OrgID:   database.PgUUID(unit.OrgID),
			Name:    unit.Name,
			Alias:   unit.Alias,
			Address: database.PgText(address),
		})
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		return toOrganizationUnitModel(updatedUnit), nil
	})
}
