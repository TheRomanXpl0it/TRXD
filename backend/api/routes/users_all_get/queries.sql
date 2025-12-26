-- name: GetUsers :many
-- Retrieve all users
SELECT id, name, email, role, score, country
  FROM users
  WHERE sqlc.arg('is_admin')::BOOLEAN
    OR role = 'Player'
  ORDER BY id ASC
  OFFSET sqlc.arg('offset')
  LIMIT sqlc.narg('limit');
