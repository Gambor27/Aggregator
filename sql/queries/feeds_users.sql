-- name: CreateFeedFollow :one
INSERT INTO feeds_users (id, created_at, updated_at, feed_id, user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserFeeds :many
SELECT * FROM feeds_users
WHERE user_id = $1;

-- name: GetFollow :one
SELECT * FROM feeds_users
WHERE id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feeds_users
WHERE id = $1;