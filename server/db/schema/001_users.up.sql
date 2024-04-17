CREATE TABLE IF NOT EXISTS users(
    id uuid PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    hashed_password VARCHAR(256) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);