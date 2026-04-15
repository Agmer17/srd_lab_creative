-- name: CreateCategory :one
INSERT INTO categories (
    name,
    slug,
    description
) VALUES (
    sqlc.arg(name),
    sqlc.arg(slug),
    sqlc.narg(description)
)
RETURNING *;

-- name: GetCategoryByID :one
SELECT *
FROM categories
WHERE id = sqlc.arg(id)
LIMIT 1;

-- name: GetCategoryBySlug :one
SELECT *
FROM categories
WHERE slug = sqlc.arg(slug)
LIMIT 1;

-- name: ListCategories :many
SELECT *
FROM categories
ORDER BY name ASC
LIMIT sqlc.arg(page_limit)
OFFSET sqlc.arg(page_offset);

-- name: UpdateCategory :one
UPDATE categories
SET
    name        = COALESCE(sqlc.narg(name), name),
    slug        = COALESCE(sqlc.narg(slug), slug),
    description = COALESCE(sqlc.narg(description), description)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = sqlc.arg(id);

-- name: SearchCategories :many
SELECT *
FROM categories
WHERE
    name ILIKE '%' || sqlc.arg(query) || '%'
    OR description ILIKE '%' || sqlc.arg(query) || '%'
ORDER BY name ASC
LIMIT  sqlc.arg(page_limit)
OFFSET sqlc.arg(page_offset);

-- name: ExistCategoryBySlug :one
SELECT EXISTS (
    SELECT 1
    FROM categories
    WHERE slug = sqlc.arg(slug)
) AS exist;