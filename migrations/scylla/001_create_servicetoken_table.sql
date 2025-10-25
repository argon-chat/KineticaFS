-- Create ServiceToken table
CREATE TABLE IF NOT EXISTS ServiceToken (
    access_key text,
    token_type int,
    id text,
    created_at timestamp,
    updated_at timestamp,
    name text,
    PRIMARY KEY (id)
);