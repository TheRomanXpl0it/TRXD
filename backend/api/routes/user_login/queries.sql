-- name: GetUserByEmail :one
-- Retrieve a user by their email address
SELECT * FROM users WHERE email = $1;
