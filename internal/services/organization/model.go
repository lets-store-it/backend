package organization

import (
	"fmt"
	"time"

	"github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/database"
	"github.com/let-store-it/backend/internal/models"
)

func toOrganization(org sqlc.Org) (*models.Organization, error) {
	id := database.UuidPtrFromPgx(org.ID)
	if id == nil {
		return nil, fmt.Errorf("failed to convert organization: %w", ErrInvalidOrganization)
	}
	return &models.Organization{
		ID:        *id,
		Name:      org.Name,
		Subdomain: org.Subdomain,
	}, nil
}

func toOrganizationUnit(unit sqlc.OrgUnit) (*models.OrganizationUnit, error) {
	id := database.UuidPtrFromPgx(unit.ID)
	if id == nil {
		return nil, fmt.Errorf("failed to convert organization unit: %w", ErrInvalidOrganizationUnit)
	}
	orgID := database.UuidPtrFromPgx(unit.OrgID)
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
