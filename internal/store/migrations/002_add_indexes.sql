-- +goose Up
-- Index for faster user queries
CREATE INDEX IF NOT EXISTS idx_todos_user_id ON todos(user_id);

-- Index for created_at for ordering
CREATE INDEX IF NOT EXISTS idx_todos_created_at ON todos(created_at DESC);

-- Index for completed status
CREATE INDEX IF NOT EXISTS idx_todos_completed ON todos(completed);

-- +goose Down
DROP INDEX IF EXISTS idx_todos_user_id;
DROP INDEX IF EXISTS idx_todos_created_at;
DROP INDEX IF EXISTS idx_todos_completed;
