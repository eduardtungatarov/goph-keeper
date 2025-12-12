-- name: SaveData :one
INSERT INTO data (user_id, title, type, data)
VALUES (sqlc.arg(user_id), sqlc.arg(title), sqlc.arg(type), sqlc.arg(data))
RETURNING *;

-- name: FindDataByUserIDAndID :one
SELECT * FROM data
WHERE user_id = sqlc.arg(user_id) and id = sqlc.arg(id)
LIMIT 1;

-- name: DeleteDataByUserIDAndID :execrows
DELETE FROM data
WHERE user_id = sqlc.arg(user_id) AND id = sqlc.arg(id);