-- Create indexes for Bucket table
CREATE INDEX IF NOT EXISTS bucket_name_idx ON bucket (name);
CREATE INDEX IF NOT EXISTS bucket_region_idx ON bucket (region);