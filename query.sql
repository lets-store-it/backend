-- name: GetOrgs :many
SELECT * FROM org WHERE is_deleted = FALSE LIMIT $1 OFFSET $2;

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


-- name: ExampleJoin :many
SELECT * FROM org
JOIN org_unit ON org.id = org_unit.org_id
WHERE org.id = $1;
