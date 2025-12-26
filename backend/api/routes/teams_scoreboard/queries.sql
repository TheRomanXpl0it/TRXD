-- name: GetTeamsScoreboard :many
-- Retrieve all teams
SELECT
    t.id,
    t.name,
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
  LEFT JOIN badges b ON b.team_id = t.id
  GROUP BY t.id, t.name, t.score, t.country
  ORDER BY t.score DESC
  OFFSET sqlc.arg('offset')
  LIMIT sqlc.narg('limit');
