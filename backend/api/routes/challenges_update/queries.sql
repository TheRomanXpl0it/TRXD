-- name: UpdateChallenge :exec
-- Updates the challenge with the given ID
UPDATE challenges
SET
  name = COALESCE(sqlc.narg('name'), name),
  category = COALESCE(sqlc.narg('category'), category),
  description = COALESCE(sqlc.narg('description'), description),
  difficulty = COALESCE(sqlc.narg('difficulty'), difficulty),
  authors = COALESCE(sqlc.narg('authors'), authors),
  type = COALESCE(sqlc.narg('type'), type),
  hidden = COALESCE(sqlc.narg('hidden'), hidden),
  max_points = COALESCE(sqlc.narg('max_points'), max_points),
  score_type = COALESCE(sqlc.narg('score_type'), score_type),
  host = COALESCE(sqlc.narg('host'), host),
  port = COALESCE(sqlc.narg('port'), port),
  attachments = COALESCE(sqlc.narg('attachments'), attachments)
WHERE id = sqlc.arg('chall_id');

-- name: UpdateDockerConfigs :exec
-- Updates the Docker configurations for the challenge with the given ID
UPDATE docker_configs
SET
  image = COALESCE(sqlc.narg('image'), image),
  compose = COALESCE(sqlc.narg('compose'), compose),
  hash_domain = COALESCE(sqlc.narg('hash_domain'), hash_domain),
  lifetime = COALESCE(sqlc.narg('lifetime'), lifetime),
  envs = COALESCE(sqlc.narg('envs'), envs),
  max_memory = COALESCE(sqlc.narg('max_memory'), max_memory),
  max_cpu = COALESCE(sqlc.narg('max_cpu'), max_cpu)
WHERE chall_id = sqlc.arg('chall_id');
