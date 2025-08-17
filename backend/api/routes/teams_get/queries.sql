-- name: GetTeamsPreview :many
-- Retrieve all teams
SELECT id, name, score, country, image
  FROM teams
  ORDER BY id;
