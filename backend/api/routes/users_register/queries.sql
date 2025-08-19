-- name: RegisterUser :one
-- Insert a new user and return the created user
INSERT INTO users (name, email, password_hash, role) VALUES ($1, $2, $3, $4) RETURNING *;
