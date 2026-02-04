-- name: GetTeamsScoreboard :many
-- Retrieve all teams or a subset if specified, ordered by score and last correct submission time
SELECT
    t.id,
    t.name,
    t.score,
    t.country,
    COALESCE(b.badges, '[]') AS badges,
    lc.last_correct_at
  FROM teams t
  LEFT JOIN ( -- Badges per team
      SELECT
        team_id,
        JSON_AGG(
          JSON_BUILD_OBJECT(
            'name', name,
            'description', description
          )
        ) AS badges
      FROM badges
      GROUP BY team_id
    ) b ON b.team_id = t.id
  LEFT JOIN ( -- Last correct submission per team
      SELECT
        u.team_id,
        MAX(s.timestamp) AS last_correct_at
      FROM users u
      JOIN submissions s
          ON s.user_id = u.id
        AND s.status = 'Correct'
      WHERE u.role = 'Player'
      GROUP BY u.team_id
    ) lc ON lc.team_id = t.id
  ORDER BY
    t.score DESC,
    lc.last_correct_at ASC NULLS LAST
  OFFSET sqlc.arg('offset')
  LIMIT sqlc.narg('limit');
