-- name: GetChallengesPreview :many
-- Retrieve all challenges
SELECT id, name, category, difficulty, type, hidden, points, solves, EXISTS(
  SELECT 1
    FROM submissions
    JOIN users ON users.id = submissions.user_id
    JOIN teams ON users.team_id = teams.id
      AND teams.id = (SELECT team_id FROM users WHERE users.id = $1)
    WHERE users.role = 'Player'
      AND submissions.status = 'Correct'
      AND submissions.chall_id = challenges.id) AS solved
  FROM challenges;
