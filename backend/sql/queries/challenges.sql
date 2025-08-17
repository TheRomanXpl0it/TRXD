-- name: GetChallengeByID :one
-- Retrieve a challenge by its ID
SELECT * FROM challenges WHERE id = $1;

-- name: GetTagsByChallenge :many
-- Retrieve all tags associated with a challenge
SELECT name FROM tags WHERE chall_id = $1;
