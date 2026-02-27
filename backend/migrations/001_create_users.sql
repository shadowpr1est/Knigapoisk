-- +goose Up
CREATE TABLE users (
                       id            BIGSERIAL PRIMARY KEY,
                       email         VARCHAR(255) UNIQUE NOT NULL,
                       username      VARCHAR(100) UNIQUE NOT NULL,
                       password_hash VARCHAR(255) NOT NULL,
                       role          VARCHAR(50) NOT NULL DEFAULT 'user',
                       created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
                       updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_users_email    ON users(email);
CREATE INDEX idx_users_username ON users(username);

-- +goose Down
DROP TABLE users;