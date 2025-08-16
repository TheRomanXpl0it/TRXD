-- name: GetUserByID :one
-- Retrieve a user by their ID
SELECT * FROM users WHERE id = $1;

-- name: GetUserByName :one
-- Retrieve a user by their name
SELECT * FROM users WHERE name = $1;
