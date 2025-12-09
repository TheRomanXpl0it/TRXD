-- name: UpdateTeam :exec
-- Update team details
UPDATE teams
SET
  name = COALESCE(sqlc.narg('name'), name),
  country = COALESCE(sqlc.narg('country'), country),
  image = COALESCE(sqlc.narg('image'), image)
WHERE id = $1;
