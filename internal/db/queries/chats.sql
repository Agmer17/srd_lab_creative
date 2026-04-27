-- name: CreateChat :one
INSERT INTO chats (room_id, sender_id, text)
VALUES (
    sqlc.arg('room_id'),
    sqlc.arg('sender_id'),
    sqlc.arg('text')
)
RETURNING *;

-- name: GetChatsByRoomID :many
SELECT
    c.id,
    c.room_id,
    c.text,
    c.created_at,
    jsonb_build_object(
        'id',              u.id,
        'full_name',       u.full_name,
        'email',           u.email,
        'gender',          u.gender,
        'profile_picture', u.profile_picture,
        'global_role',     u.global_role,
        'created_at',      u.created_at,
        'updated_at',      u.updated_at
    ) AS sender,
    COALESCE(
        jsonb_agg(
            jsonb_build_object(
                'id',         cm.id,
                'file_name',  cm.file_name,
                'media_type', cm.media_type,
                'size',       cm.size,
                'created_at', cm.created_at
            )
        ) FILTER (WHERE cm.id IS NOT NULL),
        '[]'::jsonb
    )::jsonb AS medias
FROM chats c
LEFT JOIN users u
    ON u.id = c.sender_id
   AND u.deleted_at IS NULL
LEFT JOIN chat_medias cm
    ON cm.chat_id = c.id
WHERE c.room_id = sqlc.arg('room_id')
GROUP BY
    c.id,
    c.room_id,
    c.text,
    c.created_at,
    u.id,
    u.full_name,
    u.email,
    u.gender,
    u.profile_picture,
    u.global_role,
    u.created_at,
    u.updated_at
ORDER BY c.created_at ASC;

-- name: DeleteChat :execrows
DELETE FROM chats
WHERE id = sqlc.arg('id');


-- name: GetLatestChatPreview :many
SELECT 
    cr.id   AS chatroom_id,
    cr.type AS type,
    COALESCE(lc.created_at, cr.created_at) AS last_message_at,

    -- Nama: Ambil dari nama project atau nama user lawan bicara
    (CASE 
        WHEN cr.type = 'project' THEN p.name 
        ELSE u.full_name 
    END)::text AS name,

    -- Avatar: Ambil profile picture user lawan (untuk personal)
   COALESCE(
        (CASE 
            WHEN cr.type = 'personal' THEN u.profile_picture 
            ELSE NULL 
        END),
        ''
    )::text AS avatar,

    (
        COALESCE(lc.text, '') || 
        CASE 
            WHEN has_media.exists THEN ' [media]' 
            ELSE '' 
        END
    )::text AS last_message

FROM chatrooms cr

-- SEMUA JOIN HARUS DI SINI (SEBELUM WHERE)

-- 1. Ambil detail Project jika tipenya project
LEFT JOIN projects p 
    ON cr.type = 'project' AND cr.project_id = p.id

-- 2. Ambil detail User Lawan jika tipenya personal
LEFT JOIN LATERAL (
    SELECT 
        u.full_name,
        u.profile_picture
    FROM chatroom_participants cp_other
    JOIN users u ON u.id = cp_other.user_id
    WHERE cp_other.chatroom_id = cr.id
      AND cp_other.user_id != sqlc.arg('current_user_id')
      AND cp_other.left_at IS NULL
    LIMIT 1
) u ON cr.type = 'personal'

-- 3. Ambil chat terakhir
LEFT JOIN LATERAL (
    SELECT id, text, created_at
    FROM chats c
    WHERE c.room_id = cr.id
    ORDER BY c.created_at DESC
    LIMIT 1
) lc ON TRUE

-- 4. Cek apakah ada media di chat terakhir tersebut
LEFT JOIN LATERAL (
    SELECT EXISTS (
        SELECT 1 FROM chat_medias cm 
        WHERE cm.chat_id = lc.id
    ) AS exists
) has_media ON lc.id IS NOT NULL

-- KLAUSA WHERE SETELAH SEMUA JOIN SELESAI
WHERE (
    (cr.type = 'project' AND EXISTS (
        SELECT 1 FROM project_members pm 
        WHERE pm.project_id = cr.project_id 
          AND pm.user_id = sqlc.arg('current_user_id') 
          AND pm.left_at IS NULL
    ))
    OR 
    (cr.type = 'personal' AND EXISTS (
        SELECT 1 FROM chatroom_participants cp 
        WHERE cp.chatroom_id = cr.id 
          AND cp.user_id = sqlc.arg('current_user_id') 
          AND cp.left_at IS NULL
    ))
)

ORDER BY lc.created_at DESC NULLS LAST;


-- name: GetChatID :one
SELECT
    c.id,
    c.room_id,
    c.text,
    c.created_at,
    jsonb_build_object(
        'id',              u.id,
        'full_name',       u.full_name,
        'email',           u.email,
        'gender',          u.gender,
        'profile_picture', u.profile_picture,
        'global_role',     u.global_role,
        'created_at',      u.created_at,
        'updated_at',      u.updated_at
    ) AS sender,
    COALESCE(
        jsonb_agg(
            jsonb_build_object(
                'id',         cm.id,
                'file_name',  cm.file_name,
                'media_type', cm.media_type,
                'size',       cm.size,
                'created_at', cm.created_at
            )
        ) FILTER (WHERE cm.id IS NOT NULL),
        '[]'::jsonb
    )::jsonb AS medias
FROM chats c
LEFT JOIN users u
    ON u.id = c.sender_id
   AND u.deleted_at IS NULL
LEFT JOIN chat_medias cm
    ON cm.chat_id = c.id
WHERE c.id = sqlc.arg('id')
GROUP BY
    c.id,
    c.room_id,
    c.text,
    c.created_at,
    u.id,
    u.full_name,
    u.email,
    u.gender,
    u.profile_picture,
    u.global_role,
    u.created_at,
    u.updated_at
ORDER BY c.created_at ASC;