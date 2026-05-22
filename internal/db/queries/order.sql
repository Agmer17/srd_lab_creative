-- name: CreateOrder :one
INSERT INTO orders (
    user_id,
    product_id,
    ordered_price
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetOrderByID :one
SELECT 
    sqlc.embed(o),
    sqlc.embed(u),
    sqlc.embed(p),

    COALESCE(
        jsonb_agg(
            jsonb_build_object(
                'payment_id', pay.id,
                'order_id', pay.order_id,
                'method', pay.method,
                'status', pay.status,
                'amount', pay.amount,
                'fee', pay.fee,
                'total_payment', pay.total_payment,
                'payment_number', pay.payment_number,
                'expired_at', pay.expired_at,
                'paid_at', pay.paid_at,
                'created_at', pay.created_at
            )
        ) FILTER (WHERE pay.id IS NOT NULL),
        '[]'
    )::jsonb AS payments

FROM orders o
JOIN users u ON u.id = o.user_id
JOIN products p ON p.id = o.product_id
LEFT JOIN payments pay 
    ON pay.order_id = o.id
   AND pay.deleted_at IS NULL

WHERE o.id = $1
  AND o.deleted_at IS NULL

GROUP BY o.id, u.id, p.id;

-- name: ListOrdersByUser :many
SELECT 
    sqlc.embed(o),
    sqlc.embed(u),
    sqlc.embed(p),
    -- Subquery ringkas (hanya 4 field) buat list order
    COALESCE(
        (
            SELECT jsonb_agg(
                jsonb_build_object(
                    'payment_id', pay.id,
                    'method', pay.method,
                    'status', pay.status,
                    'amount', pay.amount
                )
            )
            FROM payments pay
            WHERE pay.order_id = o.id
              AND pay.deleted_at IS NULL
        ),
        '[]'
    )::jsonb AS payments
FROM orders o
JOIN users u ON u.id = o.user_id
JOIN products p ON p.id = o.product_id
WHERE o.user_id = @user_id
  AND (
    sqlc.narg('status')::text IS NULL 
    OR sqlc.narg('status')::text = ''
    OR o.status = sqlc.narg('status')::order_status_enum
  )
  AND o.deleted_at IS NULL
ORDER BY o.created_at DESC;

-- name: UpdateOrderStatus :one
UPDATE orders
SET status = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
  AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteOrder :execrows
UPDATE orders
SET deleted_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
  AND deleted_at IS NULL;

-- name: ListOrders :many
SELECT 
    sqlc.embed(o),
    sqlc.embed(u),
    sqlc.embed(p),
    -- Subquery ringkas (hanya 4 field) buat list order
    COALESCE(
        (
            SELECT jsonb_agg(
                jsonb_build_object(
                    'payment_id', pay.id,
                    'method', pay.method,
                    'status', pay.status,
                    'amount', pay.amount
                )
            )
            FROM payments pay
            WHERE pay.order_id = o.id
              AND pay.deleted_at IS NULL
        ),
        '[]'
    )::jsonb AS payments
FROM orders o
JOIN users u ON u.id = o.user_id
JOIN products p ON p.id = o.product_id
WHERE (
   sqlc.narg('status')::text IS NULL 
   OR sqlc.narg('status')::text = ''
   OR o.status = sqlc.narg('status')::order_status_enum
)
AND o.deleted_at IS NULL
ORDER BY o.created_at DESC;