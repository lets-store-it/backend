package organization

import (
	"fmt"
	"time"

	"github.com/let-store-it/backend/generated/database"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/utils"
)

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
