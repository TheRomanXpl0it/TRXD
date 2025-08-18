-- name: UpdateFlag :exec
UPDATE flags
  SET
    flag = COALESCE(sqlc.narg('new_flag'), flag),
    regex = COALESCE(sqlc.narg('regex'), regex)
  WHERE chall_id = $1
    AND flag = $2;
