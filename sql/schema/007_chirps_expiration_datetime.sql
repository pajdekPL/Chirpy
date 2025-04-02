-- +goose Up
ALTER TABLE chirps 
ADD COLUMN expiration_datetime TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (NOW() + INTERVAL '10 years');

-- +goose Down
ALTER TABLE chirps
DROP COLUMN expiration_datetime;