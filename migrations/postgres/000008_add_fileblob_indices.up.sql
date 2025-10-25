-- Add indices for file_blob table
CREATE INDEX IF NOT EXISTS file_blob_file_id_idx ON file_blob (file_id);
