-- name: GetUserSolves :many
-- Retrieve all challenges solved by a user
SELECT c.id, c.name, c.category, s.first_blood, s.timestamp
  FROM submissions s
  JOIN challenges c ON s.chall_id = c.id
  WHERE s.user_id = $1
    AND s.status = 'Correct';
