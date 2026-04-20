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
    $1, $2, $3, $4, $5, $6, $7
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

-- name: DeleteProject :exec
DELETE FROM projects WHERE id = $1;

-- ---------------------------------------------------------------
-- GET LIST PROJECTS
-- Join: project_members -> users -> roles, progresses
-- Members filter: hanya yang left_at IS NULL (masih aktif)
-- ---------------------------------------------------------------
-- name: ListProjects :many
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

    -- Members: JSON array of active members with user & role info
    COALESCE(
        JSON_AGG(
            DISTINCT JSONB_BUILD_OBJECT(
                'id',              pm.id,
                'user_id',         u.id,
                'full_name',       u.full_name,
                'email',           u.email,
                'gender',          u.gender,
                'profile_picture', u.profile_picture,
                'role_id',         r.id,
                'role_name',       r.name,
                'joined_at',       pm.joined_at
            )
        ) FILTER (WHERE pm.id IS NOT NULL),
        '[]'
    )::jsonb AS members,

    -- Progresses: JSON array
    COALESCE(
        JSON_AGG(
            DISTINCT JSONB_BUILD_OBJECT(
                'id',           pr.id,
                'title',        pr.title,
                'weight',       pr.weight,
                'is_completed', pr.is_completed,
                'created_at',   pr.created_at
            )
        ) FILTER (WHERE pr.id IS NOT NULL),
        '[]'
    )::jsonb AS progresses

FROM projects p

LEFT JOIN project_members pm ON pm.project_id = p.id
    AND pm.left_at IS NULL

LEFT JOIN users u  ON u.id  = pm.user_id
    AND u.deleted_at IS NULL

LEFT JOIN roles r  ON r.id  = pm.role_id

LEFT JOIN progresses pr ON pr.project_id = p.id

GROUP BY p.id
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

    -- Order info (flat columns, bukan nested — lebih predictable di scan)
    o.id            AS order_id,
    o.status        AS order_status,
    o.ordered_price AS ordered_price,
    o.user_id       AS order_user_id,
    o.product_id    AS order_product_id,
    o.created_at    AS order_created_at,

    -- Members
    COALESCE(
        JSON_AGG(
            DISTINCT JSONB_BUILD_OBJECT(
                'id',              pm.id,
                'user_id',         u.id,
                'full_name',       u.full_name,
                'email',           u.email,
                'gender',          u.gender,
                'profile_picture', u.profile_picture,
                'role_id',         r.id,
                'role_name',       r.name,
                'joined_at',       pm.joined_at
            )
        ) FILTER (WHERE pm.id IS NOT NULL),
        '[]'
    )::jsonb AS members,

    -- Progresses
    COALESCE(
        JSON_AGG(
            DISTINCT JSONB_BUILD_OBJECT(
                'id',           pr.id,
                'title',        pr.title,
                'weight',       pr.weight,
                'is_completed', pr.is_completed,
                'created_at',   pr.created_at
            )
        ) FILTER (WHERE pr.id IS NOT NULL),
        '[]'
    )::jsonb AS progresses,

    -- Revision requests
    COALESCE(
        JSON_AGG(
            DISTINCT JSONB_BUILD_OBJECT(
                'id',         rv.id,
                'title',      rv.title,
                'reason',     rv.reason,
                'status',     rv.status,
                'created_at', rv.created_at
            )
        ) FILTER (WHERE rv.id IS NOT NULL),
        '[]'
    )::jsonb AS revisions

FROM projects p

INNER JOIN orders o ON o.id = p.order_id

LEFT JOIN project_members pm ON pm.project_id = p.id
    AND pm.left_at IS NULL

LEFT JOIN users u ON u.id = pm.user_id
    AND u.deleted_at IS NULL

LEFT JOIN roles r ON r.id = pm.role_id

LEFT JOIN progresses pr ON pr.project_id = p.id

LEFT JOIN revision_requests rv ON rv.project_id = p.id

WHERE p.id = $1

GROUP BY
    p.id,
    o.id,
    o.status,
    o.ordered_price,
    o.user_id,
    o.product_id,
    o.created_at;
