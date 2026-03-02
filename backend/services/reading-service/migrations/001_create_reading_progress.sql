-- +goose Up
CREATE TABLE IF NOT EXISTS reading_progress (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    book_id BIGINT NOT NULL,
    file_id BIGINT NOT NULL,
    page INT NOT NULL,
    percentage DOUBLE PRECISION NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, book_id)
);

-- +goose Down
DROP TABLE IF EXISTS reading_progress;

