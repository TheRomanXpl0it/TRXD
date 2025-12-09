-- name: CreateCategory :exec
-- Insert a new category
INSERT INTO categories (name) VALUES ($1);
