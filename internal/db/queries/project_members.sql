-- name: AddProjectMember :one
INSERT INTO project_members (
    project_id,
    user_id,
    role_id
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: RemoveProjectMember :exec
UPDATE project_members
SET left_at = CURRENT_TIMESTAMP
WHERE project_id = $1
  AND user_id    = $2
  AND left_at IS NULL;

-- name: GetActiveProjectMembers :many
-- Untuk kebutuhan internal (cek member aktif, validasi, dsb)
SELECT
    pm.id,
    pm.project_id,
    pm.joined_at,
    u.id             AS user_id,
    u.full_name,
    u.email,
    u.gender,
    u.profile_picture,
    r.id             AS role_id,
    r.name           AS role_name
FROM project_members pm
INNER JOIN users u ON u.id = pm.user_id
    AND u.deleted_at IS NULL
INNER JOIN roles r ON r.id = pm.role_id
WHERE pm.project_id = $1
  AND pm.left_at    IS NULL
ORDER BY pm.joined_at ASC;