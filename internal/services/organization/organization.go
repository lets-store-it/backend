package organization

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/let-store-it/backend/generated/database"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/utils"
)

var (
	ErrOrganizationNotFound     = errors.New("organization not found")
	ErrOrganizationUnitNotFound = errors.New("organization unit not found")
	ErrInvalidOrganization      = errors.New("invalid organization")
	ErrInvalidOrganizationUnit  = errors.New("invalid organization unit")
	ErrInvalidUserID            = errors.New("invalid user ID")

	ErrValidationError = errors.New("validation error")
)

const (
	maxNameLength      = 100
	maxSubdomainLength = 63
	maxAliasLength     = 100
)

func validateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.Join(ErrValidationError, errors.New("name cannot be empty"))
	}
	if len(name) > maxNameLength {
		return errors.Join(ErrValidationError, errors.New("name is too long"))
	}
	return nil
}

func validateSubdomain(subdomain string) error {
	if strings.TrimSpace(subdomain) == "" {
		return errors.Join(ErrValidationError, errors.New("subdomain cannot be empty"))
	}
	if len(subdomain) > maxSubdomainLength {
		return errors.Join(ErrValidationError, errors.New("subdomain is too long"))
	}
	matched, _ := regexp.MatchString("^[a-z0-9-]+$", subdomain)
	if !matched {
		return errors.Join(ErrValidationError, errors.New("subdomain can only contain lowercase letters, numbers, and hyphens"))
	}
	return nil
}

func validateAlias(alias string) error {
	if strings.TrimSpace(alias) == "" {
		return errors.Join(ErrValidationError, errors.New("alias cannot be empty"))
	}
	if len(alias) > maxAliasLength {
		return errors.Join(ErrValidationError, errors.New("alias is too long"))
	}
	matched, _ := regexp.MatchString("^[\\w-]+$", alias)
	if !matched {
		return errors.Join(ErrValidationError, errors.New("alias can only contain letters, numbers, and hyphens (no spaces)"))
	}
	return nil
}

func toOrganization(org database.Org) (*models.Organization, error) {
	id := utils.UuidFromPgx(org.ID)
	if id == nil {
		return nil, fmt.Errorf("failed to convert organization: %w", ErrInvalidOrganization)
	}
	return &models.Organization{
		ID:        *id,
		Name:      org.Name,
		Subdomain: org.Subdomain,
	}, nil
}

func toOrganizationUnit(unit database.OrgUnit) (*models.OrganizationUnit, error) {
	id := utils.UuidFromPgx(unit.ID)
	if id == nil {
		return nil, fmt.Errorf("failed to convert organization unit: %w", ErrInvalidOrganizationUnit)
	}
	orgID := utils.UuidFromPgx(unit.OrgID)
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

func (s *OrganizationService) Create(ctx context.Context, name string, subdomain string) (*models.Organization, error) {
	if err := validateName(name); err != nil {
		return nil, err
	}
	if err := validateSubdomain(subdomain); err != nil {
		return nil, err
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
		return nil, ErrInvalidUserID
	}

	res, err := s.queries.GetUserOrgs(ctx, utils.PgUUID(userID))
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

	org, err := s.queries.GetOrg(ctx, utils.PgUUID(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}
	return toOrganization(org)
}

func (s *OrganizationService) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return ErrInvalidOrganization
	}

	err := s.queries.DeleteOrg(ctx, utils.PgUUID(id))
	if err != nil {
		return fmt.Errorf("failed to delete organization: %w", err)
	}
	return nil
}

func (s *OrganizationService) Update(ctx context.Context, org *models.Organization) (*models.Organization, error) {
	if org == nil || org.ID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}

	exists, err := s.IsOrganizationExists(ctx, org.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check organization existence: %w", err)
	}
	if !exists {
		return nil, ErrOrganizationNotFound
	}

	if err := validateName(org.Name); err != nil {
		return nil, err
	}
	if err := validateSubdomain(org.Subdomain); err != nil {
		return nil, err
	}

	res, err := s.queries.UpdateOrg(ctx, database.UpdateOrgParams{
		ID:        utils.PgUUID(org.ID),
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

	exists, err := s.queries.IsOrgExists(ctx, utils.PgUUID(id))
	if err != nil {
		return false, fmt.Errorf("failed to check organization existence: %w", err)
	}
	return exists, nil
}

// Organization Unit methods
func (s *OrganizationService) CreateUnit(ctx context.Context, orgID uuid.UUID, name string, alias string, address string) (*models.OrganizationUnit, error) {
	if orgID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}
	if err := validateName(name); err != nil {
		return nil, err
	}
	if err := validateAlias(alias); err != nil {
		return nil, err
	}

	unit, err := s.queries.CreateOrgUnit(ctx, database.CreateOrgUnitParams{
		OrgID:   utils.PgUUID(orgID),
		Name:    name,
		Alias:   alias,
		Address: utils.PgText(address),
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

	units, err := s.queries.GetOrgUnits(ctx, utils.PgUUID(orgID))
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
		OrgID: utils.PgUUID(orgID),
		ID:    utils.PgUUID(id),
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
		OrgID: utils.PgUUID(orgID),
		ID:    utils.PgUUID(id),
	})
	if err != nil {
		return fmt.Errorf("failed to delete organization unit: %w", err)
	}
	return nil
}

func (s *OrganizationService) UpdateUnit(ctx context.Context, unit *models.OrganizationUnit) (*models.OrganizationUnit, error) {
	if unit == nil || unit.ID == uuid.Nil || unit.OrgID == uuid.Nil {
		return nil, ErrInvalidOrganizationUnit
	}

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

	updatedUnit, err := s.queries.UpdateOrgUnit(ctx, database.UpdateOrgUnitParams{
		ID:      utils.PgUUID(unit.ID),
		Name:    unit.Name,
		Alias:   unit.Alias,
		Address: utils.PgText(address),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update organization unit: %w", err)
	}
	return toOrganizationUnit(updatedUnit)
}
