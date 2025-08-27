-- name: GetInstance :one
-- Gets an instance by ID
SELECT * FROM instances WHERE chall_id = $1 AND team_id = $2;

-- name: GetNextInstanceToDelete :one
-- Retrieves the next instance to delete
SELECT team_id, chall_id, expires_at, docker_id
  FROM instances
  WHERE expires_at < NOW() + (
    (SELECT value
      FROM configs
      WHERE key='reclaim-instance-interval'
    ) || ' seconds')::INTERVAL
  ORDER BY expires_at ASC
  LIMIT 1;

-- name: CreateInstance :one
-- Creates a new instance for a team
WITH conf_secret AS (
    SELECT value AS secret
    FROM configs
    WHERE key='secret'
    FOR UPDATE
  ),
  info AS (
    SELECT generate_instance_remote(
      (SELECT secret FROM conf_secret),
      $1,
      $2,
      sqlc.arg(hash_domain)::BOOLEAN
    ) AS tuple
  )
INSERT INTO instances (team_id, chall_id, expires_at, host, port)
  VALUES ($1, $2, $3, (SELECT (tuple).host FROM info), (SELECT (tuple).port FROM info))
RETURNING host, port;

-- name: UpdateInstanceDockerID :exec
-- Adds the container ID to the instance
UPDATE instances
  SET docker_id = $3
  WHERE team_id = $1 AND chall_id = $2;

-- name: DeleteInstance :exec
-- Delete an instance
DELETE FROM instances
  WHERE team_id = $1 AND chall_id = $2;

