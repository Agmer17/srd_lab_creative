-- name: AddPersonalChatroomParticipant :many
INSERT INTO chatroom_participants (chatroom_id, user_id)
VALUES (unnest(@chatroom_id::uuid[]), unnest(@user_id::uuid[]))
RETURNING *;
