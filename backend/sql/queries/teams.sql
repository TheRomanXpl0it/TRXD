-- name: GetTeamByID :one
-- Retrieve a team by its ID
SELECT * FROM teams WHERE id = $1;

-- name: GetTeamFromUser :one
-- Retrieve the team associated with a user
SELECT t.* FROM teams t
  JOIN users u ON u.team_id = t.id
  WHERE u.id = $1;

-- name: GetTeamByName :one
-- Retrieve a team by its name
SELECT * FROM teams WHERE name = $1;

-- name: GetBadgesFromTeam :many
-- Retrieve all badges associated with a team
SELECT badges.name, badges.description FROM badges
  JOIN teams ON teams.id = badges.team_id
  WHERE teams.id = $1;
