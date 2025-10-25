-- Create Bucket table
CREATE TABLE IF NOT EXISTS Bucket (
    endpoint text,
    access_key text,
    custom_config text,
    updated_at timestamp,
    name text,
    secret_key text,
    use_ssl boolean,
    s3_provider text,
    storage_type int,
    id text,
    created_at timestamp,
    region text,
    PRIMARY KEY (id)
);