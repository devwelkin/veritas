-- name: CreateURL :one
INSERT INTO urls (original_url) VALUES ($1) RETURNING id;

-- name: UpdateShortCode :exec
UPDATE urls SET short_code = $1 WHERE id = $2;

-- name: GetURLByShortCode :one
SELECT original_url FROM urls WHERE short_code = $1;