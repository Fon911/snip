-- +goose Up
CREATE TABLE links(
  id SERIAL PRIMARY KEY,
  code TEXT NOT NULL UNIQUE,
  url TEXT NOT NULL
);
-- +goose Down
DROP TABLE notes;