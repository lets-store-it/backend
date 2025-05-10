-- Organizations
-- name: GetUserOrgs :many
SELECT * FROM org WHERE id IN (SELECT org_id FROM app_role_binding WHERE user_id = $1) AND deleted_at IS NULL;


-- name: CreateOrganization :one
INSERT INTO org (name, subdomain) VALUES ($1, $2) RETURNING *;

-- name: GetOrganization :one
SELECT * FROM org WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateOrganization :one
UPDATE org SET name = $2 WHERE id = $1 RETURNING *;

-- name: DeleteOrganization :exec
UPDATE org SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1;

-- Units
-- name: CreateOrgUnit :one
INSERT INTO org_unit (org_id, name, alias, address) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetOrgUnits :many
SELECT * FROM org_unit WHERE org_id = $1 AND deleted_at IS NULL;

-- name: GetOrgUnitById :one
SELECT * FROM org_unit WHERE org_id = $1 AND id = $2 AND deleted_at IS NULL;

-- name: UpdateOrgUnit :one
UPDATE org_unit SET name = $3, alias = $4, address = $5 WHERE org_id = $1 AND id = $2 AND deleted_at IS NULL RETURNING *;

-- name: DeleteOrgUnit :exec
UPDATE org_unit SET deleted_at = CURRENT_TIMESTAMP WHERE org_id = $1 AND id = $2;

-- Storage Groups
-- name: CreateStorageGroup :one
INSERT INTO storage_group (org_id, unit_id, parent_id, name, alias) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetStorageGroupById :one
SELECT * FROM storage_group WHERE org_id = $1 AND id = $2 AND deleted_at IS NULL;

-- name: GetStorageGroups :many
SELECT * FROM storage_group WHERE org_id = $1 AND deleted_at IS NULL;

-- name: UpdateStorageGroup :one
UPDATE storage_group SET name = $3, alias = $4, unit_id = $5 WHERE org_id = $1 AND id = $2 AND deleted_at IS NULL RETURNING *;

-- name: DeleteStorageGroup :exec
UPDATE storage_group SET deleted_at = CURRENT_TIMESTAMP WHERE org_id = $1 AND id = $2;

-- CellsGroups
-- name: CreateCellsGroup :one
INSERT INTO cells_group (org_id, unit_id, storage_group_id, name, alias) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetCellsGroups :many
SELECT * FROM cells_group WHERE org_id = $1 AND deleted_at IS NULL;

-- name: GetCellsGroupById :one
SELECT * FROM cells_group WHERE org_id = $1 AND id = $2 AND deleted_at IS NULL;

-- name: UpdateCellsGroup :one
UPDATE cells_group SET name = $3, alias = $4, unit_id = $5 WHERE org_id = $1 AND id = $2 AND deleted_at IS NULL RETURNING *;

-- name: DeleteCellsGroup :exec
UPDATE cells_group SET deleted_at = CURRENT_TIMESTAMP WHERE org_id = $1 AND id = $2;

-- Cells
-- name: CreateCell :one
INSERT INTO cell (org_id, cells_group_id, alias, row, level, position) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetCellById :one
SELECT * FROM cell WHERE org_id = $1 AND id = $2 AND deleted_at IS NULL;

-- name: GetCells :many
SELECT * FROM cell WHERE org_id = $1 AND cells_group_id = $2 AND deleted_at IS NULL;

-- name: UpdateCell :one
UPDATE cell SET alias = $4, row = $5, level = $6, position = $7 WHERE org_id = $1 AND cells_group_id = $2 AND id = $3 AND deleted_at IS NULL RETURNING *;

-- name: DeleteCell :exec
UPDATE cell SET deleted_at = CURRENT_TIMESTAMP WHERE org_id = $1 AND id = $2;

-- name: GetCellPath :many
WITH RECURSIVE path AS (
  SELECT
    cg.id,
    'cells_group'      AS type,
    cg.alias,
    cg.name,
    cg.storage_group_id AS parent_group_id,
    NULL::UUID         AS unit_id,
    1                  AS lvl
  FROM cell c
  JOIN cells_group cg
    ON c.cells_group_id = cg.id
   AND c.org_id         = cg.org_id
  WHERE c.org_id = $1
    AND c.id     = $2

  UNION ALL

  SELECT
    sg.id,
    'storage_group'    AS type,
    sg.alias,
    sg.name,
    sg.parent_id       AS parent_group_id,
    sg.unit_id,
    p.lvl + 1          AS lvl
  FROM path p
  JOIN storage_group sg
    ON sg.id     = p.parent_group_id
   AND sg.org_id = $1
)

SELECT id, type, alias, name
FROM (
  SELECT id, type, alias, name, lvl
  FROM path

  UNION ALL

  SELECT
    ou.id,
    'unit'            AS type,
    ou.alias,
    ou.name,
    MAX(p.lvl) + 1    AS lvl
  FROM path p
  JOIN org_unit ou
    ON ou.id     = p.unit_id
   AND ou.org_id = $1
  GROUP BY ou.id, ou.alias, ou.name
) t
ORDER BY lvl;

-- Items
-- name: CreateItem :one
INSERT INTO item (org_id, name, description) VALUES ($1, $2, $3) RETURNING *;

-- name: GetItemById :one
SELECT * FROM item WHERE org_id = $1 AND id = $2 AND deleted_at IS NULL;

-- name: GetItems :many
SELECT * FROM item WHERE org_id = $1 AND deleted_at IS NULL;

-- name: UpdateItem :one
UPDATE item SET name = $3, description = $4 WHERE org_id = $1 AND id = $2 AND deleted_at IS NULL RETURNING *;

-- name: DeleteItem :exec
UPDATE item SET deleted_at = CURRENT_TIMESTAMP WHERE org_id = $1 AND id = $2;

-- Item Variants
-- name: CreateItemVariant :one
INSERT INTO item_variant (org_id, item_id, name, article, ean13) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetItemVariantById :one
SELECT * FROM item_variant WHERE org_id = $1 AND item_id = $2 AND id = $3 AND deleted_at IS NULL;

-- name: GetItemVariants :many
SELECT * FROM item_variant WHERE org_id = $1 AND item_id = $2 AND deleted_at IS NULL;

-- name: UpdateItemVariant :one
UPDATE item_variant SET name = $4, article = $5, ean13 = $6 WHERE org_id = $1 AND item_id = $2 AND id = $3 AND deleted_at IS NULL RETURNING *;

-- name: DeleteItemVariant :exec
UPDATE item_variant SET deleted_at = CURRENT_TIMESTAMP WHERE org_id = $1 AND item_id = $2 AND id = $3;



-- User Auth
-- name: CreateUserSession :one
INSERT INTO app_user_session (user_id, token) VALUES ($1, $2) RETURNING *;

-- name: GetSessionBySecret :one
SELECT * FROM app_user_session WHERE token = $1 LIMIT 1;

-- name: InvalidateSession :exec
UPDATE app_user_session SET revoked_at = CURRENT_TIMESTAMP WHERE id = $1;

-- Item Instances
-- name: CreateItemInstance :one
INSERT INTO item_instance (org_id, item_id, variant_id, cell_id, status) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetItemInstancesForItem :many
SELECT * FROM item_instance WHERE org_id = $1 AND item_id = $2 AND deleted_at IS NULL;

-- name: GetItemInstancesForCell :many
SELECT * FROM item_instance WHERE org_id = $1 AND cell_id = $2 AND deleted_at IS NULL;

-- name: GetItemInstance :one
SELECT * FROM item_instance WHERE org_id = $1 AND id = $2 AND deleted_at IS NULL;

-- name: GetItemInstancesForCellsGroup :many
SELECT * FROM item_instance WHERE item_instance.org_id = $1 AND cell_id IN (SELECT id FROM cell WHERE cells_group_id = $2 AND deleted_at IS NULL) AND deleted_at IS NULL;


-- User
-- name: CreateUser :one
INSERT INTO app_user (email, first_name, last_name, middle_name, yandex_id) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetUserById :one
SELECT * FROM app_user WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM app_user WHERE email = $1 LIMIT 1;

-- RBAC
-- name: AssignRoleToUser :exec
INSERT INTO app_role_binding (role_id, user_id, org_id) 
VALUES ($1, $2, $3)
ON CONFLICT (user_id, org_id) 
DO UPDATE SET role_id = EXCLUDED.role_id;

-- name: UnassignRoleFromUser :exec
DELETE FROM app_role_binding WHERE org_id = $1 AND user_id = $2;

-- name: GetUserRoleInOrg :one
SELECT sqlc.embed(app_role) FROM app_role 
JOIN app_role_binding ON app_role.id = app_role_binding.role_id
WHERE app_role_binding.user_id = $2 AND app_role_binding.org_id = $1;

-- name: GetRoles :many
SELECT * FROM app_role;

-- name: GetRoleById :one
SELECT * FROM app_role WHERE id = $1;

-- name: GetEmployees :many
SELECT sqlc.embed(app_user), sqlc.embed(app_role) FROM app_user
JOIN app_role_binding ON app_user.id = app_role_binding.user_id
JOIN app_role ON app_role_binding.role_id = app_role.id
WHERE app_role_binding.org_id = $1;

-- name: GetEmployee :one
SELECT sqlc.embed(app_user), sqlc.embed(app_role) FROM app_user
JOIN app_role_binding ON app_user.id = app_role_binding.user_id
JOIN app_role ON app_role_binding.role_id = app_role.id
WHERE app_role_binding.org_id = $1 AND app_role_binding.user_id = $2;

-- Audit Log
-- name: GetObjectTypeById :one
SELECT * FROM object_type WHERE id = $1;

-- name: CreateObjectChange :one
INSERT INTO app_object_change (org_id, user_id, action, target_object_type, target_object_id, prechange_state, postchange_state) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetObjectChanges :many
SELECT * FROM app_object_change WHERE org_id = $1 AND target_object_type = $2 AND target_object_id = $3;


-- name: GetEmployeeByUserId :one
SELECT 
    u.id as user_id,
    u.email,
    u.first_name,
    u.last_name,
    u.middle_name,
    rb.role_id
FROM app_user u
JOIN app_role_binding rb ON rb.user_id = u.id
WHERE rb.org_id = $1 AND rb.user_id = $2;


-- Api Tokens
-- name: GetApiTokens :many
SELECT * FROM app_api_token WHERE org_id = $1 AND revoked_at IS NULL;

-- name: CreateApiToken :one
INSERT INTO app_api_token (org_id, name, token) VALUES ($1, $2, $3) RETURNING *;

-- name: RevokeApiToken :exec
UPDATE app_api_token SET revoked_at = CURRENT_TIMESTAMP WHERE org_id = $1 AND id = $2;

-- name: GetOrgIdByApiToken :one
SELECT org_id FROM app_api_token WHERE token = $1 AND revoked_at IS NULL;


-- Tasks
-- name: CreateTask :one
INSERT INTO task (org_id, unit_id, type, name, description, assigned_to_user_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: CreateTaskItem :one
INSERT INTO task_item (org_id, task_id, item_instance_id, source_cell_id, destination_cell_id) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetTasks :many
SELECT * FROM task WHERE org_id = $1;
