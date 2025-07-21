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

-- name: RegisterTeam :one
-- Insert a new team and return the created team
INSERT INTO teams (name, password_hash) VALUES ($1, $2) RETURNING *;

-- name: AddTeamMember :exec
-- Assign a user to a team
UPDATE users SET team_id = $1 WHERE id = $2;

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

-- name: GetChallengesByCategory :many
-- Retrieve all challenges in a specific category
SELECT * FROM challenges WHERE category = $1;

-- name: Submit :exec
-- Insert a new submission
INSERT INTO submissions (user_id, chall_id, status, flag) VALUES ($1, $2, $3, $4);
