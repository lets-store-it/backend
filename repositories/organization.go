package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/evevseev/storeit/backend/generated/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type OrganizationRepository struct {
	Queries *database.Queries
}

type AuditFields struct {
	CreatedAt time.Time
	UpdatedAt *time.Time
	CreatedBy *uuid.UUID
	UpdatedBy *uuid.UUID
}

type Organization struct {
	ID          uuid.UUID
	Name        string
	Subdomain   string
	AuditFields *AuditFields
}

func uuidFromPgx(id pgtype.UUID) *uuid.UUID {
	if !id.Valid {
		return nil
	}
	uuid := (uuid.UUID)(id.Bytes)
	return &uuid
}

func datetimeFromPgx(datetime pgtype.Timestamp) *time.Time {
	if !datetime.Valid {
		return nil
	}
	return &datetime.Time
}

func toOrganization(org database.Org) (*Organization, error) {
	id := uuidFromPgx(org.ID)
	if id == nil {
		return nil, errors.New("id is nil")
	}
	return &Organization{
		ID:        *id,
		Name:      org.Name,
		Subdomain: org.Subdomain,
		AuditFields: &AuditFields{
			CreatedAt: org.CreatedAt.Time,
			UpdatedAt: datetimeFromPgx(org.UpdatedAt),
			CreatedBy: uuidFromPgx(org.CreatedBy),
			UpdatedBy: uuidFromPgx(org.UpdatedBy),
		},
	}, nil
}

func (r *OrganizationRepository) GetOrgs(ctx context.Context, limit int32, offset int32) ([]*Organization, error) {
	res, err := r.Queries.GetOrgs(ctx, database.GetOrgsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	orgs := make([]*Organization, len(res))
	for i, org := range res {
		orgs[i], err = toOrganization(org)
		if err != nil {
			return nil, err
		}
	}

	return orgs, nil
}

func (r *OrganizationRepository) CreateOrg(ctx context.Context, name string, subdomain string) (*Organization, error) {
	org, err := r.Queries.CreateOrg(ctx, database.CreateOrgParams{
		Name:      name,
		Subdomain: subdomain,
	})
	if err != nil {
		return nil, err
	}

	return toOrganization(org)
}
