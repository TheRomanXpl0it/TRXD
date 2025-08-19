-- name: GetTeamMembers :many
-- Retrieve all members of a team by team ID
SELECT id, name, role, score FROM users WHERE team_id = $1 ORDER BY id;

-- name: GetTeamSolves :many
-- Retrieve all challenges solved by a team's members
SELECT challenges.id, challenges.name, challenges.category, submissions.timestamp
  FROM submissions
  JOIN users ON users.id = submissions.user_id
  JOIN teams ON users.team_id = teams.id
  JOIN challenges ON challenges.id = submissions.chall_id
  WHERE users.role = 'Player'
    AND teams.id = $1
    AND submissions.status = 'Correct'
  ORDER BY submissions.timestamp DESC;
