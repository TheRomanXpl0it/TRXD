-- name: GetNextInstanceToDelete :one
-- Retrieves the next instance to delete
SELECT team_id, chall_id, expires_at
  FROM instances
  WHERE expires_at < NOW() + (
    (SELECT value
      FROM configs
      WHERE key='reclaim-instance-interval'
    ) || ' seconds')::INTERVAL
  ORDER BY expires_at ASC
  LIMIT 1;
