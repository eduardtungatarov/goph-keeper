-- name: SaveData :one
INSERT INTO data (user_id, title, type, data)
VALUES (sqlc.arg(user_id), sqlc.arg(title), sqlc.arg(type), sqlc.arg(data))
RETURNING *;