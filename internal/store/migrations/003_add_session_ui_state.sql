-- +goose Up
ALTER TABLE sessions ADD COLUMN mode INTEGER DEFAULT 0;
ALTER TABLE sessions ADD COLUMN editing_idx INTEGER DEFAULT -1;

-- +goose Down
ALTER TABLE sessions DROP COLUMN editing_idx;
ALTER TABLE sessions DROP COLUMN mode;
