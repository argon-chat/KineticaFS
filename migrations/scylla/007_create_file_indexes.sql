-- Create indexes for File table
CREATE INDEX IF NOT EXISTS file_bucketid_idx ON file (bucketid);
CREATE INDEX IF NOT EXISTS file_name_idx ON file (name);