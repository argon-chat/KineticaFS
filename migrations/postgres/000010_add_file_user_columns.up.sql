-- Add JWT user information columns to file table
ALTER TABLE file ADD COLUMN user_sub UUID;
ALTER TABLE file ADD COLUMN space_id UUID;
ALTER TABLE file ADD COLUMN file_type TEXT;
