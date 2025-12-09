-- name: UpdateUser :exec
-- Update user details
UPDATE users
SET
  name = COALESCE(sqlc.narg('name'), name),
  country = COALESCE(sqlc.narg('country'), country)
WHERE id = $1;
