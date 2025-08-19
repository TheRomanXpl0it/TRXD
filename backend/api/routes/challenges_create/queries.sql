-- name: CreateChallenge :one
-- Insert a new challenge
INSERT INTO challenges (name, category, description, type, max_points, score_type)
  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;
