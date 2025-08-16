-- name: GetChallengeByID :one
-- Retrieve a challenge by its ID
SELECT * FROM challenges WHERE id = $1;

-- name: GetTagsByChallenge :many
-- Retrieve all tags associated with a challenge
SELECT name FROM tags WHERE chall_id = $1;

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
