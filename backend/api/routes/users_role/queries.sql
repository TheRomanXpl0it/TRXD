-- name: ChangeUserRole :exec
-- Changes the role of a user
UPDATE users SET role = $2 WHERE id = $1;
