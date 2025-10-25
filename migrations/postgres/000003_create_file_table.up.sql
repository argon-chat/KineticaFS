-- Create file table
CREATE TABLE IF NOT EXISTS file (
    id UUID NOT NULL,
    file_size_limit BIGINT,
    metadata TEXT,
    bucket_id UUID,
    file_size INTEGER,
    content_type TEXT,
    checksum TEXT,
    name TEXT,
    path TEXT,
    finalized BOOLEAN,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (id)
);
