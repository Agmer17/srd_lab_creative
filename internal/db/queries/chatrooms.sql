-- name: CreateProjectChatroom :one
INSERT INTO chatrooms (type, project_id)
VALUES ('project', sqlc.arg('project_id'))
RETURNING *;

-- name: CreatePersonalChatroom :one
INSERT INTO chatrooms (type, participant_key)
VALUES ('personal', sqlc.arg('participant_key'))
RETURNING *;

-- name: GetChatroomByID :one
SELECT * FROM chatrooms
WHERE id = sqlc.arg('id');

-- name: GetChatroomByProjectID :one
SELECT * FROM chatrooms
WHERE project_id = sqlc.arg('project_id');

-- name: GetChatroomByParticipantKey :one
SELECT * FROM chatrooms
WHERE participant_key = sqlc.arg('participant_key');