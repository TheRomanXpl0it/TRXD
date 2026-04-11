-- name: GetAdminStats :one
SELECT
  (SELECT COUNT(*) FROM users) AS total_users,
  (SELECT COUNT(*) FROM users WHERE role='Player') AS total_players,
  (SELECT COUNT(*) FROM teams) AS total_teams,
  (SELECT COUNT(*) FROM challenges) AS total_challenges,
  (SELECT COUNT(*) FROM challenges WHERE hidden=FALSE) AS total_released_challenges,
  (SELECT COUNT(*) FROM submissions) AS total_submissions,
  (SELECT COUNT(*) FROM submissions WHERE status='Correct') AS total_correct_submissions;
