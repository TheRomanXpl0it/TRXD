-- name: CreateConfig :exec
-- Insert a new configuration setting
INSERT INTO configs (key, type, value) VALUES ($1, $2, $3);

-- name: UpdateConfig :exec
-- Update an existing configuration setting
UPDATE configs SET value = $2 WHERE key = $1;

-- name: GetConfig :one
-- Retrieve a configuration setting by key
SELECT * FROM configs WHERE key = $1;

-- name: CreateCategory :exec
-- Insert a new category
INSERT INTO categories (name, icon) VALUES ($1, $2);

-- name: CreateChallenge :one
-- Insert a new challenge
INSERT INTO challenges (name, category, description, type, max_points, score_type)
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;

-- name: CreateFlag :exec
-- Insert a new flag for a challenge
INSERT INTO flags (flag, chall_id, regex) VALUES ($1, $2, $3);

-- name: RegisterUser :one
-- Insert a new user and return the created user
INSERT INTO users (name, email, password_hash, role) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: RegisterTeam :exec
-- Insert a new team and add the founder user to the team
WITH locked_user AS (
    SELECT id FROM users
    WHERE id = $1 AND team_id IS NULL
    FOR UPDATE
  ),
  new_team AS (
    INSERT INTO teams (name, password_hash)
    SELECT $2, $3
    FROM locked_user
    RETURNING *
  )
UPDATE users
  SET team_id = new_team.id
  FROM new_team
  WHERE users.id = $1;

-- name: GetTeamFromUser :one
-- Retrieve the team associated with a user
SELECT t.* FROM teams t
  JOIN users u ON u.team_id = t.id
  WHERE u.id = $1;

-- name: AddTeamMember :exec
-- Assign a user to a team
UPDATE users SET team_id = $1 WHERE id = $2 AND team_id IS NULL;

-- name: GetUserByID :one
-- Retrieve a user by their ID
SELECT * FROM users WHERE id = $1;

-- name: GetUserByName :one
-- Retrieve a user by their name
SELECT * FROM users WHERE name = $1;

-- name: GetUserByEmail :one
-- Retrieve a user by their email address
SELECT * FROM users WHERE email = $1;

-- name: GetTeamByName :one
-- Retrieve a team by its name
SELECT * FROM teams WHERE name = $1;

-- name: GetTeamByID :one
-- Retrieve a team by its ID
SELECT * FROM teams WHERE id = $1;

-- name: GetTeamMembers :many
-- Retrieve all members of a team by team ID
SELECT id, name, role, score FROM users WHERE team_id = $1;

-- name: GetTeams :many
-- Retrieve all teams
SELECT id FROM teams;

-- name: GetChallengeByID :one
-- Retrieve a challenge by its ID
SELECT * FROM challenges WHERE id = $1;

-- name: IsChallengeSolved :one
-- Check if a challenge is solved by a user's team
SELECT EXISTS(
  SELECT 1
    FROM submissions
    JOIN users ON users.id = submissions.user_id
    JOIN teams ON users.team_id = teams.id
      AND teams.id = (SELECT team_id FROM users WHERE users.id = $2)
    WHERE users.role = 'Player'
      AND submissions.status = 'Correct'
      AND submissions.chall_id = $1
);

-- name: GetUsers :many
-- Retrieve all users
SELECT id FROM users;

-- name: GetUserSolves :many
-- Retrieve all challenges solved by a user
SELECT c.id, c.name, c.category, s.timestamp FROM challenges c
  JOIN submissions s ON s.chall_id = c.id
    WHERE s.user_id = $1
      AND s.status = 'Correct';

-- name: GetChallenges :many
-- Retrieve all challenges
SELECT id FROM challenges;

-- name: GetFlagsByChallenge :many
-- Retrieve all flags associated with a challenge
SELECT flag, regex FROM flags WHERE chall_id = $1;

-- name: GetTagsByChallenge :many
-- Retrieve all tags associated with a challenge
SELECT name FROM tags WHERE chall_id = $1;

-- name: Submit :one
-- Insert a new submission
INSERT INTO submissions (user_id, chall_id, status, flag) VALUES ($1, $2, $3, $4) RETURNING status;

-- name: CheckFlags :one
-- Check if a flag matches any flags for a challenge
SELECT BOOL_OR(($1 = flag) OR (regex AND $1 ~ flag)) FROM flags WHERE chall_id = $2;

-- name: UpdateUser :exec
-- Update user details
UPDATE users
SET
  name = COALESCE(sqlc.narg('name'), name),
  nationality = COALESCE(sqlc.narg('nationality'), nationality),
  image = COALESCE(sqlc.narg('image'), image)
WHERE id = $1;

-- name: UpdateTeam :exec
-- Update team details
UPDATE teams
SET
  nationality = COALESCE(sqlc.narg('nationality'), nationality),
  image = COALESCE(sqlc.narg('image'), image),
  bio = COALESCE(sqlc.narg('bio'), bio)
WHERE id = $1;

-- name: DeleteCategory :exec
-- Delete a category and all associated challenges
DELETE FROM categories WHERE name = $1;

-- name: DeleteChallenge :exec
-- Delete a challenge and all associated flags
DELETE FROM challenges WHERE id = $1;

-- name: DeleteFlag :exec
-- Delete a flag from a challenge
DELETE FROM flags WHERE chall_id = $1 AND flag = $2;

-- name: ResetUserPassword :exec
-- Reset a user's password to a new password
UPDATE users SET password_hash = $1 WHERE id = $2;

-- name: ResetTeamPassword :exec
-- Reset a team's password to a new password
UPDATE teams SET password_hash = $1 WHERE id = $2;
