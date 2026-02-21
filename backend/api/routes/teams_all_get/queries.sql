-- name: GetTeamsPreview :many
-- Retrieve all teams
SELECT
    t.id,
    t.name,
    COALESCE(u.top_role, 'Player')::TEXT AS user_role,
    t.score,
    t.country,
    COALESCE(
      JSON_AGG(
        JSON_BUILD_OBJECT(
          'name', b.name,
          'description', b.description
        )
      ) FILTER (WHERE b.name IS NOT NULL),
      '[]'
    ) AS badges
  FROM teams t
  LEFT JOIN (
    SELECT team_id,
        MAX(role) AS top_role
      FROM users
      WHERE role != 'Player'
      GROUP BY team_id
    ) u ON u.team_id = t.id
  LEFT JOIN badges b ON b.team_id = t.id
  GROUP BY t.id, t.name, t.score, t.country, u.top_role
  ORDER BY t.id
  OFFSET sqlc.arg('offset')
  LIMIT sqlc.narg('limit');
