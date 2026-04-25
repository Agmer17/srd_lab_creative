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

-- name: RemoveProjectMember :execrows
UPDATE project_members
SET left_at = CURRENT_TIMESTAMP
WHERE id = $1
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
        'profile_picture', u.profile_picture,
        'global_role', u.global_role,
        'created_at', u.created_at,
        'updated_at', u.updated_at
    ) AS user,

     jsonb_build_object(
        'id', r.id,
        'role_name', r.name,
        'created_at', r.created_at
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
        'role_name', r.name,
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

-- name: UpdateProjectMemberRole :one
UPDATE project_members
SET 
    role_id = sqlc.narg('role_id'),
    is_owner =  COALESCE(sqlc.narg('is_owner'), is_owner)
WHERE id = sqlc.arg('member_id')
  AND left_at IS NULL
RETURNING *;

-- name: GetProjectMemberByID :one
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
        'role_name', r.name,
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


-- name: GetMemberDataByUserId :one
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
        'role_name', r.name,
        'created_at', r.created_at
    ) AS role

FROM project_members pm
JOIN users u 
    ON u.id = pm.user_id
   AND u.deleted_at IS NULL

JOIN roles r 
    ON r.id = pm.role_id

WHERE pm.user_id = sqlc.arg('user_id')
and pm.project_id = sqlc.arg('project_id')
LIMIT 1;

-- name: GetAllMember :many
SELECT 
    pm.id,
    pm.project_id,
    pm.user_id,
    pm.role_id,
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
        'role_name', r.name,
        'created_at', r.created_at
    ) AS role

FROM project_members pm
JOIN users u 
    ON u.id = pm.user_id
   AND u.deleted_at IS NULL

JOIN roles r 
    ON r.id = pm.role_id;