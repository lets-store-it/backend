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

-- -- name: GetOrgUnits :many
-- SELECT * FROM org_unit WHERE org_id = $1 LIMIT $2 OFFSET $3;

-- -- name: GetUnitById :one
-- SELECT * FROM org_unit WHERE id = $1;

-- -- name: CreateUnit :one
-- INSERT INTO org_unit (org_id, name, address) VALUES ($1, $2, $3) RETURNING *;

-- -- name: UpdateUnit :one
-- UPDATE org_unit SET name = $2, address = $3 WHERE id = $1 RETURNING *;

-- -- name: DeleteUnit :exec
-- UPDATE org_unit SET is_deleted = TRUE WHERE id = $1;

-- --- Storage spaces

-- -- name: GetOrgStorageSpaces :many
-- SELECT * FROM storage_space WHERE org_id = $1 LIMIT $2 OFFSET $3;

-- -- name: GetStorageSpaceById :one
-- SELECT * FROM storage_space WHERE id = $1;

-- -- name: GetStorageSpaceByUnitId :many
-- SELECT * FROM storage_space WHERE unit_id = $1;

-- -- name: CreateStorageSpace :one
-- INSERT INTO storage_space (unit_id, name, short_name) VALUES ($1, $2, $3) RETURNING *;

-- -- name: UpdateStorageSpace :one
-- UPDATE storage_space SET name = $2, short_name = $3 WHERE id = $1 RETURNING *;

-- -- name: DeleteStorageSpace :exec
-- UPDATE storage_space SET is_deleted = TRUE WHERE id = $1;

-- --- name: GetStorageSpaceChain:
-- WITH RECURSIVE storage_space_chain AS (
--     SELECT id, name, short_name, parent_id, unit_id, org_id
--     FROM storage_space
--     WHERE id = $1
--     UNION ALL
--     SELECT s.id, s.name, s.short_name, s.parent_id, s.unit_id, s.org_id
--     FROM storage_space s
--     JOIN storage_space_chain c ON s.parent_id = c.id
-- )
-- SELECT * FROM storage_space_chain;


