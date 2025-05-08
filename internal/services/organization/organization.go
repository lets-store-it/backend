package organization

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	database "github.com/let-store-it/backend/generated/sqlc"
	db "github.com/let-store-it/backend/internal/database"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrOrganizationNotFound               = errors.New("organization not found")
	ErrOrganizationUnitNotFound           = errors.New("organization unit not found")
	ErrInvalidOrganization                = errors.New("invalid organization")
	ErrInvalidOrganizationUnit            = errors.New("invalid organization unit")
	ErrInvalidUserID                      = errors.New("invalid user ID")
	ErrOrganizationSubdomainAlreadyExists = errors.New("organization already exists")

	ErrValidationError = errors.New("validation error")
)

const (
	maxNameLength      = 100
	maxSubdomainLength = 63
	maxAliasLength     = 100
)

func validateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%w: name cannot be empty", ErrValidationError)
	}
	if len(name) > maxNameLength {
		return fmt.Errorf("%w: name is too long", ErrValidationError)
	}
	return nil
}

func validateSubdomain(subdomain string) error {
	if strings.TrimSpace(subdomain) == "" {
		return fmt.Errorf("%w: subdomain cannot be empty", ErrValidationError)
	}
	if len(subdomain) > maxSubdomainLength {
		return fmt.Errorf("%w: subdomain is too long", ErrValidationError)
	}
	matched, _ := regexp.MatchString("^[a-z0-9-]+$", subdomain)
	if !matched {
		return fmt.Errorf("%w: subdomain can only contain lowercase letters, numbers, and hyphens", ErrValidationError)
	}
	return nil
}

func validateAlias(alias string) error {
	if strings.TrimSpace(alias) == "" {
		return fmt.Errorf("%w: alias cannot be empty", ErrValidationError)
	}
	if len(alias) > maxAliasLength {
		return fmt.Errorf("%w: alias is too long", ErrValidationError)
	}
	matched, _ := regexp.MatchString("^[\\w-]+$", alias)
	if !matched {
		return fmt.Errorf("%w: alias can only contain letters, numbers, and hyphens (no spaces)", ErrValidationError)
	}
	return nil
}

type OrganizationService struct {
	queries *database.Queries
	pgxPool *pgxpool.Pool
	tracer  trace.Tracer
}

func New(queries *database.Queries, pgxPool *pgxpool.Pool) *OrganizationService {
	return &OrganizationService{
		queries: queries,
		pgxPool: pgxPool,
		tracer:  otel.GetTracerProvider().Tracer("organization-service"),
	}
}

func (s *OrganizationService) Create(ctx context.Context, name string, subdomain string) (*models.Organization, error) {
	ctx, span := s.tracer.Start(ctx, "Create",
		trace.WithAttributes(
			attribute.String("name", name),
			attribute.String("subdomain", subdomain),
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

	org, err := s.queries.CreateOrg(ctx, database.CreateOrgParams{
		Name:      name,
		Subdomain: subdomain,
	})
	if err != nil {
		if db.IsUniqueViolation(err) {
			span.RecordError(err)
			span.SetStatus(codes.Error, "organization already exists")
			return nil, ErrOrganizationSubdomainAlreadyExists
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

	if userID == uuid.Nil {
		span.SetStatus(codes.Error, "invalid user ID")
		return nil, ErrInvalidUserID
	}

	res, err := s.queries.GetUserOrgs(ctx, utils.PgUUID(userID))
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

	if id == uuid.Nil {
		span.SetStatus(codes.Error, "invalid organization ID")
		return nil, ErrInvalidOrganization
	}

	org, err := s.queries.GetOrg(ctx, utils.PgUUID(id))
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

	if id == uuid.Nil {
		span.SetStatus(codes.Error, "invalid organization ID")
		return ErrInvalidOrganization
	}

	err := s.queries.DeleteOrg(ctx, utils.PgUUID(id))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to delete organization")
		return fmt.Errorf("failed to delete organization: %w", err)
	}

	span.SetStatus(codes.Ok, "organization deleted successfully")
	return nil
}

func (s *OrganizationService) Update(ctx context.Context, org *models.Organization) (*models.Organization, error) {
	ctx, span := s.tracer.Start(ctx, "Update")
	defer span.End()

	if org == nil || org.ID == uuid.Nil {
		span.SetStatus(codes.Error, "invalid organization")
		return nil, ErrInvalidOrganization
	}

	span.SetAttributes(attribute.String("org.id", org.ID.String()))

	exists, err := s.IsOrganizationExists(ctx, org.ID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to check organization existence")
		return nil, fmt.Errorf("failed to check organization existence: %w", err)
	}
	if !exists {
		span.SetStatus(codes.Error, "organization not found")
		return nil, ErrOrganizationNotFound
	}

	if err := validateName(org.Name); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "name validation failed")
		return nil, err
	}

	res, err := s.queries.UpdateOrg(ctx, database.UpdateOrgParams{
		ID:   utils.PgUUID(org.ID),
		Name: org.Name,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to update organization")
		return nil, fmt.Errorf("failed to update organization: %w", err)
	}

	span.SetStatus(codes.Ok, "organization updated successfully")
	return toOrganization(res)
}

func (s *OrganizationService) IsOrganizationExists(ctx context.Context, id uuid.UUID) (bool, error) {
	ctx, span := s.tracer.Start(ctx, "IsOrganizationExists",
		trace.WithAttributes(
			attribute.String("org.id", id.String()),
		),
	)
	defer span.End()

	if id == uuid.Nil {
		span.SetStatus(codes.Error, "invalid organization ID")
		return false, ErrInvalidOrganization
	}

	exists, err := s.queries.IsOrgExists(ctx, utils.PgUUID(id))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to check organization existence")
		return false, fmt.Errorf("failed to check organization existence: %w", err)
	}

	span.SetStatus(codes.Ok, "successfully checked organization existence")
	return exists, nil
}

// Organization Unit methods
func (s *OrganizationService) CreateUnit(ctx context.Context, orgID uuid.UUID, name string, alias string, address string) (*models.OrganizationUnit, error) {
	ctx, span := s.tracer.Start(ctx, "CreateUnit",
		trace.WithAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("name", name),
			attribute.String("alias", alias),
		),
	)
	defer span.End()

	if orgID == uuid.Nil {
		span.SetStatus(codes.Error, "invalid organization ID")
		return nil, ErrInvalidOrganization
	}
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

	unit, err := s.queries.CreateOrgUnit(ctx, database.CreateOrgUnitParams{
		OrgID:   utils.PgUUID(orgID),
		Name:    name,
		Alias:   alias,
		Address: utils.PgText(address),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create organization unit")
		return nil, fmt.Errorf("failed to create organization unit: %w", err)
	}

	span.SetAttributes(attribute.String("unit.id", unit.ID.String()))
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

	if orgID == uuid.Nil {
		span.SetStatus(codes.Error, "invalid organization ID")
		return nil, ErrInvalidOrganization
	}

	units, err := s.queries.GetOrgUnits(ctx, utils.PgUUID(orgID))
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

	if orgID == uuid.Nil {
		span.SetStatus(codes.Error, "invalid organization ID")
		return nil, ErrInvalidOrganization
	}
	if id == uuid.Nil {
		span.SetStatus(codes.Error, "invalid unit ID")
		return nil, ErrInvalidOrganizationUnit
	}

	unit, err := s.queries.GetOrgUnit(ctx, database.GetOrgUnitParams{
		OrgID: utils.PgUUID(orgID),
		ID:    utils.PgUUID(id),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get organization unit")
		return nil, fmt.Errorf("failed to get organization unit: %w", err)
	}

	unitModel, err := toOrganizationUnit(unit)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to convert organization unit")
		return nil, err
	}
	if unitModel == nil {
		span.SetStatus(codes.Error, "organization unit not found")
		return nil, ErrOrganizationUnitNotFound
	}

	span.SetStatus(codes.Ok, "successfully retrieved organization unit")
	return unitModel, nil
}

func (s *OrganizationService) DeleteUnit(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	ctx, span := s.tracer.Start(ctx, "DeleteUnit",
		trace.WithAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("unit.id", id.String()),
		),
	)
	defer span.End()

	if orgID == uuid.Nil {
		span.SetStatus(codes.Error, "invalid organization ID")
		return ErrInvalidOrganization
	}
	if id == uuid.Nil {
		span.SetStatus(codes.Error, "invalid unit ID")
		return ErrInvalidOrganizationUnit
	}

	err := s.queries.DeleteOrgUnit(ctx, database.DeleteOrgUnitParams{
		OrgID: utils.PgUUID(orgID),
		ID:    utils.PgUUID(id),
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

	if unit == nil || unit.ID == uuid.Nil || unit.OrgID == uuid.Nil {
		span.SetStatus(codes.Error, "invalid organization unit")
		return nil, ErrInvalidOrganizationUnit
	}

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

	updatedUnit, err := s.queries.UpdateOrgUnit(ctx, database.UpdateOrgUnitParams{
		ID:      utils.PgUUID(unit.ID),
		OrgID:   utils.PgUUID(unit.OrgID),
		Name:    unit.Name,
		Alias:   unit.Alias,
		Address: utils.PgText(address),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to update organization unit")
		return nil, fmt.Errorf("failed to update organization unit: %w", err)
	}

	span.SetStatus(codes.Ok, "organization unit updated successfully")
	return toOrganizationUnit(updatedUnit)
}
