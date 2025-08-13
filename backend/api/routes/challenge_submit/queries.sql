-- name: CheckFlags :one
-- Check if a flag matches any flags for a challenge
SELECT BOOL_OR(($1 = flag) OR (regex AND $1 ~ flag)) FROM flags WHERE chall_id = $2;

-- name: Submit :one
-- Insert a new submission
INSERT INTO submissions (user_id, chall_id, status, flag) VALUES ($1, $2, $3, $4) RETURNING status;
