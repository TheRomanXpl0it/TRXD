-- name: ResetTeamPassword :exec
-- Reset a team's password to a new password
UPDATE teams SET password_hash = $1 WHERE id = $2;
