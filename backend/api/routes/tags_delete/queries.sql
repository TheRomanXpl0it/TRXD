-- name: DeleteTag :exec
-- Deletes a challenge tag by name
DELETE FROM tags WHERE chall_id = $1 AND name = $2;
