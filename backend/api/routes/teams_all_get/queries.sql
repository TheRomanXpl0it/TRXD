-- name: GetTeamsPreview :many
-- Retrieve all teams
SELECT
    t.id,
    t.name,
    t.score,
    t.country,
    t.image,
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
  GROUP BY t.id, t.name, t.score, t.country, t.image
  ORDER BY t.id;
