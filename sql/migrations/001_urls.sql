-- +goose Up
-- +goose StatementBegin
CREATE TABLE urls (
    id BIGSERIAL PRIMARY KEY,
    short_code VARCHAR(255) UNIQUE,
    original_url TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Let's create an index on short_code for faster redirect service performance.
-- The UNIQUE constraint already creates a B-tree index, but it's still good practice to specify it.
CREATE INDEX IF NOT EXISTS idx_urls_short_code ON urls(short_code);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_urls_short_code;
DROP TABLE IF EXISTS urls;
-- +goose StatementEnd