-- name: CreateConfig :exec
-- Insert a new configuration setting
INSERT INTO configs (key, type, value) VALUES ($1, $2, $3);

-- name: UpdateConfig :exec
-- Update an existing configuration setting
UPDATE configs SET value = $2 WHERE key = $1;

-- name: GetConfig :one
-- Retrieve a configuration setting by key
SELECT * FROM configs WHERE key = $1;
