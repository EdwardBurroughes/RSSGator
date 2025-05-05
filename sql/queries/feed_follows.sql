-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4
    )
    RETURNING *
)
SELECT
    iff.*,
    f.name as feed_name,
    u.name as user_name
FROM inserted_feed_follow iff
JOIN feeds f ON iff.feed_id = f.id
JOIN users u ON iff.user_id = u.id;


-- name: GetFeedFollowsForUser :many
SELECT 
   f.name as feed_name
FROM feed_follows ff
JOIN feeds f ON ff.feed_id = f.id
WHERE ff.user_id = $1;


-- name: DeleteFeedFollowForUser :one
DELETE FROM feed_follows
WHERE feed_id = $1 and user_id = $2
RETURNING *
;