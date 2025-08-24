-- name: DeleteInstance :exec
-- Delete an instance
DELETE FROM instances
  WHERE team_id = $1 AND chall_id = $2;
