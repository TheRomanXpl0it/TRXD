-- name: CreateConfig :exec
-- Insert a new configuration setting
INSERT INTO configs (key, type, value) VALUES ($1, $2, $3);

-- name: UpdateConfig :exec
-- Update an existing configuration setting
UPDATE configs SET value = $2 WHERE key = $1;

-- name: GetConfig :one
-- Retrieve a configuration setting by key
SELECT * FROM configs WHERE key = $1;

-- name: GetTeamFromUser :one
-- Retrieve the team associated with a user
SELECT t.* FROM teams t
  JOIN users u ON u.team_id = t.id
  WHERE u.id = $1;

-- name: GetUserByID :one
-- Retrieve a user by their ID
SELECT * FROM users WHERE id = $1;

-- name: GetUserByName :one
-- Retrieve a user by their name
SELECT * FROM users WHERE name = $1;

-- name: GetTeamByName :one
-- Retrieve a team by its name
SELECT * FROM teams WHERE name = $1;

-- name: GetChallengeByID :one
-- Retrieve a challenge by its ID
SELECT * FROM challenges WHERE id = $1;
