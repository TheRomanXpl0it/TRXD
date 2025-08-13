-- name: UpdateTeam :exec
-- Update team details
UPDATE teams
SET
  nationality = COALESCE(sqlc.narg('nationality'), nationality),
  image = COALESCE(sqlc.narg('image'), image),
  bio = COALESCE(sqlc.narg('bio'), bio)
WHERE id = $1;
