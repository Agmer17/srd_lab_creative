-- name: CreateChatMedia :many
INSERT INTO chat_medias (chat_id, file_name, media_type)
VALUES (
    unnest(@chat_id::uuid[]),
    unnest(@filename::text[]),
    unnest(@media_type::chat_media_type[])
)
RETURNING *;

-- name: GetChatMediasByRoomID :many
SELECT
    cm.id,
    cm.chat_id,
    cm.file_name,
    cm.media_type,
    cm.size,
    cm.created_at
FROM chat_medias cm
JOIN chats c
    ON c.id = cm.chat_id
WHERE c.room_id = sqlc.arg('room_id')
ORDER BY cm.created_at DESC;

-- name: GetChatMediasByChatID :many
SELECT
    cm.id,
    cm.chat_id,
    cm.file_name,
    cm.media_type,
    cm.size,
    cm.created_at
FROM chat_medias cm
WHERE cm.chat_id = sqlc.arg('chat_id');