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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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

func (s *OrganizationService) Create(ctx context.Context, name string, subdomain string) (*models.Organization, error) {
	ctx, span := s.tracer.Start(ctx, "Create",
		trace.WithAttributes(
			attribute.String("org.name", name),
			attribute.String("org.subdomain", subdomain),
		),
	)
	defer span.End()

	if err := validateName(name); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "name validation failed")
		return nil, err
	}
	if err := validateSubdomain(subdomain); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "subdomain validation failed")
		return nil, err
	}

	org, err := s.queries.CreateOrganization(ctx, sqlc.CreateOrganizationParams{
		Name:      name,
		Subdomain: subdomain,
	})
	if err != nil {
		if database.IsUniqueViolation(err) {
			span.RecordError(err)
			span.SetStatus(codes.Error, "organization already exists")
			return nil, services.ErrDuplicationError
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create organization")
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	span.SetAttributes(attribute.String("org.id", org.ID.String()))
	span.SetStatus(codes.Ok, "organization created successfully")

	return toOrganization(org)
}

func (s *OrganizationService) GetUsersOrgs(ctx context.Context, userID uuid.UUID) ([]*models.Organization, error) {
	ctx, span := s.tracer.Start(ctx, "GetUsersOrgs",
		trace.WithAttributes(
			attribute.String("user.id", userID.String()),
		),
	)
	defer span.End()

	res, err := s.queries.GetUserOrgs(ctx, database.PgUUID(userID))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get user organizations")
		return nil, fmt.Errorf("failed to get user organizations: %w", err)
	}

	orgs := make([]*models.Organization, len(res))
	for i, org := range res {
		orgs[i], err = toOrganization(org)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed to convert organization")
			return nil, fmt.Errorf("failed to convert organization: %w", err)
		}
	}

	span.SetStatus(codes.Ok, "successfully retrieved user organizations")
	return orgs, nil
}

func (s *OrganizationService) GetByID(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	ctx, span := s.tracer.Start(ctx, "GetByID",
		trace.WithAttributes(
			attribute.String("org.id", id.String()),
		),
	)
	defer span.End()

	org, err := s.queries.GetOrganization(ctx, database.PgUUID(id))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get organization")
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}

	span.SetStatus(codes.Ok, "successfully retrieved organization")
	return toOrganization(org)
}

func (s *OrganizationService) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, span := s.tracer.Start(ctx, "Delete",
		trace.WithAttributes(
			attribute.String("org.id", id.String()),
		),
	)
	defer span.End()

	err := s.queries.DeleteOrganization(ctx, database.PgUUID(id))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to delete organization")
		return fmt.Errorf("failed to delete organization: %w", err)
	}

	span.SetStatus(codes.Ok, "organization deleted successfully")
	return nil
}

func (s *OrganizationService) Update(ctx context.Context, org *models.Organization) (*models.Organization, error) {
	ctx, span := s.tracer.Start(ctx, "Update",
		trace.WithAttributes(
			attribute.String("org.id", org.ID.String()),
		),
	)
	defer span.End()

	if err := validateName(org.Name); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "name validation failed")
		return nil, err
	}

	updatedOrg, err := s.queries.UpdateOrganization(ctx, sqlc.UpdateOrganizationParams{
		ID:   database.PgUUID(org.ID),
		Name: org.Name,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to update organization")
		return nil, fmt.Errorf("failed to update organization: %w", err)
	}

	span.SetStatus(codes.Ok, "organization updated successfully")
	return toOrganization(updatedOrg)
}


// Organization Unit methods

func (s *OrganizationService) CreateUnit(ctx context.Context, orgID uuid.UUID, name string, alias string, address string) (*models.OrganizationUnit, error) {
	ctx, span := s.tracer.Start(ctx, "CreateUnit",
		trace.WithAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("org_unit.name", name),
			attribute.String("org_unit.alias", alias),
			attribute.String("org_unit.address", address),
		),
	)
	defer span.End()

	if err := validateName(name); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "name validation failed")
		return nil, err
	}
	if err := validateAlias(alias); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "alias validation failed")
		return nil, err
	}

	unit, err := s.queries.CreateOrgUnit(ctx, sqlc.CreateOrgUnitParams{
		OrgID:   database.PgUUID(orgID),
		Name:    name,
		Alias:   alias,
		Address: database.PgText(address),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create organization unit")
		return nil, fmt.Errorf("failed to create organization unit: %w", err)
	}

	span.SetStatus(codes.Ok, "organization unit created successfully")
	return toOrganizationUnit(unit)
}

func (s *OrganizationService) GetAllUnits(ctx context.Context, orgID uuid.UUID) ([]*models.OrganizationUnit, error) {
	ctx, span := s.tracer.Start(ctx, "GetAllUnits",
		trace.WithAttributes(
			attribute.String("org.id", orgID.String()),
		),
	)
	defer span.End()

	units, err := s.queries.GetOrgUnits(ctx, database.PgUUID(orgID))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get organization units")
		return nil, fmt.Errorf("failed to get organization units: %w", err)
	}

	result := make([]*models.OrganizationUnit, len(units))
	for i, unit := range units {
		result[i], err = toOrganizationUnit(unit)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed to convert organization unit")
			return nil, fmt.Errorf("failed to convert organization unit: %w", err)
		}
	}

	span.SetStatus(codes.Ok, "successfully retrieved organization units")
	return result, nil
}

func (s *OrganizationService) GetUnitByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*models.OrganizationUnit, error) {
	ctx, span := s.tracer.Start(ctx, "GetUnitByID",
		trace.WithAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("unit.id", id.String()),
		),
	)
	defer span.End()

	unit, err := s.queries.GetOrgUnitById(ctx, sqlc.GetOrgUnitByIdParams{
		OrgID: database.PgUUID(orgID),
		ID:    database.PgUUID(id),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get organization unit")
		return nil, fmt.Errorf("failed to get organization unit: %w", err)
	}

	span.SetStatus(codes.Ok, "organization unit retrieved successfully")
	return toOrganizationUnit(unit)
}

func (s *OrganizationService) DeleteUnit(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	ctx, span := s.tracer.Start(ctx, "DeleteUnit",
		trace.WithAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("unit.id", id.String()),
		),
	)
	defer span.End()

	err := s.queries.DeleteOrgUnit(ctx, sqlc.DeleteOrgUnitParams{
		OrgID: database.PgUUID(orgID),
		ID:    database.PgUUID(id),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to delete organization unit")
		return fmt.Errorf("failed to delete organization unit: %w", err)
	}

	span.SetStatus(codes.Ok, "organization unit deleted successfully")
	return nil
}

func (s *OrganizationService) UpdateUnit(ctx context.Context, unit *models.OrganizationUnit) (*models.OrganizationUnit, error) {
	ctx, span := s.tracer.Start(ctx, "UpdateUnit",
		trace.WithAttributes(
			attribute.String("org.id", unit.OrgID.String()),
			attribute.String("unit.id", unit.ID.String()),
		),
	)
	defer span.End()

	if err := validateName(unit.Name); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "name validation failed")
		return nil, err
	}
	if err := validateAlias(unit.Alias); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "alias validation failed")
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
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to update organization unit")
		return nil, fmt.Errorf("failed to update organization unit: %w", err)
	}

	span.SetStatus(codes.Ok, "organization unit updated successfully")
	return toOrganizationUnit(updatedUnit)
}
