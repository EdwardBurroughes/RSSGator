-- name: CreatePost :one
INSERT INTO posts (created_at, updated_at, title, url, published_at, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetPosts :many
SELECT
    p.title,
    p.url,
    p.description,
    p.published_at
FROM posts p
JOIN feed_follows ff on p.feed_id = ff.feed_id
WHERE ff.user_id = $1
LIMIT $2;