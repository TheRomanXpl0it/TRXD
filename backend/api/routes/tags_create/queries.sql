-- name: CreateTag :exec
-- Creates a named tag for a challenge
INSERT INTO tags (chall_id, name) VALUES ($1, $2);
