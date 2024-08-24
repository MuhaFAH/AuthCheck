CREATE TABLE users (
    user_id UUID PRIMARY KEY,
    last_ip VARCHAR(64) NOT NULL,
    token_hash TEXT NOT NULL
);