-- name: UpdateTag :exec
-- Updates the name of a challenge tag
UPDATE tags SET name = sqlc.arg(new_name) WHERE chall_id = $1 AND name = sqlc.arg(old_name);
