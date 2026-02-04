-- name: GetTodosByUser :many
SELECT * FROM todos 
WHERE user_id = ? 
ORDER BY created_at DESC;

-- name: GetTodoByID :one
SELECT * FROM todos 
WHERE id = ?;

-- name: CreateTodo :exec
INSERT INTO todos (id, user_id, task, completed) 
VALUES (?, ?, ?, ?);

-- name: UpdateTodoTask :exec
UPDATE todos 
SET task = ?, updated_at = CURRENT_TIMESTAMP 
WHERE id = ?;

-- name: ToggleTodoCompleted :exec
UPDATE todos 
SET completed = CASE WHEN completed = 0 THEN 1 ELSE 0 END, 
    updated_at = CURRENT_TIMESTAMP 
WHERE id = ?;

-- name: DeleteTodo :exec
DELETE FROM todos 
WHERE id = ?;

-- name: DeleteAllTodosByUser :exec
DELETE FROM todos 
WHERE user_id = ?;

-- name: CountTodosByUser :one
SELECT COUNT(*) as count 
FROM todos 
WHERE user_id = ?;

-- Session queries
-- name: GetSession :one
SELECT * FROM sessions WHERE id = ?;

-- name: UpsertSession :exec
INSERT INTO sessions (id, data, mode, editing_idx, updated_at)
VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
ON CONFLICT(id) DO UPDATE SET
    data = excluded.data,
    mode = excluded.mode,
    editing_idx = excluded.editing_idx,
    updated_at = CURRENT_TIMESTAMP;

-- name: DeleteSession :exec
DELETE FROM sessions WHERE id = ?;
