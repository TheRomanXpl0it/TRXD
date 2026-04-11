-- name: GetUserIDByEmail :one
SELECT id FROM users WHERE email = $1;

-- name: GetUserIDByName :one
SELECT id FROM users WHERE name = $1;
