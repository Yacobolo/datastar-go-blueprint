-- Todo items table
CREATE TABLE IF NOT EXISTS todos (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    task TEXT NOT NULL,
    completed INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for faster user queries
CREATE INDEX IF NOT EXISTS idx_todos_user_id ON todos(user_id);

-- Index for created_at for ordering
CREATE INDEX IF NOT EXISTS idx_todos_created_at ON todos(created_at DESC);
