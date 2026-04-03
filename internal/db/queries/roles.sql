-- name: CreateRole :one
INSERT INTO roles (name)
VALUES (sqlc.arg(name))
RETURNING *;
 
-- name: GetRoleById :one
SELECT * FROM roles
WHERE id = sqlc.arg(id);

-- name: GetAllRoles :many
SELECT * FROM roles
ORDER BY created_at DESC;

-- name: GetRoleByName :one
SELECT * FROM roles
WHERE name = sqlc.arg(name);
 
-- name: SearchRoles :many
SELECT * FROM roles
WHERE name ILIKE '%' || sqlc.arg(keyword) || '%'
ORDER BY created_at DESC
LIMIT sqlc.arg(limit_val)
OFFSET sqlc.arg(offset_val);

-- name: UpdateRole :one
UPDATE roles
SET name = sqlc.arg(name)
WHERE id = sqlc.arg(id)
RETURNING *;
 
-- name: DeleteRole :exec
DELETE FROM roles
WHERE id = sqlc.arg(id);


