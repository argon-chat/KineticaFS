-- Create indexes for ServiceToken table
CREATE INDEX IF NOT EXISTS servicetoken_name_idx ON servicetoken (name);
CREATE INDEX IF NOT EXISTS servicetoken_access_key_idx ON servicetoken (access_key);