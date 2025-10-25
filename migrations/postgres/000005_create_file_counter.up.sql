-- Create file_counter table
-- Note: PostgreSQL doesn't have a native counter type like Scylla
-- Using INTEGER with proper application-level management for reference counting
CREATE TABLE IF NOT EXISTS file_counter (
    id UUID NOT NULL,
    ref INTEGER DEFAULT 0,
    PRIMARY KEY (id)
);
