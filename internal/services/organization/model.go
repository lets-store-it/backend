package organization

import (
	"github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/database"
	"github.com/let-store-it/backend/internal/models"
)

func toOrganization(org sqlc.Org) (*models.Organization, error) {
	return &models.Organization{
		ID:        database.UuidFromPgx(org.ID),
		Name:      org.Name,
		Subdomain: org.Subdomain,
	}, nil
}

func toOrganizationUnit(unit sqlc.OrgUnit) (*models.OrganizationUnit, error) {
	return &models.OrganizationUnit{
		ID:        database.UuidFromPgx(unit.ID),
		OrgID:     database.UuidFromPgx(unit.OrgID),
		Name:      unit.Name,
		Alias:     unit.Alias,
		Address:   database.PgTextPtrFromPgx(unit.Address),
		CreatedAt: unit.CreatedAt.Time,
		DeletedAt: database.PgTimePtrFromPgx(unit.DeletedAt),
	}, nil
}
