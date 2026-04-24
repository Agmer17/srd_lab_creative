-- name: CreateProgress :one
INSERT INTO progresses (
    project_id,
    title,
    weight,
    project_member_id
) VALUES (
    sqlc.arg('project_id'),
    sqlc.arg('title'),
    sqlc.arg('weight'),
    sqlc.arg('project_member_id')::uuid

)
RETURNING *;

-- name: UpdateProgress :one
UPDATE progresses
SET
    title        = COALESCE(sqlc.narg('title'), title),
    weight       = COALESCE(sqlc.narg('weight'), weight),
    is_completed = COALESCE(sqlc.narg('is_completed'), is_completed),
    project_member_id = COALESCE(sqlc.narg('project_member_id')::uuid, project_member_id)
WHERE id         = $1
RETURNING *;

-- name: DeleteProgress :execrows
DELETE FROM progresses
WHERE id         = $1;

-- name: GetProgressByProject :many
SELECT 
    p.id AS progress_id,
    p.title,
    p.weight,
    p.is_completed,
    p.created_at AS progress_created_at,
    p.project_id,
    
   jsonb_build_object(
        'id', pm.id,
        'project_id', pm.project_id,
        'is_owner', pm.is_owner,
        'joined_at', pm.joined_at,
        'user', jsonb_build_object(
            'id', u.id,
            'full_name', u.full_name,
            'email', u.email,
            'gender', u.gender,
            'profile_picture', u.profile_picture,
            'global_role', u.global_role,
            'created_at', u.created_at,
            'updated_at', u.updated_at
        ),
        'role', jsonb_build_object(
            'id', r.id,
            'role_name', r.name,
            'created_at', r.created_at
        )
    ) AS project_member

FROM progresses p
LEFT JOIN project_members pm 
    ON p.project_member_id = pm.id 
    AND pm.left_at IS NULL

LEFT JOIN users u 
    ON pm.user_id = u.id 
    AND u.deleted_at IS NULL

LEFT JOIN roles r 
    ON pm.role_id = r.id

WHERE p.project_id = sqlc.arg('project_id')
ORDER BY p.created_at DESC;

-- name: GetTotalWeightByProject :one
SELECT COALESCE(SUM(weight), 0)::DECIMAL(5,2) AS total_weight
FROM progresses
WHERE project_id = $1;

-- name: GetProgressById :one
SELECT 
    p.id AS progress_id,
    p.title,
    p.weight,
    p.is_completed,
    p.created_at AS progress_created_at,
    p.project_id,
    
   jsonb_build_object(
        'id', pm.id,
        'project_id', pm.project_id,
        'is_owner', pm.is_owner,
        'joined_at', pm.joined_at,
        'user', jsonb_build_object(
            'id', u.id,
            'full_name', u.full_name,
            'email', u.email,
            'gender', u.gender,
            'profile_picture', u.profile_picture,
            'global_role', u.global_role,
            'created_at', u.created_at,
            'updated_at', u.updated_at
        ),
        'role', jsonb_build_object(
            'id', r.id,
            'role_name', r.name,
            'created_at', r.created_at
        )
    ) AS project_member

FROM progresses p
LEFT JOIN project_members pm 
    ON p.project_member_id = pm.id 
    AND pm.left_at IS NULL

LEFT JOIN users u 
    ON pm.user_id = u.id 
    AND u.deleted_at IS NULL

LEFT JOIN roles r 
    ON pm.role_id = r.id

WHERE p.id = sqlc.arg('id')
ORDER BY p.created_at DESC;