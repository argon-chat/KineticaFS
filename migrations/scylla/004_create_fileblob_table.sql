-- Create FileBlob table
CREATE TABLE IF NOT EXISTS FileBlob (
    ID text,
    CreatedAt timestamp,
    UpdatedAt timestamp,
    FileID text,
    PRIMARY KEY (id)
);