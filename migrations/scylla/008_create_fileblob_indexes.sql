-- Create indexes for FileBlob table
CREATE INDEX IF NOT EXISTS fileblob_file_id_idx ON fileblob (file_id);