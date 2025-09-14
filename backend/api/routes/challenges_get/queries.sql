-- name: GetFlagsByChallenge :many
-- Retrieve all flags associated with a challenge
SELECT flag, regex FROM flags WHERE chall_id = $1;

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

-- name: GetChallengeSolves :many
-- Retrieve all teams that solved a challenge
SELECT teams.id, teams.name, submissions.timestamp
  FROM submissions
  JOIN users ON users.id = submissions.user_id
  JOIN teams ON users.team_id = teams.id
  WHERE users.role = 'Player'
    AND submissions.chall_id = $1
    AND submissions.status = 'Correct'
  ORDER BY submissions.timestamp ASC;

-- name: GetChallDockerConfig :one
SELECT * FROM docker_configs WHERE chall_id = $1;

-- name: GetInstanceInfo :one
-- Retrieve the instance associated with a challenge and team
SELECT expires_at, host, port FROM instances WHERE team_id = $1 AND chall_id = $2;
