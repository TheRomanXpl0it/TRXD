-- name: UpdateInstance :exec
-- Update an instance expiration time
UPDATE instances
  SET expires_at = $3
  WHERE team_id = $1 AND chall_id = $2;
