// Package app provides application-level wiring and dependency injection
// following hexagonal architecture principles.
package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	embeddednats "github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"

	"github.com/yacobolo/datastar-go-starter-kit/internal/config"
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
	NATSServer   *embeddednats.Server
	Repositories *Repositories
	Services     *Services
}

// New creates a new App instance with all dependencies wired up.
// This function initializes all infrastructure components (SessionStore, NATS, Database)
// and wires up the application following hexagonal architecture:
// 1. Initialize infrastructure (SessionStore, NATS server, NATS client, Database)
// 2. Create repositories (driven adapters) from the Store
// 3. Create services (application layer) with repository dependencies
func New(cfg *config.Config) (*App, error) {
	// 1. Create SessionStore
	sessionStore := sessions.NewCookieStore([]byte(cfg.SessionSecret))
	sessionStore.MaxAge(86400 * 30) // 30 days
	sessionStore.Options.Path = "/"
	sessionStore.Options.HttpOnly = true
	sessionStore.Options.Secure = false
	sessionStore.Options.SameSite = http.SameSiteLaxMode

	// 2. Start embedded NATS server
	natsOpts := &embeddednats.Options{
		Host:      "localhost",
		Port:      4222,
		JetStream: true,
	}
	ns, err := embeddednats.NewServer(natsOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to start NATS: %w", err)
	}
	go ns.Start()
	if !ns.ReadyForConnections(4 * time.Second) {
		return nil, fmt.Errorf("NATS not ready")
	}

	slog.Info("NATS server started", "url", ns.ClientURL())

	// 3. Connect to NATS
	nc, err := nats.Connect(ns.ClientURL())
	if err != nil {
		ns.Shutdown()
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	// 4. Open database
	dbStore, err := store.Open(cfg.DBPath)
	if err != nil {
		nc.Close()
		ns.Shutdown()
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	slog.Info("database initialized", "path", cfg.DBPath)

	// 5. Create repositories (driven adapters)
	repos := &Repositories{
		Todos:    store.NewTodoRepository(dbStore),
		Sessions: store.NewSessionRepository(dbStore),
	}

	// 6. Create services (application layer)
	// Services depend on domain interfaces, not concrete implementations
	svc := &Services{
		Todo: services.NewTodoService(repos.Todos, repos.Sessions, sessionStore),
	}

	return &App{
		Store:        dbStore,
		SessionStore: sessionStore,
		NATS:         nc,
		NATSServer:   ns,
		Repositories: repos,
		Services:     svc,
	}, nil
}

// Close closes all resources held by the application.
// This ensures graceful shutdown of all infrastructure components.
func (a *App) Close() error {
	if a.NATS != nil {
		a.NATS.Close()
	}
	if a.NATSServer != nil {
		a.NATSServer.Shutdown()
	}
	if a.Store != nil {
		return a.Store.Close()
	}
	return nil
}
