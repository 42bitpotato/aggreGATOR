-- name: CreatePost :exec
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, published_raw, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9
);

-- name: GetPostsForUser :many
SELECT * FROM posts
WHERE feed_id IN (
    SELECT feed_id
    FROM feed_follows
    INNER JOIN feeds ON feed_follows.feed_id = feeds.id
    INNER JOIN users ON feed_follows.user_id = users.id
    WHERE users.id = $1
)
ORDER BY published_at DESC NULLS LAST
LIMIT $2;