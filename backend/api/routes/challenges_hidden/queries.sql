-- name: ToggleChallengesHidden :exec
UPDATE challenges
  SET hidden = NOT hidden
  WHERE id = ANY(sqlc.arg('chall_ids')::INTEGER[]);
