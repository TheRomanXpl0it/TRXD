-- name: GetUserIDByEmail :one
SELECT id FROM users WHERE email = $1;
