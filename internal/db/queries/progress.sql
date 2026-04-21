-- name: CreateProgress :one
INSERT INTO progresses (
    project_id,
    title,
    weight
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: UpdateProgress :one
UPDATE progresses
SET
    title        = COALESCE(sqlc.narg('title'), title),
    weight       = COALESCE(sqlc.narg('weight'), weight),
    is_completed = COALESCE(sqlc.narg('is_completed'), is_completed)
WHERE id         = $1
  AND project_id = $2
RETURNING *;

-- name: DeleteProgress :exec
DELETE FROM progresses
WHERE id         = $1
  AND project_id = $2;

-- name: GetProgressesByProject :many
SELECT *
FROM progresses
WHERE project_id = $1
ORDER BY created_at ASC;

-- name: GetTotalWeightByProject :one
-- Untuk validasi di application layer agar total weight = 100
SELECT COALESCE(SUM(weight), 0)::DECIMAL(5,2) AS total_weight
FROM progresses
WHERE project_id = $1;