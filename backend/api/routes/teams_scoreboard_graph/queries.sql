-- name: GetTeamsScoreboardGraph :many
-- Get the top N teams along with their correct submissions and challenge points for scoreboard graphing
SELECT t.id AS team_id, c.id AS chall_id, c.points, s.first_blood, s."timestamp" FROM (
	SELECT * FROM teams t
        ORDER BY t.score DESC
        LIMIT CAST((SELECT value FROM configs WHERE key='scoreboard-top') AS INT)) AS t
  JOIN users u ON u.team_id = t.id
  JOIN submissions s ON s.user_id = u.id
  JOIN challenges c ON c.id = s.chall_id
  WHERE s.status = 'Correct'
	AND u.role = 'Player'
  ORDER BY s."timestamp" ASC;
