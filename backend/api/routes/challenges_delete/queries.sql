-- name: DeleteChallenge :exec
-- Delete a challenge and all associated flags
DELETE FROM challenges WHERE id = $1;
