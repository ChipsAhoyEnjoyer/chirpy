-- +goose Up
CREATE TABLE users(
    id UUID PRIMARY KEY NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email TEXT NOT NULL,
    CONSTRAINT unique_email UNIQUE (email)
);

-- +goose Down
DROP TABLE users;