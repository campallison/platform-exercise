CREATE USER docker;
CREATE DATABASE docker;
GRANT ALL PRIVILEGES ON DATABASE docker TO docker;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE users (
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    name text NOT NULL,
    email text NOT NULL,
    password text NOT NULL,
    UNIQUE(email),
    PRIMARY KEY(id)
);