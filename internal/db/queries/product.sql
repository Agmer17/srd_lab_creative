-- name: CreateProduct :one
INSERT INTO products(
    name,slug,description,price,status,is_featured
) VALUES(
    $1,$2,$3,$4,$5,$6
)
RETURNING *;

-- name: GetProductById :one
SELECT * FROM products 
WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: GetAllProduct :many
SELECT * FROM products
WHERE deleted_at IS NULL
ORDER BY created_at DESC;

-- name: DeleteProduct :exec
UPDATE products
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateProduct :one
UPDATE products
SET 
    name = COALESCE(sqlc.narg('name'), name),
    slug = COALESCE(sqlc.narg('slug'), slug),
    description = COALESCE(sqlc.narg('description'), description),
    price = COALESCE(sqlc.narg('price'), price),
    status = COALESCE(sqlc.narg('status'), status),
    is_featured = COALESCE(sqlc.narg('is_featured'), is_featured),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;


-- name: CheckProductSlugExists :one
SELECT EXISTS(
    SELECT 1 FROM products 
    WHERE slug = $1 AND deleted_at IS NULL
);

-- name: GetProductBySlug :one
SELECT * FROM products 
WHERE slug = $1 AND deleted_at IS NULL LIMIT 1;
