// Package app provides application-level wiring and dependency injection
// following hexagonal architecture principles.
package app

import (
	"github.com/gorilla/sessions"
	"github.com/nats-io/nats.go"

	"github.com/yacobolo/datastar-go-starter-kit/internal/domain"
	"github.com/yacobolo/datastar-go-starter-kit/internal/features/todo/services"
	"github.com/yacobolo/datastar-go-starter-kit/internal/store"
)

// Repositories holds all repository implementations (driven adapters).
type Repositories struct {
	Todos    domain.TodoRepository
	Sessions domain.SessionRepository
}

// Services holds all service instances (application core).
type Services struct {
	Todo *services.TodoService
}

// App is the main application struct that holds all dependencies.
// This acts as the dependency injection container for the entire application.
type App struct {
	Store        *store.SQLiteStore
	SessionStore sessions.Store
	NATS         *nats.Conn
	Repositories *Repositories
	Services     *Services
}

// New creates a new App instance with all dependencies wired up.
// This follows the dependency injection pattern where:
// 1. Infrastructure (Store, NATS, SessionStore) is passed in
// 2. Repositories are created from the Store
// 3. Services are created with repository dependencies
func New(st *store.SQLiteStore, sessionStore sessions.Store, nc *nats.Conn) *App {
	// Create repositories (driven adapters)
	repos := &Repositories{
		Todos:    store.NewTodoRepository(st),
		Sessions: store.NewSessionRepository(st),
	}

	// Create services (application layer)
	// Services depend on domain interfaces, not concrete implementations
	svc := &Services{
		Todo: services.NewTodoService(repos.Todos, sessionStore),
	}

	return &App{
		Store:        st,
		SessionStore: sessionStore,
		NATS:         nc,
		Repositories: repos,
		Services:     svc,
	}
}

// Close closes all resources held by the application.
// This ensures graceful shutdown of all infrastructure components.
func (a *App) Close() error {
	if a.NATS != nil {
		a.NATS.Close()
	}
	if a.Store != nil {
		return a.Store.Close()
	}
	return nil
}
