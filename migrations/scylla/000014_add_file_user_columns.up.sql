-- Add JWT user information columns to File table
ALTER TABLE File ADD user_sub text;
