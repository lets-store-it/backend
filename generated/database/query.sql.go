// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createOrg = `-- name: CreateOrg :one
INSERT INTO org (name, subdomain) VALUES ($1, $2) RETURNING id, name, subdomain, created_at, deleted_at
`

type CreateOrgParams struct {
	Name      string
	Subdomain string
}

func (q *Queries) CreateOrg(ctx context.Context, arg CreateOrgParams) (Org, error) {
	row := q.db.QueryRow(ctx, createOrg, arg.Name, arg.Subdomain)
	var i Org
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Subdomain,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const createOrganizationUnit = `-- name: CreateOrganizationUnit :one
INSERT INTO org_unit (org_id, name, alias, address) VALUES ($1, $2, $3, $4) RETURNING id, org_id, name, alias, address, created_at, deleted_at
`

type CreateOrganizationUnitParams struct {
	OrgID   pgtype.UUID
	Name    string
	Alias   string
	Address pgtype.Text
}

func (q *Queries) CreateOrganizationUnit(ctx context.Context, arg CreateOrganizationUnitParams) (OrgUnit, error) {
	row := q.db.QueryRow(ctx, createOrganizationUnit,
		arg.OrgID,
		arg.Name,
		arg.Alias,
		arg.Address,
	)
	var i OrgUnit
	err := row.Scan(
		&i.ID,
		&i.OrgID,
		&i.Name,
		&i.Alias,
		&i.Address,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const createStorageSpace = `-- name: CreateStorageSpace :one
INSERT INTO storage_space (org_id, unit_id, parent_id, name, alias) VALUES ($1, $2, $3, $4, $5) RETURNING id, org_id, unit_id, parent_id, name, alias, created_at, deleted_at
`

type CreateStorageSpaceParams struct {
	OrgID    pgtype.UUID
	UnitID   pgtype.UUID
	ParentID pgtype.UUID
	Name     string
	Alias    pgtype.Text
}

func (q *Queries) CreateStorageSpace(ctx context.Context, arg CreateStorageSpaceParams) (StorageSpace, error) {
	row := q.db.QueryRow(ctx, createStorageSpace,
		arg.OrgID,
		arg.UnitID,
		arg.ParentID,
		arg.Name,
		arg.Alias,
	)
	var i StorageSpace
	err := row.Scan(
		&i.ID,
		&i.OrgID,
		&i.UnitID,
		&i.ParentID,
		&i.Name,
		&i.Alias,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const deleteOrg = `-- name: DeleteOrg :exec
UPDATE org SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1
`

func (q *Queries) DeleteOrg(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteOrg, id)
	return err
}

const deleteOrganizationUnit = `-- name: DeleteOrganizationUnit :exec
UPDATE org_unit SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1
`

func (q *Queries) DeleteOrganizationUnit(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteOrganizationUnit, id)
	return err
}

const deleteStorageSpace = `-- name: DeleteStorageSpace :exec
UPDATE storage_space SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1
`

func (q *Queries) DeleteStorageSpace(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteStorageSpace, id)
	return err
}

const getOrgById = `-- name: GetOrgById :one
SELECT id, name, subdomain, created_at, deleted_at FROM org WHERE id = $1 AND deleted_at IS NULL
`

func (q *Queries) GetOrgById(ctx context.Context, id pgtype.UUID) (Org, error) {
	row := q.db.QueryRow(ctx, getOrgById, id)
	var i Org
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Subdomain,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getOrganizationUnitById = `-- name: GetOrganizationUnitById :one
SELECT id, org_id, name, alias, address, created_at, deleted_at FROM org_unit WHERE id = $1 AND deleted_at IS NULL
`

func (q *Queries) GetOrganizationUnitById(ctx context.Context, id pgtype.UUID) (OrgUnit, error) {
	row := q.db.QueryRow(ctx, getOrganizationUnitById, id)
	var i OrgUnit
	err := row.Scan(
		&i.ID,
		&i.OrgID,
		&i.Name,
		&i.Alias,
		&i.Address,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getOrganizationUnits = `-- name: GetOrganizationUnits :many

SELECT id, org_id, name, alias, address, created_at, deleted_at FROM org_unit WHERE org_id = $1 AND deleted_at IS NULL
`

// -- Units
func (q *Queries) GetOrganizationUnits(ctx context.Context, orgID pgtype.UUID) ([]OrgUnit, error) {
	rows, err := q.db.Query(ctx, getOrganizationUnits, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []OrgUnit
	for rows.Next() {
		var i OrgUnit
		if err := rows.Scan(
			&i.ID,
			&i.OrgID,
			&i.Name,
			&i.Alias,
			&i.Address,
			&i.CreatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getOrgs = `-- name: GetOrgs :many
SELECT id, name, subdomain, created_at, deleted_at FROM org WHERE deleted_at IS NULL
`

func (q *Queries) GetOrgs(ctx context.Context) ([]Org, error) {
	rows, err := q.db.Query(ctx, getOrgs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Org
	for rows.Next() {
		var i Org
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Subdomain,
			&i.CreatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getStorageSpaceById = `-- name: GetStorageSpaceById :one
SELECT id, org_id, unit_id, parent_id, name, alias, created_at, deleted_at FROM storage_space WHERE id = $1 AND deleted_at IS NULL
`

func (q *Queries) GetStorageSpaceById(ctx context.Context, id pgtype.UUID) (StorageSpace, error) {
	row := q.db.QueryRow(ctx, getStorageSpaceById, id)
	var i StorageSpace
	err := row.Scan(
		&i.ID,
		&i.OrgID,
		&i.UnitID,
		&i.ParentID,
		&i.Name,
		&i.Alias,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const isOrgExists = `-- name: IsOrgExists :one
SELECT EXISTS (SELECT 1 FROM org WHERE name = $1 OR subdomain = $2)
`

type IsOrgExistsParams struct {
	Name      string
	Subdomain string
}

func (q *Queries) IsOrgExists(ctx context.Context, arg IsOrgExistsParams) (bool, error) {
	row := q.db.QueryRow(ctx, isOrgExists, arg.Name, arg.Subdomain)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const isOrgExistsById = `-- name: IsOrgExistsById :one
SELECT EXISTS (SELECT 1 FROM org WHERE id = $1)
`

func (q *Queries) IsOrgExistsById(ctx context.Context, id pgtype.UUID) (bool, error) {
	row := q.db.QueryRow(ctx, isOrgExistsById, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const isOrganizationUnitExistsForOrganization = `-- name: IsOrganizationUnitExistsForOrganization :one
SELECT EXISTS (SELECT 1 FROM org_unit WHERE org_id = $1 AND id = $2 AND deleted_at IS NULL)
`

type IsOrganizationUnitExistsForOrganizationParams struct {
	OrgID pgtype.UUID
	ID    pgtype.UUID
}

func (q *Queries) IsOrganizationUnitExistsForOrganization(ctx context.Context, arg IsOrganizationUnitExistsForOrganizationParams) (bool, error) {
	row := q.db.QueryRow(ctx, isOrganizationUnitExistsForOrganization, arg.OrgID, arg.ID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const isStorageSpaceExistsForOrganization = `-- name: IsStorageSpaceExistsForOrganization :one
SELECT EXISTS (SELECT 1 FROM storage_space WHERE org_id = $1 AND id = $2 AND deleted_at IS NULL)
`

type IsStorageSpaceExistsForOrganizationParams struct {
	OrgID pgtype.UUID
	ID    pgtype.UUID
}

func (q *Queries) IsStorageSpaceExistsForOrganization(ctx context.Context, arg IsStorageSpaceExistsForOrganizationParams) (bool, error) {
	row := q.db.QueryRow(ctx, isStorageSpaceExistsForOrganization, arg.OrgID, arg.ID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const updateOrg = `-- name: UpdateOrg :one
UPDATE org SET name = $2, subdomain = $3 WHERE id = $1 RETURNING id, name, subdomain, created_at, deleted_at
`

type UpdateOrgParams struct {
	ID        pgtype.UUID
	Name      string
	Subdomain string
}

func (q *Queries) UpdateOrg(ctx context.Context, arg UpdateOrgParams) (Org, error) {
	row := q.db.QueryRow(ctx, updateOrg, arg.ID, arg.Name, arg.Subdomain)
	var i Org
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Subdomain,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const updateOrganizationUnit = `-- name: UpdateOrganizationUnit :one
UPDATE org_unit SET name = $2, alias = $3, address = $4 WHERE id = $1 AND deleted_at IS NULL RETURNING id, org_id, name, alias, address, created_at, deleted_at
`

type UpdateOrganizationUnitParams struct {
	ID      pgtype.UUID
	Name    string
	Alias   string
	Address pgtype.Text
}

func (q *Queries) UpdateOrganizationUnit(ctx context.Context, arg UpdateOrganizationUnitParams) (OrgUnit, error) {
	row := q.db.QueryRow(ctx, updateOrganizationUnit,
		arg.ID,
		arg.Name,
		arg.Alias,
		arg.Address,
	)
	var i OrgUnit
	err := row.Scan(
		&i.ID,
		&i.OrgID,
		&i.Name,
		&i.Alias,
		&i.Address,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const updateStorageSpace = `-- name: UpdateStorageSpace :one
UPDATE storage_space SET name = $2, alias = $3 WHERE id = $1 AND deleted_at IS NULL RETURNING id, org_id, unit_id, parent_id, name, alias, created_at, deleted_at
`

type UpdateStorageSpaceParams struct {
	ID    pgtype.UUID
	Name  string
	Alias pgtype.Text
}

func (q *Queries) UpdateStorageSpace(ctx context.Context, arg UpdateStorageSpaceParams) (StorageSpace, error) {
	row := q.db.QueryRow(ctx, updateStorageSpace, arg.ID, arg.Name, arg.Alias)
	var i StorageSpace
	err := row.Scan(
		&i.ID,
		&i.OrgID,
		&i.UnitID,
		&i.ParentID,
		&i.Name,
		&i.Alias,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}
