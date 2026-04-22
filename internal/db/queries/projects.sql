-- name: CreateProject :one
INSERT INTO projects (
    order_id,
    name,
    description,
    status,
    allowed_revision_count,
    actual_start_date,
    end_date
) VALUES (
    sqlc.arg('order_id'), 
    sqlc.arg('name'), 
    sqlc.narg('description'), 
    sqlc.arg('status'), 
    COALESCE(sqlc.narg('allowed_revision_count')::int, 3), 
    sqlc.narg('start_date'), -- Gunakan narg jika start_date bisa null/tidak diisi saat create
    sqlc.narg('end_date')
)
RETURNING *;

-- name: UpdateProject :one
UPDATE projects
SET
    name                   = COALESCE(sqlc.narg('name'), name),
    description            = COALESCE(sqlc.narg('description'), description),
    status                 = COALESCE(sqlc.narg('status'), status),
    allowed_revision_count = COALESCE(sqlc.narg('allowed_revision_count'), allowed_revision_count),
    actual_start_date      = COALESCE(sqlc.narg('actual_start_date'), actual_start_date),
    end_date               = COALESCE(sqlc.narg('end_date'), end_date),
    updated_at             = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteProject :execrows
DELETE FROM projects WHERE id = $1;

-- name: ListProjects :many
SELECT
    p.id,
    p.order_id,
    p.name,
    p.description,
    p.status,
    p.allowed_revision_count,
    p.actual_start_date,
    p.end_date,
    p.created_at,
    p.updated_at,

    -- Project Members
    COALESCE(
        (
            SELECT JSON_AGG(JSONB_BUILD_OBJECT(
                'id', pm.id,
                'project_id', pm.project_id,
                'joined_at', pm.joined_at,
                'left_at', pm.left_at,
                'is_owner', pm.is_owner,

                'user', JSONB_BUILD_OBJECT(
                    'id', u.id,
                    'global_role', u.global_role,
                    'full_name', u.full_name,
                    'email', u.email,
                    'phone_number', u.phone_number,
                    'profile_picture', u.profile_picture,
                    'gender', u.gender,
                    'created_at', u.created_at,
                    'updated_at', u.updated_at,
                    'deleted_at', u.deleted_at
                ),

                'role', JSONB_BUILD_OBJECT(
                    'id', r.id,
                    'role_name', r.name,
                    'created_at', r.created_at
                )
            ))
            FROM project_members pm
            JOIN users u ON u.id = pm.user_id
            JOIN roles r ON r.id = pm.role_id
            WHERE pm.project_id = p.id AND pm.left_at IS NULL
        ),
        '[]'
    )::jsonb AS project_members,

    -- Progress
    COALESCE(
        (
            SELECT JSON_AGG(JSONB_BUILD_OBJECT(
                'id', pr.id,
                'project_id', pr.project_id,
                'title', pr.title,
                'weight', pr.weight,
                'is_completed', pr.is_completed,
                'created_at', pr.created_at,

                'member', JSONB_BUILD_OBJECT(
                    'id', pm_task.id,
                    'project_id', pm_task.project_id,
                    'joined_at', pm_task.joined_at,
                    'left_at', pm_task.left_at,
                    'is_owner', pm_task.is_owner,

                    'user', JSONB_BUILD_OBJECT(
                        'id', u_task.id,
                        'global_role', u_task.global_role,
                        'full_name', u_task.full_name,
                        'email', u_task.email,
                        'phone_number', u_task.phone_number,
                        'profile_picture', u_task.profile_picture,
                        'gender', u_task.gender,
                        'created_at', u_task.created_at,
                        'updated_at', u_task.updated_at,
                        'deleted_at', u_task.deleted_at
                    ),

                    'role', JSONB_BUILD_OBJECT(
                        'id', r_task.id,
                        'role_name', r_task.name,
                        'created_at', r_task.created_at
                    )
                )
            ))
            FROM progresses pr
            JOIN project_members pm_task ON pm_task.id = pr.project_member_id
            JOIN users u_task ON u_task.id = pm_task.user_id
            JOIN roles r_task ON r_task.id = pm_task.role_id
            WHERE pr.project_id = p.id
        ),
        '[]'
    )::jsonb AS progress

FROM projects p
ORDER BY p.created_at DESC;

-- ---------------------------------------------------------------
-- GET DETAIL PROJECT
-- Join: orders, project_members -> users -> roles,
--       progresses, revision_requests
-- Members filter: hanya yang left_at IS NULL
-- ---------------------------------------------------------------
-- name: GetProjectDetail :one
SELECT
    p.id,
    p.name,
    p.description,
    p.status,
    p.allowed_revision_count,
    p.actual_start_date,
    p.end_date,
    p.created_at,
    p.updated_at,

    -- Order (tetap flat)
    o.id            AS order_id,
    o.status        AS order_status,
    o.ordered_price AS ordered_price,
    o.user_id       AS order_user_id,
    o.product_id    AS order_product_id,
    o.created_at    AS order_created_at,

    COALESCE(
        (
            SELECT JSON_AGG(JSONB_BUILD_OBJECT(
                'id', pm.id,
                'project_id', pm.project_id,
                'joined_at', pm.joined_at,
                'left_at', pm.left_at,
                'is_owner', pm.is_owner,

                'user', JSONB_BUILD_OBJECT(
                    'id', u.id,
                    'global_role', u.global_role,
                    'full_name', u.full_name,
                    'email', u.email,
                    'phone_number', u.phone_number,
                    'profile_picture', u.profile_picture,
                    'gender', u.gender,
                    'created_at', u.created_at,
                    'updated_at', u.updated_at,
                    'deleted_at', u.deleted_at
                ),

                'role', JSONB_BUILD_OBJECT(
                    'id', r.id,
                    'role_name', r.name,
                    'created_at', r.created_at
                )
            ))
            FROM project_members pm
            JOIN users u ON u.id = pm.user_id
            JOIN roles r ON r.id = pm.role_id
            WHERE pm.project_id = p.id AND pm.left_at IS NULL
        ),
        '[]'
    )::jsonb AS project_members,

    COALESCE(
        (
            SELECT JSON_AGG(JSONB_BUILD_OBJECT(
                'id', pr.id,
                'project_id', pr.project_id,
                'title', pr.title,
                'weight', pr.weight,
                'is_completed', pr.is_completed,
                'created_at', pr.created_at,

                'member', JSONB_BUILD_OBJECT(
                    'id', pm_task.id,
                    'project_id', pm_task.project_id,
                    'joined_at', pm_task.joined_at,
                    'left_at', pm_task.left_at,
                    'is_owner', pm_task.is_owner,

                    'user', JSONB_BUILD_OBJECT(
                        'id', u_task.id,
                        'global_role', u_task.global_role,
                        'full_name', u_task.full_name,
                        'email', u_task.email,
                        'phone_number', u_task.phone_number,
                        'profile_picture', u_task.profile_picture,
                        'gender', u_task.gender,
                        'created_at', u_task.created_at,
                        'updated_at', u_task.updated_at,
                        'deleted_at', u_task.deleted_at
                    ),

                    'role', JSONB_BUILD_OBJECT(
                        'id', r_task.id,
                        'role_name', r_task.name,
                        'created_at', r_task.created_at
                    )
                )
            ))
            FROM progresses pr
            JOIN project_members pm_task ON pm_task.id = pr.project_member_id
            JOIN users u_task ON u_task.id = pm_task.user_id
            JOIN roles r_task ON r_task.id = pm_task.role_id
            WHERE pr.project_id = p.id
        ),
        '[]'
    )::jsonb AS progress,

    COALESCE(
        (
            SELECT JSON_AGG(JSONB_BUILD_OBJECT(
                'id', rv.id,
                'title', rv.title,
                'reason', rv.reason,
                'status', rv.status,
                'created_at', rv.created_at
            ))
            FROM revision_requests rv
            WHERE rv.project_id = p.id
        ),
        '[]'
    )::jsonb AS revisions

FROM projects p
JOIN orders o ON o.id = p.order_id

WHERE p.id = $1;