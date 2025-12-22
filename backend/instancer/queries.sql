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
WITH info AS (
    SELECT generate_instance_remote(
      sqlc.arg(chall_id),
      sqlc.arg(hash_domain)::BOOLEAN
    ) AS remote
  )
INSERT INTO instances (team_id, chall_id, expires_at, host, port)
  VALUES (sqlc.arg(team_id), sqlc.arg(chall_id), sqlc.arg(expires_at),
    (SELECT (remote).host FROM info), (SELECT (remote).port FROM info))
RETURNING host, port;

-- name: UpdateInstanceDockerID :exec
-- Adds the container ID to the instance
UPDATE instances
  SET docker_id = $3
  WHERE team_id = $1 AND chall_id = $2;

-- name: UpdateInstanceExpire :exec
-- Update an instance expiration time
UPDATE instances
  SET expires_at = $3
  WHERE team_id = $1 AND chall_id = $2;

-- name: DeleteInstance :exec
-- Delete an instance
DELETE FROM instances
  WHERE team_id = $1 AND chall_id = $2;

