CREATE TABLE users (
    user_id UUID PRIMARY KEY UNIQUE,
    last_ip VARCHAR(64) NOT NULL,
    token_hash TEXT NOT NULL
);