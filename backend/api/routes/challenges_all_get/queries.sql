-- name: GetChallengesPreview :many
-- Retrieve all challenges
SELECT id, name, category, difficulty, type, hidden, points, solves,
    EXISTS(
      SELECT 1
        FROM submissions
        JOIN users ON users.id = submissions.user_id
        JOIN teams ON users.team_id = teams.id
          AND teams.id = (SELECT team_id FROM users WHERE users.id = $1)
        WHERE users.role = 'Player'
          AND submissions.status = 'Correct'
        AND submissions.chall_id = challenges.id
    ) AS solved,
    EXISTS(
      SELECT 1
        FROM submissions
        JOIN users ON users.id = submissions.user_id
        JOIN teams ON users.team_id = teams.id
          AND teams.id = (SELECT team_id FROM users WHERE users.id = $1)
        WHERE users.role = 'Player'
          AND submissions.status = 'Correct'
          AND submissions.chall_id = challenges.id
        ORDER BY submissions.timestamp
        LIMIT 1
    ) AS first_blood
  FROM challenges;

-- name: GetInstanceExpire :one
-- Retrieve the expiration time of a specific instance
SELECT expires_at
  FROM instances
  JOIN teams ON teams.id = instances.team_id
  JOIN users ON users.team_id = teams.id
  WHERE users.id = sqlc.arg(user_id) AND chall_id = $1;
