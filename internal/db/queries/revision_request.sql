-- name: CreateRevisionRequest :one
INSERT INTO revision_requests (
    project_id,
    title,
    reason
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: UpdateRevisionStatus :one
UPDATE revision_requests
SET 
    status = COALESCE(sqlc.narg('status'), status)
WHERE id   = $1
RETURNING *;

-- name: GetRevisionsByProject :many
SELECT *
FROM revision_requests
WHERE project_id = $1
ORDER BY created_at DESC;

-- name: CountPendingRevisions :one
-- Untuk validasi allowed_revision_count di application layer
SELECT COUNT(*)::INT AS pending_count
FROM revision_requests
WHERE project_id = $1
  AND status     = 'pending';