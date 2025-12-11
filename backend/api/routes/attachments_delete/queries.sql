-- name: DeleteAttachment :exec
-- Deletes a challenge attachment by name
DELETE FROM attachments WHERE chall_id = $1 AND name = $2;

-- name: GetAttachmentHash :one
-- Retrieves the hash of a challenge attachment by name
SELECT hash FROM attachments WHERE chall_id = $1 AND name = $2;
