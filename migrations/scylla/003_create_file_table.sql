-- Create File table
CREATE TABLE IF NOT EXISTS File (
    FileSizeLimit varint,
    Metadata text,
    ID text,
    CreatedAt timestamp,
    UpdatedAt timestamp,
    BucketID text,
    FileSize int,
    ContentType text,
    Checksum text,
    Name text,
    Path text,
    Finalized boolean,
    PRIMARY KEY (id)
);