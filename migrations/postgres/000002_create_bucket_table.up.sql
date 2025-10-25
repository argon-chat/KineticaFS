-- Create bucket table
CREATE TABLE IF NOT EXISTS bucket (
    endpoint TEXT,
    access_key TEXT,
    custom_config TEXT,
    updated_at TIMESTAMP,
    name TEXT,
    secret_key TEXT,
    use_ssl BOOLEAN,
    s3_provider TEXT,
    storage_type INTEGER,
    id UUID NOT NULL,
    created_at TIMESTAMP,
    region TEXT,
    PRIMARY KEY (id)
);
