-- +goose Up
CREATE TABLE IF NOT EXISTS files (
    id BIGSERIAL PRIMARY KEY,
    book_id BIGINT NOT NULL,
    format TEXT NOT NULL,
    storage_key TEXT NOT NULL UNIQUE,
    size_bytes BIGINT NOT NULL,
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    checksum TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_files_book_id ON files(book_id);

-- +goose Down
DROP TABLE IF EXISTS files;

