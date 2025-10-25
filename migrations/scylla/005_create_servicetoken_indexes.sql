-- Create indexes for ServiceToken table
CREATE INDEX IF NOT EXISTS servicetoken_name_idx ON servicetoken (name);
CREATE INDEX IF NOT EXISTS servicetoken_accesskey_idx ON servicetoken (accesskey);