-- name: GetConfigs :many
-- Fetches all configuration settings
SELECT * FROM configs ORDER BY key;
