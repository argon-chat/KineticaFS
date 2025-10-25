-- Create file_blob table
CREATE TABLE IF NOT EXISTS file_blob (
    id UUID NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    file_id UUID,
    PRIMARY KEY (id)
);
