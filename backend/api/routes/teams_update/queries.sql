-- name: UpdateTeam :exec
-- Update team details
UPDATE teams
SET
  name = COALESCE(sqlc.narg('name'), name),
  country = COALESCE(sqlc.narg('country'), country)
WHERE id = $1;
