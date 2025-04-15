package repositories

import (
	"context"
	"errors"

	"github.com/let-store-it/backend/generated/database"
	"github.com/let-store-it/backend/internal/storeit/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type OrganizationRepository struct {
	Queries *database.Queries
}

func uuidFromPgx(id pgtype.UUID) *uuid.UUID {
	if !id.Valid {
		return nil
	}
	uuid := (uuid.UUID)(id.Bytes)
	return &uuid
}

// func datetimeFromPgx(datetime pgtype.Timestamp) *time.Time {
// 	if !datetime.Valid {
// 		return nil
// 	}
// 	return &datetime.Time
// }

// Organization errors
// Organization Exists Error
type OrganizationExistsError struct {
}

func (e *OrganizationExistsError) Error() string {
	return "organization already exists"
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

func (r *OrganizationRepository) GetOrganizationByID(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	org, err := r.Queries.GetOrgById(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}
	return toOrganization(org)
}

func (r *OrganizationRepository) GetOrganizations(ctx context.Context) ([]*models.Organization, error) {
	res, err := r.Queries.GetOrgs(ctx)
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

func (r *OrganizationRepository) CreateOrganization(ctx context.Context, name string, subdomain string) (*models.Organization, error) {
	org, err := r.Queries.CreateOrg(ctx, database.CreateOrgParams{
		Name:      name,
		Subdomain: subdomain,
	})
	if err != nil {
		return nil, err
	}

	return toOrganization(org)
}

func (r *OrganizationRepository) DeleteOrganization(ctx context.Context, id uuid.UUID) error {
	err := r.Queries.DeleteOrg(ctx, pgtype.UUID{Bytes: id, Valid: true})
	return err
}

func (r *OrganizationRepository) IsOrganizationExists(ctx context.Context, name string, subdomain string) (bool, error) {
	exists, err := r.Queries.IsOrgExists(ctx, database.IsOrgExistsParams{
		Name:      name,
		Subdomain: subdomain,
	})
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *OrganizationRepository) IsOrganizationExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	exists, err := r.Queries.IsOrgExistsById(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *OrganizationRepository) UpdateOrganization(ctx context.Context, org *models.Organization) (*models.Organization, error) {
	res, err := r.Queries.UpdateOrg(ctx, database.UpdateOrgParams{
		ID:        pgtype.UUID{Bytes: org.ID, Valid: true},
		Name:      org.Name,
		Subdomain: org.Subdomain,
	})
	if err != nil {
		return nil, err
	}
	return toOrganization(res)
}
