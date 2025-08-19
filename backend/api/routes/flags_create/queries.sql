-- name: CreateFlag :exec
-- Insert a new flag for a challenge
INSERT INTO flags (flag, chall_id, regex) VALUES ($1, $2, $3);
