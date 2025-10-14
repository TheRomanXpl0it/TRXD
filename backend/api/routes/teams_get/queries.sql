-- name: GetTeamMembers :many
-- Retrieve all members of a team by team ID
SELECT id, name, role, score FROM users WHERE team_id = $1 ORDER BY id;

-- name: GetTeamSolves :many
-- Retrieve all challenges solved by a team's members
SELECT c.id, c.name, c.category, s.first_blood, s.timestamp, s.user_id
  FROM submissions s
  JOIN users u ON u.id = s.user_id
  JOIN teams t ON u.team_id = t.id
  JOIN challenges c ON c.id = s.chall_id
  WHERE u.role = 'Player'
    AND t.id = $1
    AND s.status = 'Correct'
  ORDER BY s.timestamp DESC;
