-- name: GetCategory :one
-- fetch category by name
SELECT * FROM categories WHERE name = $1;

-- name: UpdateChallengesCategory :exec
-- update category name in challenges table
UPDATE challenges SET category = sqlc.arg(new_category) WHERE category = sqlc.arg(old_category);
