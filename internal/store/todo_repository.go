package store

import (
	"context"

	"github.com/yacobolo/datastar-go-starter-kit/internal/domain"
	"github.com/yacobolo/datastar-go-starter-kit/internal/store/queries"
)

// TodoRepository is the concrete implementation of domain.TodoRepository.
// It wraps sqlc-generated queries and acts as a driven adapter in hexagonal architecture.
type TodoRepository struct {
	store *SQLiteStore
}

// Ensure TodoRepository implements domain.TodoRepository at compile time.
var _ domain.TodoRepository = (*TodoRepository)(nil)

// NewTodoRepository creates a new TodoRepository instance.
func NewTodoRepository(st *SQLiteStore) *TodoRepository {
	return &TodoRepository{store: st}
}

// GetTodosByUser retrieves all todos for a given user ID.
func (r *TodoRepository) GetTodosByUser(ctx context.Context, userID string) ([]queries.Todo, error) {
	return r.store.Queries().GetTodosByUser(ctx, userID)
}

// CreateTodo creates a new todo in the database.
func (r *TodoRepository) CreateTodo(ctx context.Context, arg queries.CreateTodoParams) error {
	return r.store.Queries().CreateTodo(ctx, arg)
}

// DeleteAllTodosByUser deletes all todos for a given user ID.
func (r *TodoRepository) DeleteAllTodosByUser(ctx context.Context, userID string) error {
	return r.store.Queries().DeleteAllTodosByUser(ctx, userID)
}
