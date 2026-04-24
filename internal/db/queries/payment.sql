-- name: GetLatestPayment :one
SELECT p.* FROM payments p
JOIN orders o ON o.id = p.order_id
WHERE o.user_id = $1 
  AND p.order_id = $2
  AND p.deleted_at IS NULL
ORDER BY p.created_at DESC 
LIMIT 1;


-- name: CreateNewPayment :one
INSERT INTO payments (
    order_id,
    method,
    amount
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: UpdatePaymentWithGatewayData :one
UPDATE payments 
SET 
    fee = $2,
    total_payment = $3,
    payment_number = $4,
    expired_at = $5
WHERE id = $1
RETURNING *;

-- name: SetPaymentExpired :exec
UPDATE payments 
SET 
    status = 'expired'
WHERE id = $1;
