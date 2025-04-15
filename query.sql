-- name: GetOrgs :many
SELECT * FROM org WHERE deleted_at IS NULL;

-- name: IsOrgExistsById :one
SELECT EXISTS (SELECT 1 FROM org WHERE id = $1);

-- name: IsOrgExists :one
SELECT EXISTS (SELECT 1 FROM org WHERE name = $1 OR subdomain = $2);

-- name: GetOrgById :one
SELECT * FROM org WHERE id = $1 AND deleted_at IS NULL;

-- name: CreateOrg :one
INSERT INTO org (name, subdomain) VALUES ($1, $2) RETURNING *;

-- name: UpdateOrg :one
UPDATE org SET name = $2, subdomain = $3 WHERE id = $1 RETURNING *;

-- name: DeleteOrg :exec
UPDATE org SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1;

-- -- Units

-- name: GetOrganizationUnits :many
SELECT * FROM org_unit WHERE org_id = $1 AND deleted_at IS NULL;

-- name: IsOrganizationUnitExistsForOrganization :one
SELECT EXISTS (SELECT 1 FROM org_unit WHERE org_id = $1 AND id = $2 AND deleted_at IS NULL);

-- name: GetOrganizationUnitById :one
SELECT * FROM org_unit WHERE id = $1 AND deleted_at IS NULL;

-- name: CreateOrganizationUnit :one
INSERT INTO org_unit (org_id, name, alias, address) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateOrganizationUnit :one
UPDATE org_unit SET name = $2, alias = $3, address = $4 WHERE id = $1 AND deleted_at IS NULL RETURNING *;

-- name: DeleteOrganizationUnit :exec
UPDATE org_unit SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1;

-- --- Storage spaces
-- name: GetOrganizationStorageGroups :many
SELECT * FROM storage_space WHERE org_id = $1 AND deleted_at IS NULL;

-- name: IsStorageGroupExistsForOrganization :one
SELECT EXISTS (SELECT 1 FROM storage_space WHERE org_id = $1 AND id = $2 AND deleted_at IS NULL);

-- name: GetStorageGroupById :one
SELECT * FROM storage_space WHERE id = $1 AND deleted_at IS NULL;

-- name: CreateStorageGroup :one
INSERT INTO storage_space (org_id, unit_id, parent_id, name, alias) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdateStorageGroup :one
UPDATE storage_space SET name = $2, alias = $3 WHERE id = $1 AND deleted_at IS NULL RETURNING *;

-- name: DeleteStorageGroup :exec
UPDATE storage_space SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1;
