-- name: GetUsersPreview :many
-- Retrieve all users
SELECT id, name, email, role, score, country
  FROM users
  ORDER BY id ASC;
