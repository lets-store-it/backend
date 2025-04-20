package organization

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/let-store-it/backend/generated/database"
	"github.com/let-store-it/backend/internal/models"
)

var (
	ErrOrganizationNotFound     = errors.New("organization not found")
	ErrOrganizationUnitNotFound = errors.New("organization unit not found")
	ErrInvalidOrganization      = errors.New("invalid organization")
	ErrInvalidOrganizationUnit  = errors.New("invalid organization unit")
	ErrInvalidName              = errors.New("invalid name")
	ErrInvalidSubdomain         = errors.New("invalid subdomain")
	ErrInvalidAlias             = errors.New("invalid alias")
)

func uuidFromPgx(id pgtype.UUID) *uuid.UUID {
	if !id.Valid {
		return nil
	}
	result := uuid.UUID(id.Bytes)
	return &result
}

func toOrganization(org database.Org) (*models.Organization, error) {
	id := uuidFromPgx(org.ID)
	if id == nil {
		return nil, fmt.Errorf("failed to convert organization: %w", ErrInvalidOrganization)
	}
	return &models.Organization{
		ID:        *id,
		Name:      org.Name,
		Subdomain: org.Subdomain,
	}, nil
}

type OrganizationService struct {
	queries *database.Queries
	pgxPool *pgxpool.Pool
}

func New(queries *database.Queries, pgxPool *pgxpool.Pool) *OrganizationService {
	return &OrganizationService{
		queries: queries,
		pgxPool: pgxPool,
	}
}

func validateOrganizationData(name, subdomain string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%w: name cannot be empty", ErrInvalidName)
	}
	if len(name) > 100 {
		return fmt.Errorf("%w: name is too long (max 100 characters)", ErrInvalidName)
	}
	if strings.TrimSpace(subdomain) == "" {
		return fmt.Errorf("%w: subdomain cannot be empty", ErrInvalidSubdomain)
	}
	if len(subdomain) > 63 {
		return fmt.Errorf("%w: subdomain is too long (max 63 characters)", ErrInvalidSubdomain)
	}
	matched, _ := regexp.MatchString("^[a-z0-9-]+$", subdomain)
	if !matched {
		return fmt.Errorf("%w: subdomain can only contain lowercase letters, numbers, and hyphens", ErrInvalidSubdomain)
	}
	return nil
}

func (s *OrganizationService) Create(ctx context.Context, name string, subdomain string) (*models.Organization, error) {
	if err := validateOrganizationData(name, subdomain); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	org, err := s.queries.CreateOrg(ctx, database.CreateOrgParams{
		Name:      name,
		Subdomain: subdomain,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	return toOrganization(org)
}

func (s *OrganizationService) GetUsersOrgs(ctx context.Context, userID uuid.UUID) ([]*models.Organization, error) {
	if userID == uuid.Nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	res, err := s.queries.GetUserOrgs(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get user organizations: %w", err)
	}

	orgs := make([]*models.Organization, len(res))
	for i, org := range res {
		orgs[i], err = toOrganization(org)
		if err != nil {
			return nil, fmt.Errorf("failed to convert organization: %w", err)
		}
	}

	return orgs, nil
}

func (s *OrganizationService) GetByID(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	if id == uuid.Nil {
		return nil, ErrInvalidOrganization
	}

	org, err := s.queries.GetOrg(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}
	return toOrganization(org)
}

func (s *OrganizationService) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return ErrInvalidOrganization
	}

	err := s.queries.DeleteOrg(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return fmt.Errorf("failed to delete organization: %w", err)
	}
	return nil
}

func (s *OrganizationService) Update(ctx context.Context, org *models.Organization) (*models.Organization, error) {
	if org == nil {
		return nil, ErrInvalidOrganization
	}
	if org.ID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}

	exists, err := s.IsOrganizationExists(ctx, org.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check organization existence: %w", err)
	}
	if !exists {
		return nil, ErrOrganizationNotFound
	}

	if err := validateOrganizationData(org.Name, org.Subdomain); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	res, err := s.queries.UpdateOrg(ctx, database.UpdateOrgParams{
		ID:        pgtype.UUID{Bytes: org.ID, Valid: true},
		Name:      org.Name,
		Subdomain: org.Subdomain,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update organization: %w", err)
	}
	return toOrganization(res)
}

func (s *OrganizationService) IsOrganizationExists(ctx context.Context, id uuid.UUID) (bool, error) {
	if id == uuid.Nil {
		return false, ErrInvalidOrganization
	}

	exists, err := s.queries.IsOrgExists(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return false, fmt.Errorf("failed to check organization existence: %w", err)
	}
	return exists, nil
}

func toOrganizationUnit(unit database.OrgUnit) (*models.OrganizationUnit, error) {
	id := uuidFromPgx(unit.ID)
	if id == nil {
		return nil, fmt.Errorf("failed to convert organization unit: %w", ErrInvalidOrganizationUnit)
	}
	orgID := uuidFromPgx(unit.OrgID)
	if orgID == nil {
		return nil, fmt.Errorf("failed to convert organization unit: %w", ErrInvalidOrganization)
	}

	var address *string
	if unit.Address.Valid {
		address = &unit.Address.String
	}

	var deletedAt *time.Time
	if unit.DeletedAt.Valid {
		deletedAt = &unit.DeletedAt.Time
	}

	return &models.OrganizationUnit{
		ID:        *id,
		OrgID:     *orgID,
		Name:      unit.Name,
		Alias:     unit.Alias,
		Address:   address,
		CreatedAt: unit.CreatedAt.Time,
		DeletedAt: deletedAt,
	}, nil
}

func validateOrganizationUnitData(name string, alias string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%w: name cannot be empty", ErrInvalidName)
	}
	if len(name) > 100 {
		return fmt.Errorf("%w: name is too long (max 100 characters)", ErrInvalidName)
	}

	if strings.TrimSpace(alias) == "" {
		return fmt.Errorf("%w: alias cannot be empty", ErrInvalidAlias)
	}
	if len(alias) > 100 {
		return fmt.Errorf("%w: alias is too long (max 100 characters)", ErrInvalidAlias)
	}
	matched, _ := regexp.MatchString("^[\\w-]+$", alias)
	if !matched {
		return fmt.Errorf("%w: alias can only contain letters, numbers, and hyphens (no spaces)", ErrInvalidAlias)
	}
	return nil
}

func (s *OrganizationService) CreateUnit(ctx context.Context, orgID uuid.UUID, name string, alias string, address string) (*models.OrganizationUnit, error) {
	if orgID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}
	if err := validateOrganizationUnitData(name, alias); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	unit, err := s.queries.CreateOrgUnit(ctx, database.CreateOrgUnitParams{
		OrgID:   pgtype.UUID{Bytes: orgID, Valid: true},
		Name:    name,
		Alias:   alias,
		Address: pgtype.Text{String: address, Valid: address != ""},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create organization unit: %w", err)
	}

	return toOrganizationUnit(unit)
}

func (s *OrganizationService) GetAllUnits(ctx context.Context, orgID uuid.UUID) ([]*models.OrganizationUnit, error) {
	if orgID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}

	units, err := s.queries.GetOrgUnits(ctx, pgtype.UUID{Bytes: orgID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get organization units: %w", err)
	}

	result := make([]*models.OrganizationUnit, len(units))
	for i, unit := range units {
		result[i], err = toOrganizationUnit(unit)
		if err != nil {
			return nil, fmt.Errorf("failed to convert organization unit: %w", err)
		}
	}

	return result, nil
}

func (s *OrganizationService) GetUnitByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*models.OrganizationUnit, error) {
	if orgID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}
	if id == uuid.Nil {
		return nil, ErrInvalidOrganizationUnit
	}

	unit, err := s.queries.GetOrgUnit(ctx, database.GetOrgUnitParams{
		OrgID: pgtype.UUID{Bytes: orgID, Valid: true},
		ID:    pgtype.UUID{Bytes: id, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get organization unit: %w", err)
	}

	unitModel, err := toOrganizationUnit(unit)
	if err != nil {
		return nil, err
	}
	if unitModel == nil {
		return nil, ErrOrganizationUnitNotFound
	}
	return unitModel, nil
}

func (s *OrganizationService) DeleteUnit(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if orgID == uuid.Nil {
		return ErrInvalidOrganization
	}
	if id == uuid.Nil {
		return ErrInvalidOrganizationUnit
	}

	err := s.queries.DeleteOrgUnit(ctx, database.DeleteOrgUnitParams{
		OrgID: pgtype.UUID{Bytes: orgID, Valid: true},
		ID:    pgtype.UUID{Bytes: id, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to delete organization unit: %w", err)
	}
	return nil
}

func (s *OrganizationService) UpdateUnit(ctx context.Context, unit *models.OrganizationUnit) (*models.OrganizationUnit, error) {
	if unit == nil {
		return nil, ErrInvalidOrganizationUnit
	}
	if unit.ID == uuid.Nil {
		return nil, ErrInvalidOrganizationUnit
	}
	if unit.OrgID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}

	if err := validateOrganizationUnitData(unit.Name, unit.Alias); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	var address string
	if unit.Address != nil {
		address = *unit.Address
	}

	updatedUnit, err := s.queries.UpdateOrgUnit(ctx, database.UpdateOrgUnitParams{
		ID:      pgtype.UUID{Bytes: unit.ID, Valid: true},
		Name:    unit.Name,
		Alias:   unit.Alias,
		Address: pgtype.Text{String: address, Valid: address != ""},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update organization unit: %w", err)
	}
	return toOrganizationUnit(updatedUnit)
}
