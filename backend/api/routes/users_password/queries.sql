-- name: ResetUserPassword :exec
-- Reset a user's password to a new password
UPDATE users SET password_hash = $2, password_salt = $3 WHERE id = $1;
