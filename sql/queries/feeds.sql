-- name: CreateFeed :one
INSERT INTO feeds (id, name, url, user_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetFeedByUrl :one
SELECT * FROM feeds WHERE url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds SET last_fetched_at = $1 WHERE id = $2;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds 
WHERE last_fetched_at IS NULL OR last_fetched_at < $1 
ORDER BY last_fetched_at NULLS FIRST 
LIMIT 1;