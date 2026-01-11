-- Add indices for file table
CREATE INDEX IF NOT EXISTS file_bucket_id_idx ON file (bucket_id);
CREATE INDEX IF NOT EXISTS file_name_idx ON file (name);
