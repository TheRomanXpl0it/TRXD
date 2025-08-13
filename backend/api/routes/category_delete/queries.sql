-- name: DeleteCategory :exec
-- Delete a category and all associated challenges
DELETE FROM categories WHERE name = $1;
