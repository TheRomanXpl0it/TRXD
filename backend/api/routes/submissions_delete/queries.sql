-- name: DeleteSubmission :exec
-- Delete a submission by its ID
DELETE FROM submissions WHERE id = $1;
