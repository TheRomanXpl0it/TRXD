-- name: CreateCategory :exec
-- Insert a new category
INSERT INTO categories (name, icon) VALUES ($1, $2);
