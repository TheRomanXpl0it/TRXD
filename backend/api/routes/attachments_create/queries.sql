-- name: CreateAttachment :exec
-- Creates an attachment for a challenge
INSERT INTO attachments (chall_id, name, hash) VALUES ($1, $2, $3);
