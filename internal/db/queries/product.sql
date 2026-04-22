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


-- name: CreateProductImage :many
INSERT INTO product_images (product_id, image_url, is_primary, sort_order)
SELECT $1, unnest($2::text[]), unnest($3::bool[]), unnest($4::int[])
RETURNING *;


-- name: GetProductImageByImageId :one
SELECT * FROM product_images
WHERE id = $1 LIMIT 1;

-- name: GetAllProductImageByProductId :many
SELECT * FROM product_images
WHERE product_id = $1 ORDER BY sort_order ASC;

-- name: DeleteProductImageByImageId :exec
DELETE FROM product_images
WHERE id = $1;

-- name: DeleteProductImageByProductId :exec
DELETE FROM product_images
WHERE product_id = $1;

-- name: UpdateProductImageByImageId :one
UPDATE product_images
SET
    product_id = COALESCE(sqlc.narg('product_id'), product_id),
    image_url = COALESCE(sqlc.narg('image_url'), image_url),
    is_primary = COALESCE(sqlc.narg('is_primary'), is_primary),
    sort_order = COALESCE(sqlc.narg('sort_order'), sort_order)
WHERE id = $1
RETURNING *;

-- name: GetTotalImageOfProductId :one
SELECT COUNT(*) FROM product_images
WHERE product_id = $1;

-- name: GetImageIdsAndOrderByProductId :many
SELECT id, sort_order FROM product_images
WHERE product_id = $1 ORDER BY sort_order ASC;

-- name: ImageIdOrderChange :exec
UPDATE product_images
SET sort_order = t.new_order
FROM (
    SELECT unnest($1::uuid[]) AS image_id,
           unnest($2::int[])  AS new_order
) AS t
WHERE product_images.id = t.image_id;



-- name: AssignProductToCategory :exec
INSERT INTO product_categories(product_id, category_id)
VALUES ($1, $2)
ON CONFLICT (product_id, category_id) DO NOTHING;

-- name: RemoveProductFromCategory :exec
DELETE FROM product_categories
WHERE product_id = $1 AND category_id = $2;

-- name: RemoveProductFromAllCategory :exec
DELETE FROM product_categories
WHERE product_id = $1;

-- name: GetProductCategory :many
SELECT c.* FROM categories c
JOIN product_categories pc ON c.id = pc.category_id
WHERE pc.product_id = $1;

-- name: GetProductsByCategory :many
SELECT p.* FROM products p
JOIN product_categories pc ON p.id = pc.product_id
WHERE pc.category_id = $1 AND p.deleted_at IS NULL;

