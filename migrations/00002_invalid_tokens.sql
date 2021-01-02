-- +goose Up
CREATE TABLE invalid_tokens (
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    token text,
    PRIMARY KEY (token)
);

-- +goose Down
DROP TABLE invalid_tokens;