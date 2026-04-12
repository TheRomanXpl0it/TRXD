-- name: GetTeamIDByEmail :one
SELECT team_id FROM users WHERE email = $1;

-- name: GetTeamIDByName :one
SELECT id FROM teams WHERE name = $1;
