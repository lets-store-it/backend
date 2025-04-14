-- name: GetOrgs :many
SELECT * FROM org WHERE is_deleted = FALSE;

-- name: IsOrgExistsById :one
SELECT EXISTS (SELECT 1 FROM org WHERE id = $1);

-- name: IsOrgExists :one
SELECT EXISTS (SELECT 1 FROM org WHERE name = $1 OR subdomain = $2);

-- name: GetOrgById :one
SELECT * FROM org WHERE id = $1 AND is_deleted = FALSE;

-- name: CreateOrg :one
INSERT INTO org (name, subdomain) VALUES ($1, $2) RETURNING *;

-- name: UpdateOrg :one
UPDATE org SET name = $2, subdomain = $3 WHERE id = $1 RETURNING *;

-- name: DeleteOrg :exec
UPDATE org SET is_deleted = TRUE WHERE id = $1;

-- -- Units

-- name: GetOrganizationUnits :many
SELECT * FROM org_unit WHERE org_id = $1 AND is_deleted = FALSE;

-- name: IsOrganizationUnitExistsForOrganization :one
SELECT EXISTS (SELECT 1 FROM org_unit WHERE org_id = $1 AND id = $2 AND is_deleted = FALSE);

-- name: GetOrganizationUnitById :one
SELECT * FROM org_unit WHERE id = $1 AND is_deleted = FALSE;

-- name: CreateOrganizationUnit :one
INSERT INTO org_unit (org_id, name, alias, address) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateOrganizationUnit :one
UPDATE org_unit SET name = $2, alias = $3, address = $4 WHERE id = $1 AND is_deleted = FALSE RETURNING *;

-- name: DeleteOrganizationUnit :exec
UPDATE org_unit SET is_deleted = TRUE WHERE id = $1;

-- --- Storage spaces
-- -- name: GetOrganizationStorageSpaces :many
SELECT * FROM storage_space WHERE org_id = $1 AND is_deleted = FALSE;

-- name: IsStorageSpaceExistsForOrganization :one
SELECT EXISTS (SELECT 1 FROM storage_space WHERE org_id = $1 AND id = $2 AND is_deleted = FALSE);

-- name: GetStorageSpaceById :one
SELECT * FROM storage_space WHERE id = $1 AND is_deleted = FALSE;

-- name: CreateStorageSpace :one
INSERT INTO storage_space (org_id, unit_id, parent_id, name, alias) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdateStorageSpace :one
UPDATE storage_space SET name = $2, alias = $3 WHERE id = $1 AND is_deleted = FALSE RETURNING *;

-- name: DeleteStorageSpace :exec
UPDATE storage_space SET is_deleted = TRUE WHERE id = $1;