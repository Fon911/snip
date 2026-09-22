-- name: SaveLink :exec
INSERT INTO links (code, url)
VALUES ($1, $2);
-- name: GetLink :one
SELECT url
FROM links
WHERE code = $1;
-- name: IncrementClicks :exec
UPDATE links
SET clicks = clicks + 1
WHERE code = $1;
-- name: GetStats :one
SELECT url,
  clicks
FROM links
WHERE code = $1;