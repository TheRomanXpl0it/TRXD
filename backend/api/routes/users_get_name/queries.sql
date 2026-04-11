-- name: GetUserIDByName :one
SELECT id FROM users WHERE name = $1;
