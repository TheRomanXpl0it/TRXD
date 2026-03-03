-- name: GetTotalSubmissions :one
-- fetches total number of submissions
SELECT COUNT(*) FROM submissions;

-- name: GetSubmissions :many
-- fetches all submissions, with pagination
SELECT
    s.id,
    s.user_id,
    u.name AS user_name,
    t.id AS team_id,
    t.name AS team_name,
    s.chall_id,
    c.name AS chall_name,
    s.status,
    s.first_blood,
    s.flag,
    s.timestamp
  FROM submissions s
  JOIN users u ON s.user_id = u.id
  JOIN teams t ON u.team_id = t.id
  JOIN challenges c ON s.chall_id = c.id
  ORDER BY s.id DESC
  OFFSET sqlc.arg('offset')
  LIMIT sqlc.narg('limit');
