-- Create File table
CREATE TABLE IF NOT EXISTS File (
    file_size_limit varint,
    metadata text,
    id text,
    created_at timestamp,
    updated_at timestamp,
    bucket_id text,
    file_size int,
    content_type text,
    checksum text,
    name text,
    path text,
    finalized boolean,
    PRIMARY KEY (id)
);