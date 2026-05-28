-- name: CreateReview :one
INSERT INTO reviews (
    user_id,
    product_id,
    order_id,
    rating,
    comment,
    show
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetReviewByOrderID :one
SELECT * FROM reviews
WHERE order_id = $1;

-- name: ListReviewsByProduct :many
SELECT 
    sqlc.embed(r),
    sqlc.embed(u)
FROM reviews r
JOIN users u ON u.id = r.user_id
WHERE r.product_id = $1
  AND r.show = TRUE
ORDER BY r.created_at DESC;

-- name: ListFeaturedReviews :many
SELECT 
    sqlc.embed(r),
    sqlc.embed(u),
    sqlc.embed(p)
FROM reviews r
JOIN users u ON u.id = r.user_id
JOIN products p ON p.id = r.product_id
WHERE r.show = TRUE
ORDER BY r.rating DESC, r.created_at DESC
LIMIT $1;

-- name: UpdateReviewShowStatus :one
UPDATE reviews
SET show = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: ListReviewsByUser :many
SELECT
    sqlc.embed(r),
    sqlc.embed(p)
FROM reviews r
JOIN products p ON p.id = r.product_id
WHERE r.user_id = $1
ORDER BY r.created_at DESC;
