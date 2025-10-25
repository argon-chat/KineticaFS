-- Create Bucket table
CREATE TABLE IF NOT EXISTS Bucket (
    Endpoint text,
    AccessKey text,
    CustomConfig text,
    UpdatedAt timestamp,
    Name text,
    SecretKey text,
    UseSSL boolean,
    S3Provider text,
    StorageType int,
    ID text,
    CreatedAt timestamp,
    Region text,
    PRIMARY KEY (id)
);