-- +goose Up
CREATE TABLE IF NOT EXISTS book_ratings (
    book_id BIGINT PRIMARY KEY,
    rating_sum BIGINT NOT NULL,
    rating_count BIGINT NOT NULL,
    average DOUBLE PRECISION NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS book_ratings;

