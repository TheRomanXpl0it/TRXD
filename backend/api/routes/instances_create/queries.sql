-- name: GetInstance :one
-- Gets an instance by ID
SELECT * FROM instances WHERE chall_id = $1 AND team_id = $2;

-- name: CreateInstance :exec
-- Creates a new instance for a team
INSERT INTO instances (team_id, chall_id, expires_at, host, port)
	VALUES ($1, $2, $3, sqlc.narg(host), $4);