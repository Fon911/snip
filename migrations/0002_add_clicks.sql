-- +goose Up
ALTER TABLE links
ADD COLUMN clicks INT NOT NULL DEFAULT 0;
-- +goose Down
ALTER TABLE links DROP COLUMN clicks;