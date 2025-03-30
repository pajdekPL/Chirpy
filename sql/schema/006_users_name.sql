-- +goose Up
ALTER TABLE users 
ADD COLUMN user_name TEXT NOT NULL UNIQUE;

-- +goose Down
ALTER TABLE users
DROP COLUMN user_name;