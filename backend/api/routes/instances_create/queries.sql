-- name: GetInstance :one
-- Gets an instance by ID
SELECT * FROM instances WHERE chall_id = $1 AND team_id = $2;

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



-- WITH conf_secret AS (
--     SELECT value AS secret
--     FROM configs
--     WHERE key='secret'
--     FOR UPDATE
--   ),
--   available_port AS (
--     SELECT get_random_available_port() AS port
--   ),
--   conf_domain AS (
--     SELECT value AS domain
--     FROM configs
--     WHERE key='domain'
--   )
-- INSERT INTO instances (team_id, chall_id, expires_at, host, port)
--   VALUES (
--     $1,
--     $2,
--     $3,
--     CASE
--       WHEN sqlc.arg(hash_domain)::BOOLEAN THEN (
--         SELECT generate_team_hash(
--           (SELECT secret FROM conf_secret),
--           $1,
--           (SELECT port FROM available_port)
--         ) || '.' || (SELECT domain FROM conf_domain)
--       )
--       ELSE (SELECT domain FROM conf_domain)
--     END,
--     (SELECT port FROM available_port)
--   )
-- RETURNING host, port;
