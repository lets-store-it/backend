-- name: IsOrgExists :one
SELECT EXISTS (SELECT 1 FROM org WHERE id = $1 AND deleted_at IS NULL);

-- name: GetOrg :one
SELECT * FROM org WHERE id = $1 AND deleted_at IS NULL;

-- name: CreateOrg :one
INSERT INTO org (name, subdomain) VALUES ($1, $2) RETURNING *;

-- name: UpdateOrg :one
UPDATE org SET name = $2, subdomain = $3 WHERE id = $1 RETURNING *;

-- name: DeleteOrg :exec
UPDATE org SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1;

-- Units
-- name: GetOrgUnits :many
SELECT * FROM org_unit WHERE org_id = $1 AND deleted_at IS NULL;

-- name: GetOrgUnit :one
SELECT * FROM org_unit WHERE org_id = $1 AND id = $2 AND deleted_at IS NULL;

-- name: CreateOrgUnit :one
INSERT INTO org_unit (org_id, name, alias, address) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateOrgUnit :one
UPDATE org_unit SET name = $3, alias = $4, address = $5 WHERE org_id = $1 AND id = $2 AND deleted_at IS NULL RETURNING *;

-- name: DeleteOrgUnit :exec
UPDATE org_unit SET deleted_at = CURRENT_TIMESTAMP WHERE org_id = $1 AND id = $2;

-- Storage spaces
-- name: GetStorageGroups :many
SELECT * FROM storage_group WHERE org_id = $1 AND deleted_at IS NULL;

-- name: GetStorageGroup :one
SELECT * FROM storage_group WHERE org_id = $1 AND id = $2 AND deleted_at IS NULL;

-- name: CreateStorageGroup :one
INSERT INTO storage_group (org_id, unit_id, parent_id, name, alias) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdateStorageGroup :one
UPDATE storage_group SET name = $3, alias = $4 WHERE org_id = $1 AND id = $2 AND deleted_at IS NULL RETURNING *;

-- name: DeleteStorageGroup :exec
UPDATE storage_group SET deleted_at = CURRENT_TIMESTAMP WHERE org_id = $1 AND id = $2;

-- Items
-- name: GetItems :many
SELECT * FROM item WHERE org_id = $1 AND deleted_at IS NULL;

-- name: GetItem :one
SELECT * FROM item WHERE org_id = $1 AND id = $2 AND deleted_at IS NULL;

-- -- name: GetActiveItemWithVariants :many
-- SELECT * FROM item 
-- JOIN item_variant ON item.id = item_variant.item_id
-- WHERE item.org_id = $1 AND item.deleted_at IS NULL
-- GROUP BY item.id;

-- name: GetItemVariants :many
SELECT * FROM item_variant WHERE org_id = $1 AND item_id = $2 AND deleted_at IS NULL;

-- name: CreateItem :one
INSERT INTO item (org_id, name, description) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateItem :one
UPDATE item SET name = $3, description = $4 WHERE org_id = $1 AND id = $2 AND deleted_at IS NULL RETURNING *;

-- name: DeleteItem :exec
UPDATE item SET deleted_at = CURRENT_TIMESTAMP WHERE org_id = $1 AND id = $2;

-- name: CreateItemVariant :one
INSERT INTO item_variant (org_id, item_id, name, article, ean13) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdateItemVariant :one
UPDATE item_variant SET name = $3, article = $4, ean13 = $5 WHERE org_id = $1 AND item_id = $2 AND id = $3 AND deleted_at IS NULL RETURNING *;

-- name: DeleteItemVariant :exec
UPDATE item_variant SET deleted_at = CURRENT_TIMESTAMP WHERE org_id = $1 AND item_id = $2 AND id = $3;

-- name: IsItemExists :one
SELECT EXISTS (SELECT 1 FROM item WHERE org_id = $1 AND id = $2 AND deleted_at IS NULL);

-- Auth
-- name: GetUserBySessionSecret :one
SELECT * FROM app_user WHERE id = (SELECT user_id FROM app_user_session WHERE token = $1 AND expires_at > CURRENT_TIMESTAMP LIMIT 1);

-- name: GetUserById :one
SELECT * FROM app_user WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM app_user WHERE email = $1 LIMIT 1;

-- name: CreateUserSession :one
INSERT INTO app_user_session (user_id, token) VALUES ($1, $2) RETURNING *;

-- name: CreateUser :one
INSERT INTO app_user (email, first_name, last_name, middle_name, yandex_id) VALUES ($1, $2, $3, $4, $5) RETURNING *;


-- Role Bindings
-- name: AssignRoleToUser :exec
INSERT INTO app_role_binding (role_id, user_id, org_id) VALUES ($1, $2, $3);

-- name: UnassignRoleFromUser :exec
DELETE FROM app_role_binding WHERE role_id = $1 AND user_id = $2 AND org_id = $3;

-- name: GetUserRolesInOrg :many
SELECT * FROM app_role_binding WHERE user_id = $1 AND org_id = $2;

-- name: GetUserOrgs :many
SELECT * FROM org WHERE id IN (SELECT org_id FROM app_role_binding WHERE user_id = $1);


-- CellsGroups
-- name: GetCellsGroups :many
SELECT * FROM cells_group WHERE org_id = $1 AND deleted_at IS NULL;

-- name: GetCellsGroup :one
SELECT * FROM cells_group WHERE org_id = $1 AND id = $2 AND deleted_at IS NULL;

-- name: CreateCellsGroup :one
INSERT INTO cells_group (org_id, storage_group_id, name, alias) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateCellsGroup :one
UPDATE cells_group SET name = $3, alias = $4 WHERE org_id = $1 AND id = $2 AND deleted_at IS NULL RETURNING *;

-- name: DeleteCellsGroup :exec
UPDATE cells_group SET deleted_at = CURRENT_TIMESTAMP WHERE org_id = $1 AND id = $2;


-- Cells
-- name: GetCells :many
SELECT * FROM cell WHERE org_id = $1 AND cells_group_id = $2 AND deleted_at IS NULL;

-- name: GetCell :one
SELECT * FROM cell WHERE org_id = $1 AND cells_group_id = $2 AND id = $3 AND deleted_at IS NULL;

-- name: CreateCell :one
INSERT INTO cell (org_id, cells_group_id, alias, row, level, position) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: UpdateCell :one
UPDATE cell SET alias = $4, row = $5, level = $6, position = $7 WHERE org_id = $1 AND cells_group_id = $2 AND id = $3 AND deleted_at IS NULL RETURNING *;

-- name: DeleteCell :exec
UPDATE cell SET deleted_at = CURRENT_TIMESTAMP WHERE org_id = $1 AND cells_group_id = $2 AND id = $3;


-- Audit Log
-- name: CreateObjectChange :one
INSERT INTO app_object_changes (org_id, user_id, action, target_object_type, target_object_id, prechange_state, postchange_state) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetObjectChanges :many
SELECT * FROM app_object_changes WHERE org_id = $1 AND target_object_type = $2 AND target_object_id = $3 AND deleted_at IS NULL;


-- Item Instances
-- name: CreateItemInstance :one
INSERT INTO item_instance (org_id, item_id, variant_id, cell_id, status, affected_by_operation_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetItemInstancesForItem :many
SELECT * FROM item_instance WHERE org_id = $1 AND item_id = $2 AND deleted_at IS NULL;

-- name: GetItemInstancesForCell :many
SELECT * FROM item_instance WHERE org_id = $1 AND cell_id = $2 AND deleted_at IS NULL;

-- name: GetItemInstancesForCellsGroup :many
SELECT * FROM item_instance WHERE item_instance.org_id = $1 AND cell_id IN (SELECT id FROM cell WHERE cells_group_id = $2 AND deleted_at IS NULL) AND deleted_at IS NULL;


-- Api Tokens
-- name: GetApiTokens :many
SELECT * FROM app_api_token WHERE org_id = $1 AND revoked_at IS NULL;

-- name: CreateApiToken :one
INSERT INTO app_api_token (org_id, token) VALUES ($1, $2) RETURNING *;

-- name: RevokeApiToken :exec
UPDATE app_api_token SET revoked_at = CURRENT_TIMESTAMP WHERE org_id = $1 AND token = $2;

-- name: GetOrgIdByApiToken :one
SELECT org_id FROM app_api_token WHERE token = $1 AND revoked_at IS NULL;

-- -- name: GetActiveItemVariants :many
-- SELECT * FROM item_variant WHERE item_id = $1 AND deleted_at IS NULL;

-- -- name: GetItemVariant :one
-- SELECT * FROM item_variant WHERE item_id = $1 AND id = $2 AND deleted_at IS NULL;


-- -- Custom fields
-- -- name: GetCustomFields :many
-- SELECT * FROM custom_field WHERE org_id = $1 AND deleted_at IS NULL;

-- -- Object Types
-- -- name: GetObjectTypes :many
-- SELECT id, group, name FROM object_type;

-- CREATE TABLE custom_field (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     org_id UUID NOT NULL REFERENCES org(id),
--     type VARCHAR(100) NOT NULL CHECK (type IN ('text', 'integer', 'decimal' 'boolean')),
--     name VARCHAR(100) NOT NULL,
--     label VARCHAR(100) NOT NULL,
--     description VARCHAR(255),
--     group_name VARCHAR(100),
--     created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
--     deleted_at TIMESTAMP,
--     UNIQUE (org_id, name)
-- );

-- -- name: IsCustomFieldExistsForOrganization :one
-- SELECT EXISTS (SELECT 1 FROM custom_field WHERE org_id = $1 AND name = $2 AND deleted_at IS NULL);

-- -- name: CreateCustomField :one
-- INSERT INTO custom_field (org_id, name, label, type, group_name, description) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- -- name: UpdateCustomField :one
-- UPDATE custom_field SET name = $2, label = $3, group_name = $5, description = $6 WHERE id = $1 AND deleted_at IS NULL RETURNING *;

-- -- name: DeleteCustomField :exec
-- UPDATE custom_field SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1;

-- -- name: GetCustomFieldById :one
-- SELECT * FROM custom_field WHERE id = $1 AND deleted_at IS NULL;

-- -- name: GetCustomFieldRelatedTypes :many
-- SELECT object_type_id FROM custom_field_related_types WHERE custom_field_id = $1;

-- -- name: AddCustomFieldRelatedType :exec
-- INSERT INTO custom_field_related_types (custom_field_id, object_type_id) VALUES ($1, $2);

-- -- name: DeleteCustomFieldRelatedType :exec
-- DELETE FROM custom_field_related_types WHERE custom_field_id = $1 AND object_type_id = $2;
