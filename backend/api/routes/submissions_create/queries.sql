-- name: CheckFlags :one
-- Check if a flag matches any flags for a challenge
SELECT COALESCE(BOOL_OR(($1 = flag) OR (regex AND $1 ~ flag)), false)::BOOLEAN FROM flags WHERE chall_id = $2;

-- name: Submit :one
-- Insert a new submission
WITH challenge AS (
    SELECT challenges.id FROM challenges
    WHERE challenges.id = $2 FOR UPDATE
  )
INSERT INTO submissions (user_id, chall_id, status, flag)
  VALUES ($1, (SELECT id FROM challenge), sqlc.arg(status), $3)
  RETURNING status, first_blood;
