-- name: GetChallenges :many
-- Retrieve all challenges
SELECT * FROM challenges;

-- name: IsChallengeFirstBlood :one
-- Check if a challenge is solved (and is a first blood) by a user's team
SELECT first_blood
  FROM submissions
  JOIN users ON users.id = submissions.user_id
  JOIN teams ON users.team_id = teams.id
    AND teams.id = (SELECT team_id FROM users WHERE users.id = $2)
  WHERE users.role = 'Player'
    AND submissions.status = 'Correct'
    AND submissions.chall_id = $1;

-- name: GetInstanceInfo :one
-- Retrieve the instance associated with a challenge and team
SELECT expires_at, host, port, docker_id FROM instances WHERE team_id = $1 AND chall_id = $2;
