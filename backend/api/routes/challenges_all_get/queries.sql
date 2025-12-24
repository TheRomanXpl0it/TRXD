-- name: GetAllChallengesInfo :many
-- Retrieve all challenges along with first blood status and instance info for a user
WITH tid AS (SELECT team_id FROM users WHERE users.id = $1)
SELECT
    c.*,
    (s.first_blood IS NOT NULL)::BOOLEAN AS solved,
    COALESCE(s.first_blood, FALSE) AS first_blood,
    (ARRAY_AGG('/' || a.hash || '/' || a.name ORDER BY a.name)
      FILTER (WHERE a.name IS NOT NULL))::TEXT[] AS attachments,
    i.expires_at,
    i.host AS instance_host,
    i.port AS instance_port,
    i.docker_id
  FROM challenges c
  LEFT JOIN attachments a
    ON a.chall_id = c.id
  LEFT JOIN (
      SELECT submissions.chall_id, submissions.first_blood
        FROM submissions
        JOIN users ON users.id = submissions.user_id
        WHERE users.team_id = (SELECT team_id FROM tid)
          AND users.role = 'Player'
          AND submissions.status = 'Correct') s
    ON s.chall_id = c.id
  LEFT JOIN instances i
    ON i.chall_id = c.id
      AND i.team_id = (SELECT team_id FROM tid)
  GROUP BY c.id, s.first_blood, i.expires_at, i.host, i.port, i.docker_id 
  ORDER BY c.points ASC, c.id ASC;
