-- Set TTL for FileBlob table (10 minutes = 600 seconds)
ALTER TABLE fileblob WITH default_time_to_live = 600;