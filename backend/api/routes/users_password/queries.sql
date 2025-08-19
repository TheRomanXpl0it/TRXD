-- name: ResetUserPassword :exec
-- Reset a user's password to a new password
UPDATE users SET password_hash = $1 WHERE id = $2;
