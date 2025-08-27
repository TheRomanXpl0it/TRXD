-- name: GetInstance :one
-- Gets an instance by ID
SELECT * FROM instances WHERE chall_id = $1 AND team_id = $2;
