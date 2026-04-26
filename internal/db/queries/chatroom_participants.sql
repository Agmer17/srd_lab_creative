-- name: AddPersonalChatroomParticipant :one
INSERT INTO chatroom_participants (chatroom_id, user_id)
VALUES (sqlc.arg('chatroom_id'), sqlc.arg('user_id'))
RETURNING *;
