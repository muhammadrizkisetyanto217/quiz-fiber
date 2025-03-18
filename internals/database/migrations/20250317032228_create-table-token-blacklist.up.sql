CREATE TABLE IF NOT EXISTS token_blacklist (
    id SERIAL PRIMARY KEY,
    token TEXT NOT NULL UNIQUE,
    expired_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL
);
