-- name: DeleteFlag :exec
-- Delete a flag from a challenge
DELETE FROM flags WHERE chall_id = $1 AND flag = $2;
