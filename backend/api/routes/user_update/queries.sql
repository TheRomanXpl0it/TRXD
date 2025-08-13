-- name: UpdateUser :exec
-- Update user details
UPDATE users
SET
  name = COALESCE(sqlc.narg('name'), name),
  nationality = COALESCE(sqlc.narg('nationality'), nationality),
  image = COALESCE(sqlc.narg('image'), image)
WHERE id = $1;
