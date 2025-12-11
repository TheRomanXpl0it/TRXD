-- name: GetChallengeByID :one
-- Retrieve a challenge by its ID
SELECT * FROM challenges WHERE id = $1;

-- name: GetDockerConfigsByID :one
-- Retrieve Docker configurations by challenge ID
SELECT
  image,
  compose,
  hash_domain,
  envs,
  COALESCE(NULLIF(lifetime, 0), (SELECT value::INTEGER FROM configs WHERE key='instance-lifetime')) AS lifetime,
  COALESCE(NULLIF(max_memory, 0), (SELECT value::INTEGER FROM configs WHERE key='instance-max-memory')) AS max_memory,
  COALESCE(NULLIF(max_cpu, ''), (SELECT value FROM configs WHERE key='instance-max-cpu')) AS max_cpu
FROM docker_configs
WHERE chall_id = $1;

-- name: GetHiddenAndAttachments :one
-- Checks if a challenge is hidden
SELECT
  c.hidden, 
  (ARRAY_AGG(a.hash || '/' || a.name) FILTER (WHERE a.name IS NOT NULL))::TEXT[] AS attachments
FROM challenges c
LEFT JOIN attachments a
  ON a.chall_id = c.id
WHERE c.id = $1
GROUP BY c.hidden;
