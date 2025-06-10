-- +goose Up
-- +goose StatementBegin
CREATE TABLE urls (
    id BIGSERIAL PRIMARY KEY,
    short_code VARCHAR(255) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Redirect servisinin hızlı çalışması için short_code üzerinde bir index oluşturalım.
-- UNIQUE kısıtlaması zaten bir B-tree index oluşturur, ancak bunu belirtmek yine de iyi bir pratiktir.
CREATE INDEX IF NOT EXISTS idx_urls_short_code ON urls(short_code);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_urls_short_code;
DROP TABLE IF EXISTS urls;
-- +goose StatementEnd