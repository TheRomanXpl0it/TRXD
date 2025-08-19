-- name: UpdateUser :exec
-- Update user details
UPDATE users
SET
  name = COALESCE(sqlc.narg('name'), name),
  country = COALESCE(sqlc.narg('country'), country),
  image = COALESCE(sqlc.narg('image'), image)
WHERE id = $1;
