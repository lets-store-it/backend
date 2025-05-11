package organization

import (
	"github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/database"
	"github.com/let-store-it/backend/internal/models"
)

func toOrganizationModel(org sqlc.Org) *models.Organization {
	return &models.Organization{
		ID:        database.UUIDFromPgx(org.ID),
		Name:      org.Name,
		Subdomain: org.Subdomain,
	}
}

func toOrganizationUnitModel(unit sqlc.OrgUnit) *models.OrganizationUnit {
	return &models.OrganizationUnit{
		ID:        database.UUIDFromPgx(unit.ID),
		OrgID:     database.UUIDFromPgx(unit.OrgID),
		Name:      unit.Name,
		Alias:     unit.Alias,
		Address:   database.PgTextPtrFromPgx(unit.Address),
		CreatedAt: unit.CreatedAt.Time,
		DeletedAt: database.PgTimePtrFromPgx(unit.DeletedAt),
	}
}
