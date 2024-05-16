-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPostByUser :many
SELECT posts.title, posts.url, posts.description, posts.published_at
FROM posts
INNER JOIN feeds_users
ON posts.feed_id = feeds_users.feed_id
WHERE feeds_users.user_id = $1
LIMIT $2;