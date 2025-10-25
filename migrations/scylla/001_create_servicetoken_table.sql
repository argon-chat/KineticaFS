-- Create ServiceToken table
CREATE TABLE IF NOT EXISTS ServiceToken (
    AccessKey text,
    TokenType int,
    ID text,
    CreatedAt timestamp,
    UpdatedAt timestamp,
    Name text,
    PRIMARY KEY (id)
);