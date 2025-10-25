-- Add indices for service_token table
CREATE INDEX IF NOT EXISTS service_token_name_idx ON service_token (name);
CREATE INDEX IF NOT EXISTS service_token_access_key_idx ON service_token (access_key);
