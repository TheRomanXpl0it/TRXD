-- name: GetInstances :many
-- Get all instances
SELECT
  i.team_id,
  t.name AS team_name,
  i.chall_id,
  c.name AS chall_name,
  i.expires_at,
  c.conn_type,
  i.host,
  COALESCE(i.port, 0) AS port,
  COALESCE(i.docker_id, '') AS docker_id
FROM instances i
JOIN teams t ON i.team_id = t.id
JOIN challenges c ON i.chall_id = c.id
ORDER BY i.expires_at ASC;
