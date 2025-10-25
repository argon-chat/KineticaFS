-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
-- Create service_token table
CREATE TABLE IF NOT EXISTS service_token (
    access_key TEXT,
    token_type INTEGER,
    id UUID NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    name TEXT,
    PRIMARY KEY (id)
);
