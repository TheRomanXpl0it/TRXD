-- name: AddTeamMember :exec
-- Assign a user to a team
UPDATE users SET team_id = $1 WHERE id = $2 AND team_id IS NULL;
