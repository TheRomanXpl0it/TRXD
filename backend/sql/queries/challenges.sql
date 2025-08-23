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
  COALESCE(lifetime, (SELECT value::INTEGER FROM configs WHERE key='instance-lifetime')) AS lifetime,
  COALESCE(max_memory, (SELECT value::INTEGER FROM configs WHERE key='instance-max-memory')) AS max_memory,
  COALESCE(max_cpu, (SELECT value FROM configs WHERE key='instance-max-cpu')) AS max_cpu
FROM docker_configs
WHERE chall_id = $1;

-- name: GetTagsByChallenge :many
-- Retrieve all tags associated with a challenge
SELECT name FROM tags WHERE chall_id = $1;
