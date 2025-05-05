-- name: CreateFeed :one
INSERT INTO feeds (name, url, user_id)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetFeeds :many
SELECT 
    f.name as feed_name,
    f.url,
    u.name as user_name 
from feeds f
join users u on f.user_id = u.id;


-- name: GetFeed :one 
SELECT * FROM feeds WHERE url = $1;

-- name: MarkFeedFetched :one
UPDATE feeds SET last_fetched_at = $1 WHERE id = $2
RETURNING name;

-- name: GetNextFeedToFetch :one
SELECT
    id,
    url,
    name
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;
