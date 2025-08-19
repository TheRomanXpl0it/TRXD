-- name: GetUserSolves :many
-- Retrieve all challenges solved by a user
SELECT s.chall_id, s.timestamp FROM submissions s
    WHERE s.user_id = $1
      AND s.status = 'Correct';
