-- Create File table
CREATE TABLE IF NOT EXISTS File (
    id text,
    file_size_limit varint,
    metadata text,
    bucket_id text,
    file_size int,
    content_type text,
    ref_count bigint,
    checksum text,
    name text,
    path text,
    finalized boolean,
    created_at timestamp,
    updated_at timestamp,
    PRIMARY KEY (id)
);