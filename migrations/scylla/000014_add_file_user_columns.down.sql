-- Remove JWT user information columns from File table
ALTER TABLE File DROP user_sub;
