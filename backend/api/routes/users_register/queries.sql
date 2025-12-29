-- name: RegisterUser :one
-- Insert a new user and return the created user
INSERT INTO users (name, email, password_hash, password_salt, role) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UserExistsByEmail :one
-- Check if a user exists by email
SELECT EXISTS(SELECT 1 FROM users WHERE email = $1) AS exists;
