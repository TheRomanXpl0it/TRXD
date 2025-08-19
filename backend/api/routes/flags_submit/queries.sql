-- name: CheckFlags :one
-- Check if a flag matches any flags for a challenge
SELECT BOOL_OR(($1 = flag) OR (regex AND $1 ~ flag)) FROM flags WHERE chall_id = $2;

-- name: Submit :one
-- Insert a new submission
WITH challenge AS (
    SELECT challenges.id FROM challenges
    WHERE challenges.id = $2 FOR UPDATE
  ),
  inserted AS (
    INSERT INTO submissions (user_id, chall_id, status, flag)
    VALUES ($1, (SELECT id FROM challenge), sqlc.arg(status), $3)
    RETURNING status
  ),
  blood_check AS (
    SELECT COUNT(*)=0 AS first_blood FROM submissions
    WHERE chall_id = (SELECT id FROM challenge) AND status = 'Correct'
  )
SELECT inserted.status, (blood_check.first_blood AND (sqlc.arg(status) = 'Correct')) AS first_blood FROM inserted, blood_check;
