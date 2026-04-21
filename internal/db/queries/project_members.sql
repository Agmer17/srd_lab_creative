-- name: AddProjectMember :one
INSERT INTO project_members (
    project_id,
    user_id,
    role_id,
    is_owner
) VALUES (
    sqlc.arg('project_id') ,sqlc.arg('user_id'), sqlc.arg('role_id'),sqlc.arg('is_owner')
)
RETURNING *;

-- name: RemoveProjectMember :exec
UPDATE project_members
SET left_at = CURRENT_TIMESTAMP
WHERE project_id = $1
  AND user_id    = $2
  AND left_at IS NULL;

-- name: GetActiveProjectMembers :many
SELECT 
    pm.id,
    pm.project_id,
    pm.is_owner,
    pm.joined_at,

    jsonb_build_object(
        'id', u.id,
        'full_name', u.full_name,
        'email', u.email,
        'gender', u.gender,
        'profile_picture', u.profile_picture
    ) AS user,

    jsonb_build_object(
        'id', r.id,
        'name', r.name
    ) AS role

FROM project_members pm
JOIN users u 
    ON u.id = pm.user_id
   AND u.deleted_at IS NULL

JOIN roles r 
    ON r.id = pm.role_id

WHERE pm.project_id = sqlc.arg('project_id')
  AND pm.left_at IS NULL

ORDER BY pm.joined_at ASC;

-- name: GetProjectMemberWithUser :one
SELECT 
    pm.id,
    pm.project_id,
    pm.is_owner,
    pm.joined_at,
    pm.left_at,

    jsonb_build_object(
        'id', u.id,
        'full_name', u.full_name,
        'email', u.email,
        'gender', u.gender,
        'profile_picture', u.profile_picture,
        'global_role', u.global_role,
        'created_at', u.created_at,
        'updated_at', u.updated_at
    ) AS user,

    jsonb_build_object(
        'id', r.id,
        'name', r.name,
        'created_at', r.created_at
    ) AS role

FROM project_members pm
JOIN users u 
    ON u.id = pm.user_id
   AND u.deleted_at IS NULL

JOIN roles r 
    ON r.id = pm.role_id

WHERE pm.id = sqlc.arg('id')
LIMIT 1;

-- name: ListProjectMembersWithUser :many
SELECT 
    pm.id,
    pm.project_id,
    pm.is_owner,
    pm.joined_at,
    pm.left_at,

    jsonb_build_object(
        'id', u.id,
        'full_name', u.full_name,
        'email', u.email,
        'gender', u.gender,
        'profile_picture', u.profile_picture,
        'global_role', u.global_role,
        'created_at', u.created_at,
        'updated_at', u.updated_at
    ) AS user,

    jsonb_build_object(
        'id', r.id,
        'name', r.name,
        'created_at', r.created_at
    ) AS role

FROM project_members pm
JOIN users u 
    ON u.id = pm.user_id
   AND u.deleted_at IS NULL

JOIN roles r 
    ON r.id = pm.role_id

WHERE pm.project_id = sqlc.arg('project_id')
ORDER BY pm.joined_at ASC;