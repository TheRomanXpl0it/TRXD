-- name: UpdateTeam :exec
-- Update team details
UPDATE teams
SET
  country = COALESCE(sqlc.narg('country'), country),
  image = COALESCE(sqlc.narg('image'), image),
  bio = COALESCE(sqlc.narg('bio'), bio)
WHERE id = $1;
