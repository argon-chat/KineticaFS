-- Remove JWT user information columns from file table
ALTER TABLE file DROP COLUMN IF EXISTS user_sub;
ALTER TABLE file DROP COLUMN IF EXISTS space_id;
ALTER TABLE file DROP COLUMN IF EXISTS file_type;
