-- Create FileBlob table
CREATE TABLE IF NOT EXISTS FileBlob (
    id text,
    created_at timestamp,
    updated_at timestamp,
    file_id text,
    PRIMARY KEY (id)
);