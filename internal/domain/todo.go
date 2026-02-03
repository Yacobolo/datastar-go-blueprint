// Package domain defines the core business interfaces following hexagonal architecture.
// These interfaces define ports that are implemented by adapters in the store layer.
package domain

import (
	"context"

	"github.com/yacobolo/datastar-go-starter-kit/internal/store/queries"
)

// TodoRepository defines the interface for todo data access.
// This is a port in hexagonal architecture, implemented by store adapters.
type TodoRepository interface {
	GetTodosByUser(ctx context.Context, userID string) ([]queries.Todo, error)
	CreateTodo(ctx context.Context, arg queries.CreateTodoParams) error
	DeleteAllTodosByUser(ctx context.Context, userID string) error
}

// SessionRepository defines the interface for session data access.
// This is a port in hexagonal architecture, implemented by store adapters.
type SessionRepository interface {
	GetSession(ctx context.Context, sessionID string) (queries.Session, error)
	UpsertSession(ctx context.Context, arg queries.UpsertSessionParams) error
}
