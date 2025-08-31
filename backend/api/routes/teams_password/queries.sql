-- name: ResetTeamPassword :exec
-- Reset a team's password to a new password
UPDATE teams SET password_hash = $2, password_salt = $3 WHERE id = $1;
