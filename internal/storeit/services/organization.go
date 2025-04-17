package services

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
	"github.com/let-store-it/backend/internal/storeit/models"
)

var (
	ErrOrganizationNotFound = errors.New("organization not found")
)

func uuidFromPgx(id pgtype.UUID) *uuid.UUID {
	if !id.Valid {
		return nil
	}
	uuid := (uuid.UUID)(id.Bytes)
	return &uuid
}

func toOrganization(org database.Org) (*models.Organization, error) {
	id := uuidFromPgx(org.ID)
	if id == nil {
		return nil, errors.New("id is nil")
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

func NewOrganizationService(queries *database.Queries, pgxPool *pgxpool.Pool) *OrganizationService {
	return &OrganizationService{
		queries: queries,
		pgxPool: pgxPool,
	}
}

func validateOrganizationData(name, subdomain string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("organization name cannot be empty")
	}
	if len(name) > 100 {
		return fmt.Errorf("organization name is too long (max 100 characters)")
	}
	if strings.TrimSpace(subdomain) == "" {
		return fmt.Errorf("subdomain cannot be empty")
	}
	if len(subdomain) > 63 {
		return fmt.Errorf("subdomain is too long (max 63 characters)")
	}
	matched, _ := regexp.MatchString("^[a-z0-9-]+$", subdomain)
	if !matched {
		return fmt.Errorf("subdomain can only contain lowercase letters, numbers, and hyphens")
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
		return nil, err
	}

	return toOrganization(org)
}

func (s *OrganizationService) GetUsersOrgs(ctx context.Context, userID uuid.UUID) ([]*models.Organization, error) {
	res, err := s.queries.GetUserOrgs(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, err
	}

	orgs := make([]*models.Organization, len(res))
	for i, org := range res {
		orgs[i], err = toOrganization(org)
		if err != nil {
			return nil, err
		}
	}

	return orgs, nil
}

func (s *OrganizationService) GetAll(ctx context.Context) ([]*models.Organization, error) {
	res, err := s.queries.GetActiveOrgs(ctx)
	if err != nil {
		return nil, err
	}

	orgs := make([]*models.Organization, len(res))
	for i, org := range res {
		orgs[i], err = toOrganization(org)
		if err != nil {
			return nil, err
		}
	}

	return orgs, nil
}

func (s *OrganizationService) GetByID(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	org, err := s.queries.GetOrg(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}
	return toOrganization(org)
}

func (s *OrganizationService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.queries.DeleteOrg(ctx, pgtype.UUID{Bytes: id, Valid: true})
	return err
}

func (s *OrganizationService) Update(ctx context.Context, org *models.Organization) (*models.Organization, error) {
	exists, err := s.IsOrganizationExists(ctx, org.ID)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	return toOrganization(res)
}

func (s *OrganizationService) IsOrganizationExists(ctx context.Context, id uuid.UUID) (bool, error) {
	exists, err := s.queries.IsOrgExists(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return false, err
	}
	return exists, nil
}

var (
	ErrOrganizationUnitNotFound = errors.New("organization unit not found")
)

func toOrganizationUnit(unit database.OrgUnit) (*models.OrganizationUnit, error) {
	id := uuidFromPgx(unit.ID)
	if id == nil {
		return nil, errors.New("id is nil")
	}
	orgID := uuidFromPgx(unit.OrgID)
	if orgID == nil {
		return nil, errors.New("org_id is nil")
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
		return fmt.Errorf("organization unit name cannot be empty")
	}
	if len(name) > 100 {
		return fmt.Errorf("organization unit name is too long (max 100 characters)")
	}

	if strings.TrimSpace(alias) == "" {
		return fmt.Errorf("organization unit alias cannot be empty")
	}
	if len(alias) > 100 {
		return fmt.Errorf("organization unit alias is too long (max 100 characters)")
	}
	matched, _ := regexp.MatchString("^[\\w-]+$", alias)
	if !matched {
		return fmt.Errorf("organization unit alias can only contain letters, numbers, and hyphens (no spaces)")
	}
	return nil
}

func (s *OrganizationService) CreateUnit(ctx context.Context, orgID uuid.UUID, name string, alias string, address string) (*models.OrganizationUnit, error) {
	unit, err := s.queries.CreateOrgUnit(ctx, database.CreateOrgUnitParams{
		OrgID:   pgtype.UUID{Bytes: orgID, Valid: true},
		Name:    name,
		Alias:   alias,
		Address: pgtype.Text{String: address, Valid: address != ""},
	})
	if err != nil {
		return nil, err
	}

	return toOrganizationUnit(unit)
}

func (s *OrganizationService) GetAllUnits(ctx context.Context, orgID uuid.UUID) ([]*models.OrganizationUnit, error) {
	units, err := s.queries.GetActiveOrgUnits(ctx, pgtype.UUID{Bytes: orgID, Valid: true})
	if err != nil {
		return nil, err
	}

	result := make([]*models.OrganizationUnit, len(units))
	for i, unit := range units {
		result[i], err = toOrganizationUnit(unit)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (s *OrganizationService) GetUnitByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*models.OrganizationUnit, error) {
	unit, err := s.queries.GetOrgUnit(ctx, database.GetOrgUnitParams{
		OrgID: pgtype.UUID{Bytes: orgID, Valid: true},
		ID:    pgtype.UUID{Bytes: id, Valid: true},
	})
	if err != nil {
		return nil, err
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
	return s.queries.DeleteOrgUnit(ctx, database.DeleteOrgUnitParams{
		OrgID: pgtype.UUID{Bytes: orgID, Valid: true},
		ID:    pgtype.UUID{Bytes: id, Valid: true},
	})
}

func (s *OrganizationService) UpdateUnit(ctx context.Context, unit *models.OrganizationUnit) (*models.OrganizationUnit, error) {
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
		return nil, err
	}
	return toOrganizationUnit(updatedUnit)
}

func (s *OrganizationService) IsOrganizationUnitExists(ctx context.Context, orgID uuid.UUID, unitID uuid.UUID) (bool, error) {
	return s.queries.IsOrgUnitExists(ctx, database.IsOrgUnitExistsParams{
		OrgID: pgtype.UUID{Bytes: orgID, Valid: true},
		ID:    pgtype.UUID{Bytes: unitID, Valid: true},
	})
}
