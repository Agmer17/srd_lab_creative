-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = sqlc.arg(email)
AND deleted_at IS NULL;
 
-- name: CreateUser :one
INSERT INTO users (
    full_name,
    email,
    profile_picture,
    provider,
    provider_user_id
) VALUES (
    sqlc.arg(full_name),
    sqlc.arg(email),
    sqlc.arg(profile_picture),
    sqlc.arg(provider),
    sqlc.arg(provider_user_id)
)
RETURNING *;
 
-- name: GetUserAuthInfoByProviderID :one
SELECT 
    id,
    global_role
FROM users
WHERE provider = $1
  AND provider_user_id = $2
LIMIT 1;
 
-- name: GetUserById :one
SELECT * FROM users
WHERE id = sqlc.arg(id)
AND deleted_at IS NULL;
 
-- name: ListUsers :many
SELECT * FROM users
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT sqlc.arg(limit_val)
OFFSET sqlc.arg(offset_val);
 
-- name: SearchUsers :many
SELECT * FROM users
WHERE deleted_at IS NULL
  AND (
    full_name    ILIKE '%' || sqlc.arg(keyword) || '%'
    OR email     ILIKE '%' || sqlc.arg(keyword) || '%'
    OR phone_number ILIKE '%' || sqlc.arg(keyword) || '%'
  )
ORDER BY created_at DESC
LIMIT sqlc.arg(limit_val)::int
OFFSET sqlc.arg(offset_val)::int;

-- name: UpdateUser :one
UPDATE users
SET
    full_name       = COALESCE(sqlc.narg(full_name),      full_name),
    gender          = COALESCE(sqlc.narg(gender), gender),
    phone_number    = COALESCE(sqlc.narg(phone_number), phone_number),
    updated_at      = CURRENT_TIMESTAMP
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL
RETURNING *;

-- name: UpdateUserGlobalRole :one
UPDATE users
SET
    global_role = sqlc.arg(global_role),
    updated_at  = CURRENT_TIMESTAMP
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteUser :one
UPDATE users
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL
RETURNING *;
 
-- name: HardDeleteUser :exec
DELETE FROM users
WHERE id = sqlc.arg(id);