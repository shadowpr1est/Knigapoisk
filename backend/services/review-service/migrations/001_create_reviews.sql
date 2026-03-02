-- +goose Up
CREATE TABLE IF NOT EXISTS reviews (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    book_id BIGINT NOT NULL,
    rating SMALLINT NOT NULL,
    text TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, book_id)
);

CREATE INDEX IF NOT EXISTS idx_reviews_book_id ON reviews(book_id);

-- +goose Down
DROP TABLE IF EXISTS reviews;

