-- name: AddConfig :exec
-- Insert a new configuration setting
INSERT INTO configs (key, type, value) VALUES ($1, $2, $3);

-- name: CreateCategory :exec
-- Insert a new category
INSERT INTO categories (name, icon) VALUES ($1, $2);

-- name: CreateChallenge :exec
-- Insert a new challenge
INSERT INTO challenges (name, category, description, type, max_points, score_type)
	VALUES ($1, $2, $3, $4, $5, $6);

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

-- name: GetChallengeByID :one
-- Retrieve a challenge by its ID
SELECT * FROM challenges WHERE id = $1;

-- name: Submit :one
-- Insert a new submission
INSERT INTO submissions (user_id, chall_id, status, flag) VALUES ($1, $2, $3, $4) RETURNING status;

-- name: CheckFlags :one
SELECT BOOL_OR(($1 = flag) OR (regex AND $1 ~ flag)) FROM flags WHERE chall_id = $2;
