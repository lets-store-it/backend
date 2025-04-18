// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const assignRoleToUser = `-- name: AssignRoleToUser :exec
INSERT INTO app_role_binding (role_id, user_id, org_id) VALUES ($1, $2, $3)
`

type AssignRoleToUserParams struct {
	RoleID int32
	UserID pgtype.UUID
	OrgID  pgtype.UUID
}

// Role Bindings
func (q *Queries) AssignRoleToUser(ctx context.Context, arg AssignRoleToUserParams) error {
	_, err := q.db.Exec(ctx, assignRoleToUser, arg.RoleID, arg.UserID, arg.OrgID)
	return err
}

const createApiToken = `-- name: CreateApiToken :one
INSERT INTO app_api_token (org_id, token) VALUES ($1, $2) RETURNING id, org_id, name, token, created_at, revoked_at
`

type CreateApiTokenParams struct {
	OrgID pgtype.UUID
	Token string
}

func (q *Queries) CreateApiToken(ctx context.Context, arg CreateApiTokenParams) (AppApiToken, error) {
	row := q.db.QueryRow(ctx, createApiToken, arg.OrgID, arg.Token)
	var i AppApiToken
	err := row.Scan(
		&i.ID,
		&i.OrgID,
		&i.Name,
		&i.Token,
		&i.CreatedAt,
		&i.RevokedAt,
	)
	return i, err
}

const createCell = `-- name: CreateCell :one
INSERT INTO cell (org_id, cells_group_id, alias, row, level, position) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, org_id, cells_group_id, alias, row, level, position, created_at, deleted_at
`

type CreateCellParams struct {
	OrgID        pgtype.UUID
	CellsGroupID pgtype.UUID
	Alias        string
	Row          int32
	Level        int32
	Position     int32
}

func (q *Queries) CreateCell(ctx context.Context, arg CreateCellParams) (Cell, error) {
	row := q.db.QueryRow(ctx, createCell,
		arg.OrgID,
		arg.CellsGroupID,
		arg.Alias,
		arg.Row,
		arg.Level,
		arg.Position,
	)
	var i Cell
	err := row.Scan(
		&i.ID,
		&i.OrgID,
		&i.CellsGroupID,
		&i.Alias,
		&i.Row,
		&i.Level,
		&i.Position,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const createCellsGroup = `-- name: CreateCellsGroup :one
INSERT INTO cells_group (org_id, storage_group_id, name, alias) VALUES ($1, $2, $3, $4) RETURNING id, org_id, storage_group_id, name, alias, created_at, deleted_at
`

type CreateCellsGroupParams struct {
	OrgID          pgtype.UUID
	StorageGroupID pgtype.UUID
	Name           string
	Alias          string
}

func (q *Queries) CreateCellsGroup(ctx context.Context, arg CreateCellsGroupParams) (CellsGroup, error) {
	row := q.db.QueryRow(ctx, createCellsGroup,
		arg.OrgID,
		arg.StorageGroupID,
		arg.Name,
		arg.Alias,
	)
	var i CellsGroup
	err := row.Scan(
		&i.ID,
		&i.OrgID,
		&i.StorageGroupID,
		&i.Name,
		&i.Alias,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const createItem = `-- name: CreateItem :one
INSERT INTO item (org_id, name, description) VALUES ($1, $2, $3) RETURNING id, org_id, name, description, width, depth, height, weight, created_at, deleted_at
`

type CreateItemParams struct {
	OrgID       pgtype.UUID
	Name        string
	Description pgtype.Text
}

func (q *Queries) CreateItem(ctx context.Context, arg CreateItemParams) (Item, error) {
	row := q.db.QueryRow(ctx, createItem, arg.OrgID, arg.Name, arg.Description)
	var i Item
	err := row.Scan(
		&i.ID,
		&i.OrgID,
		&i.Name,
		&i.Description,
		&i.Width,
		&i.Depth,
		&i.Height,
		&i.Weight,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const createItemInstance = `-- name: CreateItemInstance :one
INSERT INTO item_instance (org_id, item_id, variant_id, cell_id, status, affected_by_operation_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, org_id, item_id, variant_id, cell_id, status, affected_by_operation_id, created_at, deleted_at
`

type CreateItemInstanceParams struct {
	OrgID                 pgtype.UUID
	ItemID                pgtype.UUID
	VariantID             pgtype.UUID
	CellID                pgtype.UUID
	Status                string
	AffectedByOperationID pgtype.UUID
}

// Item Instances
func (q *Queries) CreateItemInstance(ctx context.Context, arg CreateItemInstanceParams) (ItemInstance, error) {
	row := q.db.QueryRow(ctx, createItemInstance,
		arg.OrgID,
		arg.ItemID,
		arg.VariantID,
		arg.CellID,
		arg.Status,
		arg.AffectedByOperationID,
	)
	var i ItemInstance
	err := row.Scan(
		&i.ID,
		&i.OrgID,
		&i.ItemID,
		&i.VariantID,
		&i.CellID,
		&i.Status,
		&i.AffectedByOperationID,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const createItemVariant = `-- name: CreateItemVariant :one
INSERT INTO item_variant (org_id, item_id, name, article, ean13) VALUES ($1, $2, $3, $4, $5) RETURNING id, org_id, item_id, name, article, ean13, created_at, deleted_at
`

type CreateItemVariantParams struct {
	OrgID   pgtype.UUID
	ItemID  pgtype.UUID
	Name    string
	Article pgtype.Text
	Ean13   pgtype.Int4
}

func (q *Queries) CreateItemVariant(ctx context.Context, arg CreateItemVariantParams) (ItemVariant, error) {
	row := q.db.QueryRow(ctx, createItemVariant,
		arg.OrgID,
		arg.ItemID,
		arg.Name,
		arg.Article,
		arg.Ean13,
	)
	var i ItemVariant
	err := row.Scan(
		&i.ID,
		&i.OrgID,
		&i.ItemID,
		&i.Name,
		&i.Article,
		&i.Ean13,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const createObjectChange = `-- name: CreateObjectChange :one
INSERT INTO app_object_changes (org_id, user_id, action, target_object_type, target_object_id, prechange_state, postchange_state) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, org_id, user_id, action, time, target_object_type, target_object_id, prechange_state, postchange_state
`

type CreateObjectChangeParams struct {
	OrgID            pgtype.UUID
	UserID           pgtype.UUID
	Action           string
	TargetObjectType int32
	TargetObjectID   pgtype.UUID
	PrechangeState   []byte
	PostchangeState  []byte
}

// Audit Log
func (q *Queries) CreateObjectChange(ctx context.Context, arg CreateObjectChangeParams) (AppObjectChange, error) {
	row := q.db.QueryRow(ctx, createObjectChange,
		arg.OrgID,
		arg.UserID,
		arg.Action,
		arg.TargetObjectType,
		arg.TargetObjectID,
		arg.PrechangeState,
		arg.PostchangeState,
	)
	var i AppObjectChange
	err := row.Scan(
		&i.ID,
		&i.OrgID,
		&i.UserID,
		&i.Action,
		&i.Time,
		&i.TargetObjectType,
		&i.TargetObjectID,
		&i.PrechangeState,
		&i.PostchangeState,
	)
	return i, err
}

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

const createOrgUnit = `-- name: CreateOrgUnit :one
INSERT INTO org_unit (org_id, name, alias, address) VALUES ($1, $2, $3, $4) RETURNING id, org_id, name, alias, address, created_at, deleted_at
`

type CreateOrgUnitParams struct {
	OrgID   pgtype.UUID
	Name    string
	Alias   string
	Address pgtype.Text
}

func (q *Queries) CreateOrgUnit(ctx context.Context, arg CreateOrgUnitParams) (OrgUnit, error) {
	row := q.db.QueryRow(ctx, createOrgUnit,
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

const createStorageGroup = `-- name: CreateStorageGroup :one
INSERT INTO storage_group (org_id, unit_id, parent_id, name, alias) VALUES ($1, $2, $3, $4, $5) RETURNING id, org_id, unit_id, parent_id, name, alias, description, created_at, deleted_at
`

type CreateStorageGroupParams struct {
	OrgID    pgtype.UUID
	UnitID   pgtype.UUID
	ParentID pgtype.UUID
	Name     string
	Alias    string
}

func (q *Queries) CreateStorageGroup(ctx context.Context, arg CreateStorageGroupParams) (StorageGroup, error) {
	row := q.db.QueryRow(ctx, createStorageGroup,
		arg.OrgID,
		arg.UnitID,
		arg.ParentID,
		arg.Name,
		arg.Alias,
	)
	var i StorageGroup
	err := row.Scan(
		&i.ID,
		&i.OrgID,
		&i.UnitID,
		&i.ParentID,
		&i.Name,
		&i.Alias,
		&i.Description,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO app_user (email, first_name, last_name, middle_name, yandex_id) VALUES ($1, $2, $3, $4, $5) RETURNING id, email, first_name, last_name, middle_name, yandex_id, created_at
`

type CreateUserParams struct {
	Email      string
	FirstName  string
	LastName   string
	MiddleName pgtype.Text
	YandexID   pgtype.Text
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (AppUser, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.Email,
		arg.FirstName,
		arg.LastName,
		arg.MiddleName,
		arg.YandexID,
	)
	var i AppUser
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.FirstName,
		&i.LastName,
		&i.MiddleName,
		&i.YandexID,
		&i.CreatedAt,
	)
	return i, err
}

const createUserSession = `-- name: CreateUserSession :one
INSERT INTO app_user_session (user_id, token) VALUES ($1, $2) RETURNING id, user_id, token, created_at, expires_at, revoked_at
`

type CreateUserSessionParams struct {
	UserID pgtype.UUID
	Token  string
}

func (q *Queries) CreateUserSession(ctx context.Context, arg CreateUserSessionParams) (AppUserSession, error) {
	row := q.db.QueryRow(ctx, createUserSession, arg.UserID, arg.Token)
	var i AppUserSession
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Token,
		&i.CreatedAt,
		&i.ExpiresAt,
		&i.RevokedAt,
	)
	return i, err
}

const deleteCell = `-- name: DeleteCell :exec
UPDATE cell SET deleted_at = CURRENT_TIMESTAMP WHERE org_id = $1 AND cells_group_id = $2 AND id = $3
`

type DeleteCellParams struct {
	OrgID        pgtype.UUID
	CellsGroupID pgtype.UUID
	ID           pgtype.UUID
}

func (q *Queries) DeleteCell(ctx context.Context, arg DeleteCellParams) error {
	_, err := q.db.Exec(ctx, deleteCell, arg.OrgID, arg.CellsGroupID, arg.ID)
	return err
}

const deleteCellsGroup = `-- name: DeleteCellsGroup :exec
UPDATE cells_group SET deleted_at = CURRENT_TIMESTAMP WHERE org_id = $1 AND id = $2
`

type DeleteCellsGroupParams struct {
	OrgID pgtype.UUID
	ID    pgtype.UUID
}

func (q *Queries) DeleteCellsGroup(ctx context.Context, arg DeleteCellsGroupParams) error {
	_, err := q.db.Exec(ctx, deleteCellsGroup, arg.OrgID, arg.ID)
	return err
}

const deleteItem = `-- name: DeleteItem :exec
UPDATE item SET deleted_at = CURRENT_TIMESTAMP WHERE org_id = $1 AND id = $2
`

type DeleteItemParams struct {
	OrgID pgtype.UUID
	ID    pgtype.UUID
}

func (q *Queries) DeleteItem(ctx context.Context, arg DeleteItemParams) error {
	_, err := q.db.Exec(ctx, deleteItem, arg.OrgID, arg.ID)
	return err
}

const deleteItemVariant = `-- name: DeleteItemVariant :exec
UPDATE item_variant SET deleted_at = CURRENT_TIMESTAMP WHERE org_id = $1 AND item_id = $2 AND id = $3
`

type DeleteItemVariantParams struct {
	OrgID  pgtype.UUID
	ItemID pgtype.UUID
	ID     pgtype.UUID
}

func (q *Queries) DeleteItemVariant(ctx context.Context, arg DeleteItemVariantParams) error {
	_, err := q.db.Exec(ctx, deleteItemVariant, arg.OrgID, arg.ItemID, arg.ID)
	return err
}

const deleteOrg = `-- name: DeleteOrg :exec
UPDATE org SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1
`

func (q *Queries) DeleteOrg(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteOrg, id)
	return err
}

const deleteOrgUnit = `-- name: DeleteOrgUnit :exec
UPDATE org_unit SET deleted_at = CURRENT_TIMESTAMP WHERE org_id = $1 AND id = $2
`

type DeleteOrgUnitParams struct {
	OrgID pgtype.UUID
	ID    pgtype.UUID
}

func (q *Queries) DeleteOrgUnit(ctx context.Context, arg DeleteOrgUnitParams) error {
	_, err := q.db.Exec(ctx, deleteOrgUnit, arg.OrgID, arg.ID)
	return err
}

const deleteStorageGroup = `-- name: DeleteStorageGroup :exec
UPDATE storage_group SET deleted_at = CURRENT_TIMESTAMP WHERE org_id = $1 AND id = $2
`

type DeleteStorageGroupParams struct {
	OrgID pgtype.UUID
	ID    pgtype.UUID
}

func (q *Queries) DeleteStorageGroup(ctx context.Context, arg DeleteStorageGroupParams) error {
	_, err := q.db.Exec(ctx, deleteStorageGroup, arg.OrgID, arg.ID)
	return err
}

const getApiTokens = `-- name: GetApiTokens :many
SELECT id, org_id, name, token, created_at, revoked_at FROM app_api_token WHERE org_id = $1 AND revoked_at IS NULL
`

// Api Tokens
func (q *Queries) GetApiTokens(ctx context.Context, orgID pgtype.UUID) ([]AppApiToken, error) {
	rows, err := q.db.Query(ctx, getApiTokens, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AppApiToken
	for rows.Next() {
		var i AppApiToken
		if err := rows.Scan(
			&i.ID,
			&i.OrgID,
			&i.Name,
			&i.Token,
			&i.CreatedAt,
			&i.RevokedAt,
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

const getCell = `-- name: GetCell :one
SELECT id, org_id, cells_group_id, alias, row, level, position, created_at, deleted_at FROM cell WHERE org_id = $1 AND cells_group_id = $2 AND id = $3 AND deleted_at IS NULL
`

type GetCellParams struct {
	OrgID        pgtype.UUID
	CellsGroupID pgtype.UUID
	ID           pgtype.UUID
}

func (q *Queries) GetCell(ctx context.Context, arg GetCellParams) (Cell, error) {
	row := q.db.QueryRow(ctx, getCell, arg.OrgID, arg.CellsGroupID, arg.ID)
	var i Cell
	err := row.Scan(
		&i.ID,
		&i.OrgID,
		&i.CellsGroupID,
		&i.Alias,
		&i.Row,
		&i.Level,
		&i.Position,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getCells = `-- name: GetCells :many
SELECT id, org_id, cells_group_id, alias, row, level, position, created_at, deleted_at FROM cell WHERE org_id = $1 AND cells_group_id = $2 AND deleted_at IS NULL
`

type GetCellsParams struct {
	OrgID        pgtype.UUID
	CellsGroupID pgtype.UUID
}

// Cells
func (q *Queries) GetCells(ctx context.Context, arg GetCellsParams) ([]Cell, error) {
	rows, err := q.db.Query(ctx, getCells, arg.OrgID, arg.CellsGroupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Cell
	for rows.Next() {
		var i Cell
		if err := rows.Scan(
			&i.ID,
			&i.OrgID,
			&i.CellsGroupID,
			&i.Alias,
			&i.Row,
			&i.Level,
			&i.Position,
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

const getCellsGroup = `-- name: GetCellsGroup :one
SELECT id, org_id, storage_group_id, name, alias, created_at, deleted_at FROM cells_group WHERE org_id = $1 AND id = $2 AND deleted_at IS NULL
`

type GetCellsGroupParams struct {
	OrgID pgtype.UUID
	ID    pgtype.UUID
}

func (q *Queries) GetCellsGroup(ctx context.Context, arg GetCellsGroupParams) (CellsGroup, error) {
	row := q.db.QueryRow(ctx, getCellsGroup, arg.OrgID, arg.ID)
	var i CellsGroup
	err := row.Scan(
		&i.ID,
		&i.OrgID,
		&i.StorageGroupID,
		&i.Name,
		&i.Alias,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getCellsGroups = `-- name: GetCellsGroups :many
SELECT id, org_id, storage_group_id, name, alias, created_at, deleted_at FROM cells_group WHERE org_id = $1 AND deleted_at IS NULL
`

// CellsGroups
func (q *Queries) GetCellsGroups(ctx context.Context, orgID pgtype.UUID) ([]CellsGroup, error) {
	rows, err := q.db.Query(ctx, getCellsGroups, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CellsGroup
	for rows.Next() {
		var i CellsGroup
		if err := rows.Scan(
			&i.ID,
			&i.OrgID,
			&i.StorageGroupID,
			&i.Name,
			&i.Alias,
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

const getItem = `-- name: GetItem :one
SELECT id, org_id, name, description, width, depth, height, weight, created_at, deleted_at FROM item WHERE org_id = $1 AND id = $2 AND deleted_at IS NULL
`

type GetItemParams struct {
	OrgID pgtype.UUID
	ID    pgtype.UUID
}

func (q *Queries) GetItem(ctx context.Context, arg GetItemParams) (Item, error) {
	row := q.db.QueryRow(ctx, getItem, arg.OrgID, arg.ID)
	var i Item
	err := row.Scan(
		&i.ID,
		&i.OrgID,
		&i.Name,
		&i.Description,
		&i.Width,
		&i.Depth,
		&i.Height,
		&i.Weight,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getItemInstancesForCell = `-- name: GetItemInstancesForCell :many
SELECT id, org_id, item_id, variant_id, cell_id, status, affected_by_operation_id, created_at, deleted_at FROM item_instance WHERE org_id = $1 AND cell_id = $2 AND deleted_at IS NULL
`

type GetItemInstancesForCellParams struct {
	OrgID  pgtype.UUID
	CellID pgtype.UUID
}

func (q *Queries) GetItemInstancesForCell(ctx context.Context, arg GetItemInstancesForCellParams) ([]ItemInstance, error) {
	rows, err := q.db.Query(ctx, getItemInstancesForCell, arg.OrgID, arg.CellID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ItemInstance
	for rows.Next() {
		var i ItemInstance
		if err := rows.Scan(
			&i.ID,
			&i.OrgID,
			&i.ItemID,
			&i.VariantID,
			&i.CellID,
			&i.Status,
			&i.AffectedByOperationID,
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

const getItemInstancesForCellsGroup = `-- name: GetItemInstancesForCellsGroup :many
SELECT id, org_id, item_id, variant_id, cell_id, status, affected_by_operation_id, created_at, deleted_at FROM item_instance WHERE item_instance.org_id = $1 AND cell_id IN (SELECT id FROM cell WHERE cells_group_id = $2 AND deleted_at IS NULL) AND deleted_at IS NULL
`

type GetItemInstancesForCellsGroupParams struct {
	OrgID        pgtype.UUID
	CellsGroupID pgtype.UUID
}

func (q *Queries) GetItemInstancesForCellsGroup(ctx context.Context, arg GetItemInstancesForCellsGroupParams) ([]ItemInstance, error) {
	rows, err := q.db.Query(ctx, getItemInstancesForCellsGroup, arg.OrgID, arg.CellsGroupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ItemInstance
	for rows.Next() {
		var i ItemInstance
		if err := rows.Scan(
			&i.ID,
			&i.OrgID,
			&i.ItemID,
			&i.VariantID,
			&i.CellID,
			&i.Status,
			&i.AffectedByOperationID,
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

const getItemInstancesForItem = `-- name: GetItemInstancesForItem :many
SELECT id, org_id, item_id, variant_id, cell_id, status, affected_by_operation_id, created_at, deleted_at FROM item_instance WHERE org_id = $1 AND item_id = $2 AND deleted_at IS NULL
`

type GetItemInstancesForItemParams struct {
	OrgID  pgtype.UUID
	ItemID pgtype.UUID
}

func (q *Queries) GetItemInstancesForItem(ctx context.Context, arg GetItemInstancesForItemParams) ([]ItemInstance, error) {
	rows, err := q.db.Query(ctx, getItemInstancesForItem, arg.OrgID, arg.ItemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ItemInstance
	for rows.Next() {
		var i ItemInstance
		if err := rows.Scan(
			&i.ID,
			&i.OrgID,
			&i.ItemID,
			&i.VariantID,
			&i.CellID,
			&i.Status,
			&i.AffectedByOperationID,
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

const getItemVariants = `-- name: GetItemVariants :many

SELECT id, org_id, item_id, name, article, ean13, created_at, deleted_at FROM item_variant WHERE org_id = $1 AND item_id = $2 AND deleted_at IS NULL
`

type GetItemVariantsParams struct {
	OrgID  pgtype.UUID
	ItemID pgtype.UUID
}

// -- name: GetActiveItemWithVariants :many
// SELECT * FROM item
// JOIN item_variant ON item.id = item_variant.item_id
// WHERE item.org_id = $1 AND item.deleted_at IS NULL
// GROUP BY item.id;
func (q *Queries) GetItemVariants(ctx context.Context, arg GetItemVariantsParams) ([]ItemVariant, error) {
	rows, err := q.db.Query(ctx, getItemVariants, arg.OrgID, arg.ItemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ItemVariant
	for rows.Next() {
		var i ItemVariant
		if err := rows.Scan(
			&i.ID,
			&i.OrgID,
			&i.ItemID,
			&i.Name,
			&i.Article,
			&i.Ean13,
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

const getItems = `-- name: GetItems :many
SELECT id, org_id, name, description, width, depth, height, weight, created_at, deleted_at FROM item WHERE org_id = $1 AND deleted_at IS NULL
`

// Items
func (q *Queries) GetItems(ctx context.Context, orgID pgtype.UUID) ([]Item, error) {
	rows, err := q.db.Query(ctx, getItems, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Item
	for rows.Next() {
		var i Item
		if err := rows.Scan(
			&i.ID,
			&i.OrgID,
			&i.Name,
			&i.Description,
			&i.Width,
			&i.Depth,
			&i.Height,
			&i.Weight,
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

const getObjectChanges = `-- name: GetObjectChanges :many
SELECT id, org_id, user_id, action, time, target_object_type, target_object_id, prechange_state, postchange_state FROM app_object_changes WHERE org_id = $1 AND target_object_type = $2 AND target_object_id = $3 AND deleted_at IS NULL
`

type GetObjectChangesParams struct {
	OrgID            pgtype.UUID
	TargetObjectType int32
	TargetObjectID   pgtype.UUID
}

func (q *Queries) GetObjectChanges(ctx context.Context, arg GetObjectChangesParams) ([]AppObjectChange, error) {
	rows, err := q.db.Query(ctx, getObjectChanges, arg.OrgID, arg.TargetObjectType, arg.TargetObjectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AppObjectChange
	for rows.Next() {
		var i AppObjectChange
		if err := rows.Scan(
			&i.ID,
			&i.OrgID,
			&i.UserID,
			&i.Action,
			&i.Time,
			&i.TargetObjectType,
			&i.TargetObjectID,
			&i.PrechangeState,
			&i.PostchangeState,
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

const getOrg = `-- name: GetOrg :one
SELECT id, name, subdomain, created_at, deleted_at FROM org WHERE id = $1 AND deleted_at IS NULL
`

func (q *Queries) GetOrg(ctx context.Context, id pgtype.UUID) (Org, error) {
	row := q.db.QueryRow(ctx, getOrg, id)
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

const getOrgIdByApiToken = `-- name: GetOrgIdByApiToken :one
SELECT org_id FROM app_api_token WHERE token = $1 AND revoked_at IS NULL
`

func (q *Queries) GetOrgIdByApiToken(ctx context.Context, token string) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, getOrgIdByApiToken, token)
	var org_id pgtype.UUID
	err := row.Scan(&org_id)
	return org_id, err
}

const getOrgUnit = `-- name: GetOrgUnit :one
SELECT id, org_id, name, alias, address, created_at, deleted_at FROM org_unit WHERE org_id = $1 AND id = $2 AND deleted_at IS NULL
`

type GetOrgUnitParams struct {
	OrgID pgtype.UUID
	ID    pgtype.UUID
}

func (q *Queries) GetOrgUnit(ctx context.Context, arg GetOrgUnitParams) (OrgUnit, error) {
	row := q.db.QueryRow(ctx, getOrgUnit, arg.OrgID, arg.ID)
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

const getOrgUnits = `-- name: GetOrgUnits :many
SELECT id, org_id, name, alias, address, created_at, deleted_at FROM org_unit WHERE org_id = $1 AND deleted_at IS NULL
`

// Units
func (q *Queries) GetOrgUnits(ctx context.Context, orgID pgtype.UUID) ([]OrgUnit, error) {
	rows, err := q.db.Query(ctx, getOrgUnits, orgID)
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

const getStorageGroup = `-- name: GetStorageGroup :one
SELECT id, org_id, unit_id, parent_id, name, alias, description, created_at, deleted_at FROM storage_group WHERE org_id = $1 AND id = $2 AND deleted_at IS NULL
`

type GetStorageGroupParams struct {
	OrgID pgtype.UUID
	ID    pgtype.UUID
}

func (q *Queries) GetStorageGroup(ctx context.Context, arg GetStorageGroupParams) (StorageGroup, error) {
	row := q.db.QueryRow(ctx, getStorageGroup, arg.OrgID, arg.ID)
	var i StorageGroup
	err := row.Scan(
		&i.ID,
		&i.OrgID,
		&i.UnitID,
		&i.ParentID,
		&i.Name,
		&i.Alias,
		&i.Description,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getStorageGroups = `-- name: GetStorageGroups :many
SELECT id, org_id, unit_id, parent_id, name, alias, description, created_at, deleted_at FROM storage_group WHERE org_id = $1 AND deleted_at IS NULL
`

// Storage spaces
func (q *Queries) GetStorageGroups(ctx context.Context, orgID pgtype.UUID) ([]StorageGroup, error) {
	rows, err := q.db.Query(ctx, getStorageGroups, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []StorageGroup
	for rows.Next() {
		var i StorageGroup
		if err := rows.Scan(
			&i.ID,
			&i.OrgID,
			&i.UnitID,
			&i.ParentID,
			&i.Name,
			&i.Alias,
			&i.Description,
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

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, email, first_name, last_name, middle_name, yandex_id, created_at FROM app_user WHERE email = $1 LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (AppUser, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i AppUser
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.FirstName,
		&i.LastName,
		&i.MiddleName,
		&i.YandexID,
		&i.CreatedAt,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, email, first_name, last_name, middle_name, yandex_id, created_at FROM app_user WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUserById(ctx context.Context, id pgtype.UUID) (AppUser, error) {
	row := q.db.QueryRow(ctx, getUserById, id)
	var i AppUser
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.FirstName,
		&i.LastName,
		&i.MiddleName,
		&i.YandexID,
		&i.CreatedAt,
	)
	return i, err
}

const getUserBySessionSecret = `-- name: GetUserBySessionSecret :one
SELECT id, email, first_name, last_name, middle_name, yandex_id, created_at FROM app_user WHERE id = (SELECT user_id FROM app_user_session WHERE token = $1 AND expires_at > CURRENT_TIMESTAMP LIMIT 1)
`

// Auth
func (q *Queries) GetUserBySessionSecret(ctx context.Context, token string) (AppUser, error) {
	row := q.db.QueryRow(ctx, getUserBySessionSecret, token)
	var i AppUser
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.FirstName,
		&i.LastName,
		&i.MiddleName,
		&i.YandexID,
		&i.CreatedAt,
	)
	return i, err
}

const getUserOrgs = `-- name: GetUserOrgs :many
SELECT id, name, subdomain, created_at, deleted_at FROM org WHERE id IN (SELECT org_id FROM app_role_binding WHERE user_id = $1)
`

func (q *Queries) GetUserOrgs(ctx context.Context, userID pgtype.UUID) ([]Org, error) {
	rows, err := q.db.Query(ctx, getUserOrgs, userID)
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

const getUserRolesInOrg = `-- name: GetUserRolesInOrg :many
SELECT id, org_id, role_id, user_id FROM app_role_binding WHERE user_id = $1 AND org_id = $2
`

type GetUserRolesInOrgParams struct {
	UserID pgtype.UUID
	OrgID  pgtype.UUID
}

func (q *Queries) GetUserRolesInOrg(ctx context.Context, arg GetUserRolesInOrgParams) ([]AppRoleBinding, error) {
	rows, err := q.db.Query(ctx, getUserRolesInOrg, arg.UserID, arg.OrgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AppRoleBinding
	for rows.Next() {
		var i AppRoleBinding
		if err := rows.Scan(
			&i.ID,
			&i.OrgID,
			&i.RoleID,
			&i.UserID,
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

const isItemExists = `-- name: IsItemExists :one
SELECT EXISTS (SELECT 1 FROM item WHERE org_id = $1 AND id = $2 AND deleted_at IS NULL)
`

type IsItemExistsParams struct {
	OrgID pgtype.UUID
	ID    pgtype.UUID
}

func (q *Queries) IsItemExists(ctx context.Context, arg IsItemExistsParams) (bool, error) {
	row := q.db.QueryRow(ctx, isItemExists, arg.OrgID, arg.ID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const isOrgExists = `-- name: IsOrgExists :one
SELECT EXISTS (SELECT 1 FROM org WHERE id = $1 AND deleted_at IS NULL)
`

func (q *Queries) IsOrgExists(ctx context.Context, id pgtype.UUID) (bool, error) {
	row := q.db.QueryRow(ctx, isOrgExists, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const revokeApiToken = `-- name: RevokeApiToken :exec
UPDATE app_api_token SET revoked_at = CURRENT_TIMESTAMP WHERE org_id = $1 AND token = $2
`

type RevokeApiTokenParams struct {
	OrgID pgtype.UUID
	Token string
}

func (q *Queries) RevokeApiToken(ctx context.Context, arg RevokeApiTokenParams) error {
	_, err := q.db.Exec(ctx, revokeApiToken, arg.OrgID, arg.Token)
	return err
}

const unassignRoleFromUser = `-- name: UnassignRoleFromUser :exec
DELETE FROM app_role_binding WHERE role_id = $1 AND user_id = $2 AND org_id = $3
`

type UnassignRoleFromUserParams struct {
	RoleID int32
	UserID pgtype.UUID
	OrgID  pgtype.UUID
}

func (q *Queries) UnassignRoleFromUser(ctx context.Context, arg UnassignRoleFromUserParams) error {
	_, err := q.db.Exec(ctx, unassignRoleFromUser, arg.RoleID, arg.UserID, arg.OrgID)
	return err
}

const updateCell = `-- name: UpdateCell :one
UPDATE cell SET alias = $4, row = $5, level = $6, position = $7 WHERE org_id = $1 AND cells_group_id = $2 AND id = $3 AND deleted_at IS NULL RETURNING id, org_id, cells_group_id, alias, row, level, position, created_at, deleted_at
`

type UpdateCellParams struct {
	OrgID        pgtype.UUID
	CellsGroupID pgtype.UUID
	ID           pgtype.UUID
	Alias        string
	Row          int32
	Level        int32
	Position     int32
}

func (q *Queries) UpdateCell(ctx context.Context, arg UpdateCellParams) (Cell, error) {
	row := q.db.QueryRow(ctx, updateCell,
		arg.OrgID,
		arg.CellsGroupID,
		arg.ID,
		arg.Alias,
		arg.Row,
		arg.Level,
		arg.Position,
	)
	var i Cell
	err := row.Scan(
		&i.ID,
		&i.OrgID,
		&i.CellsGroupID,
		&i.Alias,
		&i.Row,
		&i.Level,
		&i.Position,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const updateCellsGroup = `-- name: UpdateCellsGroup :one
UPDATE cells_group SET name = $3, alias = $4 WHERE org_id = $1 AND id = $2 AND deleted_at IS NULL RETURNING id, org_id, storage_group_id, name, alias, created_at, deleted_at
`

type UpdateCellsGroupParams struct {
	OrgID pgtype.UUID
	ID    pgtype.UUID
	Name  string
	Alias string
}

func (q *Queries) UpdateCellsGroup(ctx context.Context, arg UpdateCellsGroupParams) (CellsGroup, error) {
	row := q.db.QueryRow(ctx, updateCellsGroup,
		arg.OrgID,
		arg.ID,
		arg.Name,
		arg.Alias,
	)
	var i CellsGroup
	err := row.Scan(
		&i.ID,
		&i.OrgID,
		&i.StorageGroupID,
		&i.Name,
		&i.Alias,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const updateItem = `-- name: UpdateItem :one
UPDATE item SET name = $3, description = $4 WHERE org_id = $1 AND id = $2 AND deleted_at IS NULL RETURNING id, org_id, name, description, width, depth, height, weight, created_at, deleted_at
`

type UpdateItemParams struct {
	OrgID       pgtype.UUID
	ID          pgtype.UUID
	Name        string
	Description pgtype.Text
}

func (q *Queries) UpdateItem(ctx context.Context, arg UpdateItemParams) (Item, error) {
	row := q.db.QueryRow(ctx, updateItem,
		arg.OrgID,
		arg.ID,
		arg.Name,
		arg.Description,
	)
	var i Item
	err := row.Scan(
		&i.ID,
		&i.OrgID,
		&i.Name,
		&i.Description,
		&i.Width,
		&i.Depth,
		&i.Height,
		&i.Weight,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const updateItemVariant = `-- name: UpdateItemVariant :one
UPDATE item_variant SET name = $3, article = $4, ean13 = $5 WHERE org_id = $1 AND item_id = $2 AND id = $3 AND deleted_at IS NULL RETURNING id, org_id, item_id, name, article, ean13, created_at, deleted_at
`

type UpdateItemVariantParams struct {
	OrgID   pgtype.UUID
	ItemID  pgtype.UUID
	Name    string
	Article pgtype.Text
	Ean13   pgtype.Int4
}

func (q *Queries) UpdateItemVariant(ctx context.Context, arg UpdateItemVariantParams) (ItemVariant, error) {
	row := q.db.QueryRow(ctx, updateItemVariant,
		arg.OrgID,
		arg.ItemID,
		arg.Name,
		arg.Article,
		arg.Ean13,
	)
	var i ItemVariant
	err := row.Scan(
		&i.ID,
		&i.OrgID,
		&i.ItemID,
		&i.Name,
		&i.Article,
		&i.Ean13,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
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

const updateOrgUnit = `-- name: UpdateOrgUnit :one
UPDATE org_unit SET name = $3, alias = $4, address = $5 WHERE org_id = $1 AND id = $2 AND deleted_at IS NULL RETURNING id, org_id, name, alias, address, created_at, deleted_at
`

type UpdateOrgUnitParams struct {
	OrgID   pgtype.UUID
	ID      pgtype.UUID
	Name    string
	Alias   string
	Address pgtype.Text
}

func (q *Queries) UpdateOrgUnit(ctx context.Context, arg UpdateOrgUnitParams) (OrgUnit, error) {
	row := q.db.QueryRow(ctx, updateOrgUnit,
		arg.OrgID,
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

const updateStorageGroup = `-- name: UpdateStorageGroup :one
UPDATE storage_group SET name = $3, alias = $4 WHERE org_id = $1 AND id = $2 AND deleted_at IS NULL RETURNING id, org_id, unit_id, parent_id, name, alias, description, created_at, deleted_at
`

type UpdateStorageGroupParams struct {
	OrgID pgtype.UUID
	ID    pgtype.UUID
	Name  string
	Alias string
}

func (q *Queries) UpdateStorageGroup(ctx context.Context, arg UpdateStorageGroupParams) (StorageGroup, error) {
	row := q.db.QueryRow(ctx, updateStorageGroup,
		arg.OrgID,
		arg.ID,
		arg.Name,
		arg.Alias,
	)
	var i StorageGroup
	err := row.Scan(
		&i.ID,
		&i.OrgID,
		&i.UnitID,
		&i.ParentID,
		&i.Name,
		&i.Alias,
		&i.Description,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}
