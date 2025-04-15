package repositories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/let-store-it/backend/generated/database"
	"github.com/let-store-it/backend/internal/storeit/models"
)

type OrganizationUnitRepository struct {
	Queries *database.Queries
}

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
	return &models.OrganizationUnit{
		ID:      *id,
		OrgID:   *orgID,
		Name:    unit.Name,
		Alias:   unit.Alias,
		Address: address,
	}, nil
}

func (r *OrganizationUnitRepository) GetOrganizationUnitByID(ctx context.Context, id uuid.UUID) (*models.OrganizationUnit, error) {
	unit, err := r.Queries.GetOrganizationUnitById(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}
	return toOrganizationUnit(unit)
}

func (r *OrganizationUnitRepository) GetOrganizationUnits(ctx context.Context, orgID uuid.UUID) ([]*models.OrganizationUnit, error) {
	units, err := r.Queries.GetOrganizationUnits(ctx, pgtype.UUID{Bytes: orgID, Valid: true})
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

func (r *OrganizationUnitRepository) CreateOrganizationUnit(ctx context.Context, orgID uuid.UUID, name string, alias string, address string) (*models.OrganizationUnit, error) {
	unit, err := r.Queries.CreateOrganizationUnit(ctx, database.CreateOrganizationUnitParams{
		OrgID:   pgtype.UUID{Bytes: orgID, Valid: true},
		Name:    name,
		Alias:   alias,
		Address: pgtype.Text{String: address, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return toOrganizationUnit(unit)
}

func (r *OrganizationUnitRepository) DeleteOrganizationUnit(ctx context.Context, id uuid.UUID) error {
	return r.Queries.DeleteOrganizationUnit(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func (r *OrganizationUnitRepository) UpdateOrganizationUnit(ctx context.Context, unit *models.OrganizationUnit) (*models.OrganizationUnit, error) {
	var address string
	if unit.Address != nil {
		address = *unit.Address
	}

	updatedUnit, err := r.Queries.UpdateOrganizationUnit(ctx, database.UpdateOrganizationUnitParams{
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

func (r *OrganizationUnitRepository) IsOrganizationUnitExistsForOrganization(ctx context.Context, orgID uuid.UUID, unitID uuid.UUID) (bool, error) {
	return r.Queries.IsOrganizationUnitExistsForOrganization(ctx, database.IsOrganizationUnitExistsForOrganizationParams{
		OrgID: pgtype.UUID{Bytes: orgID, Valid: true},
		ID:    pgtype.UUID{Bytes: unitID, Valid: true},
	})
}
