-- name: RegisterTeam :exec
-- Insert a new team and add the founder user to the team
WITH locked_user AS (
    SELECT id FROM users
    WHERE id = $1 AND team_id IS NULL
    FOR UPDATE
  ),
  new_team AS (
    INSERT INTO teams (name, password_hash)
    SELECT $2, $3
    FROM locked_user
    RETURNING *
  )
UPDATE users
  SET team_id = new_team.id
  FROM new_team
  WHERE users.id = $1;
