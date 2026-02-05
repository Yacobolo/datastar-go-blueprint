# Agent Guidelines for Datastar Go Blueprint

This document provides essential information for AI coding agents working on this codebase.

## Build, Lint, and Test Commands

### Running Tests
```bash
# Run all tests (package summary format)
task test
# or
gotestsum --format pkgname-and-test-fails -- ./...

# Run tests with coverage
task test:coverage

# Run a single test
go test -run TestFunctionName ./path/to/package

# Run tests in specific package
go test ./internal/features/todo
```

### Linting
```bash
# Run Go linters (golangci-lint)
task lint

# Run CSS linting (checks for unused/invalid CSS classes)
task css:lint

# Run all checks (tests + lint + css lint)
task check
```

### Building
```bash
# Full production build (with code generation + asset bundling)
task build

# Development mode (hot reload with Air, Templ watcher, esbuild)
task dev

# Stop development processes
task dev:stop

# Code generation only
task generate:all       # All generators (sqlc, templ, cssgen)
task generate:templ     # Generate templ templates
task generate:sqlc      # Generate type-safe DB queries
task generate:css       # Generate type-safe CSS constants
```

### Asset Bundling
```bash
# Bundle web assets (CSS + JS/TS with esbuild)
task build:web:assets
```

## Code Style Guidelines

### Import Organization
Organize imports in three groups with blank lines between:
1. Standard library (alphabetical)
2. External packages (alphabetical)
3. Internal packages (alphabetical)

```go
import (
    "context"
    "fmt"
    "log/slog"

    "github.com/go-chi/chi/v5"
    "github.com/google/uuid"

    "github.com/yacobolo/datastar-go-blueprint/internal/domain"
    "github.com/yacobolo/datastar-go-blueprint/internal/store"
)
```

Use import aliases for clarity:
- `ds "github.com/Yacobolo/datastar-templ"` - Datastar attributes in templ files
- Package aliases to avoid conflicts: `commoncomponents`, `todocomponents`

### Naming Conventions

**Types:**
- Exported: `PascalCase` - `Handlers`, `TodoService`, `TodoRepository`
- Interfaces: `PascalCase` - `TodoRepository`, `SessionRepository`
- Unexported: `camelCase` - `txKey`

**Functions:**
- Exported: `PascalCase` - `IndexPage`, `GetTodosByUser`, `NewHandlers`
- Constructors: `New` + type name - `NewHandlers`, `NewTodoService`
- Unexported: `camelCase` - `subject`, `refreshTodos`, `notifyUpdate`

**Variables:**
- Constants (exported): `PascalCase` - `Dev`, `Prod`, `ToastSuccess`
- Constants (SQL/internal): `SCREAMING_SNAKE_CASE` (sqlc generated)
- Variables: `camelCase` (unexported), `PascalCase` (exported)

**Packages:**
- Singular lowercase nouns: `app`, `config`, `domain`, `store`, `todo`
- Exception: `services`, `queries` (grouping/generated)

### Struct Organization
Order fields by importance: configuration → core dependencies → infrastructure → state

```go
type Handlers struct {
    logger       *slog.Logger              // 1. Observability
    todoService  *services.TodoService     // 2. Core business
    nats         *nats.Conn                // 3. Infrastructure
    sessionStore sessions.Store            // 4. Infrastructure
}
```

### Error Handling

**Always wrap errors with context using `fmt.Errorf` and `%w`:**
```go
if err != nil {
    return fmt.Errorf("failed to get todos: %w", err)
}
```

**Check specific errors with `errors.Is`:**
```go
if err != nil && !errors.Is(err, sql.ErrNoRows) {
    return nil, fmt.Errorf("failed to get todos: %w", err)
}
```

**Join multiple errors with `errors.Join`:**
```go
if rbErr := tx.Rollback(); rbErr != nil {
    return errors.Join(err, fmt.Errorf("rollback failed: %w", rbErr))
}
```

**HTTP error responses:**
```go
if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
}
```

### Documentation Comments

**Follow Godoc conventions:**
- Full sentences starting with the item name
- Package comment at top of main file
- Exported items must have comments

```go
// Package todo implements the todo feature handlers and routes.
package todo

// Handlers holds dependencies for todo HTTP handlers.
type Handlers struct { ... }

// NewHandlers creates a new Handlers instance with the given dependencies.
func NewHandlers(...) *Handlers { ... }

// RequireSession retrieves or creates a session ID for the current request.
func RequireSession(...) (string, bool) { ... }
```

**Inline comments for complex logic:**
```go
// 1. Create SessionStore
sessionStore := sessions.NewCookieStore([]byte(cfg.SessionSecret))

// 2. Start embedded NATS server
natsOpts := &embeddednats.Options{...}
```

**Nolint comments must explain WHY:**
```go
//nolint:unused // Reserved for future transaction support
func (s *SQLiteStore) conn(ctx context.Context) *queries.Queries {
```

## Architecture Patterns

### Hexagonal Architecture (Ports & Adapters)
This codebase follows hexagonal architecture:

**Domain (Core):**
- `internal/domain/` - Define interfaces (ports)
- `internal/features/*/services/` - Business logic

**Adapters:**
- **Driving** (Primary): `internal/features/*/handlers.go`, `routes.go`
- **Driven** (Secondary): `internal/store/*_repository.go`

**Infrastructure:**
- `internal/platform/` - Shared infrastructure (router, pubsub)
- `internal/app/` - Dependency injection container

### Repository-Service-Handler Pattern

**Repositories** (Data Layer):
```go
// 1. Define interface in domain/
type TodoRepository interface {
    GetTodosByUser(ctx context.Context, userID string) ([]queries.Todo, error)
}

// 2. Implement in store/
type TodoRepository struct {
    store *SQLiteStore
}

// 3. Compile-time interface check
var _ domain.TodoRepository = (*TodoRepository)(nil)

// 4. Constructor
func NewTodoRepository(st *SQLiteStore) *TodoRepository {
    return &TodoRepository{store: st}
}
```

**Services** (Business Logic):
```go
type TodoService struct {
    todoRepo    domain.TodoRepository    // Depend on interfaces
    sessionRepo domain.SessionRepository
}

func NewTodoService(todoRepo domain.TodoRepository, ...) *TodoService {
    return &TodoService{todoRepo: todoRepo, ...}
}
```

**Handlers** (HTTP Layer):
```go
type Handlers struct {
    logger      *slog.Logger
    todoService *services.TodoService
}

func NewHandlers(logger *slog.Logger, svc *services.TodoService, ...) *Handlers {
    return &Handlers{logger: logger, todoService: svc, ...}
}

// Handler signature
func (h *Handlers) IndexPage(w http.ResponseWriter, r *http.Request) {
    // Extract params → Get session → Call service → Respond
}
```

### Feature-Based Organization
Organize code by feature (vertical slices), not layer:

```
internal/features/todo/
├── components/        # UI components (templ)
├── pages/             # Full page layouts (templ)
├── services/          # Business logic
├── handlers.go        # HTTP handlers
└── routes.go          # Route registration
```

### Database Patterns

**SQLC for type-safe queries:**
```sql
-- name: GetTodosByUser :many
SELECT * FROM todos WHERE user_id = ? ORDER BY created_at DESC;
```

**Repository implementation:**
```go
func (r *TodoRepository) GetTodosByUser(ctx context.Context, userID string) ([]queries.Todo, error) {
    return r.store.Queries().GetTodosByUser(ctx, userID)
}
```

**Migrations with Goose + embed:**
```go
//go:embed migrations/*.sql
var migrations embed.FS
```

**Transaction pattern:**
```go
err := store.WithinTransaction(ctx, func(txCtx context.Context) error {
    // Use txCtx for all operations within transaction
    return nil
})
```

### Handler Pattern
Standard flow: Extract params → Get session → Call service → Notify → Respond

```go
func (h *Handlers) ToggleTodo(w http.ResponseWriter, r *http.Request) {
    // 1. Get session
    sessionID, ok := RequireSession(h.sessionStore, w, r)
    if !ok {
        return
    }

    // 2. Extract URL parameter
    idx, ok := RequireIntParam(w, r, "idx")
    if !ok {
        return
    }

    // 3. Call business logic
    _, mvc, err := h.todoService.GetSessionMVC(w, r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    h.todoService.ToggleTodo(mvc, idx)
    
    // 4. Persist
    if err := h.todoService.SaveMVC(r.Context(), sessionID, mvc); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // 5. Notify update (triggers SSE)
    h.notifyUpdate(sessionID, pubsub.WithRefresh())
    w.WriteHeader(http.StatusOK)
}
```

### Templ Component Patterns

**Type-safe CSS classes from generated constants:**
```go
class={ ui.Btn, ui.BtnLg, ui.BtnPrimary }
```

**Datastar attributes (import as `ds`):**
```go
import ds "github.com/Yacobolo/datastar-templ"

// Signals
{ ds.Signals(ds.String("input", input))... }
{ ds.Bind("input")... }

// Events
{ ds.OnClick(ds.Post("/api/todos/%d/toggle", i))... }
{ ds.OnKeyDown(...)... }

// HTTP methods
{ ds.Get("/api/todos/updates")... }
{ ds.Post("/api/todos/%d", i)... }
{ ds.Put("/api/todos/reset")... }
{ ds.Delete("/api/todos/%d", i)... }

// Conditional rendering
{ ds.Show("$theme === 'light'")... }
```

**Component composition:**
```go
templ Base(title string) {
    <body>
        { children... }  // Child content injection
    </body>
}

// Usage
@layouts.Base("My Page") {
    <div>Content here</div>
}
```

## Development Workflow

### When modifying files:

**Templ files** (`.templ`):
- Generate with `task generate:templ` or `templ generate`
- Auto-watched in dev mode (`task dev`)
- Creates `*_templ.go` files (don't edit these)

**SQL queries** (`internal/store/queries/*.sql`):
- Generate with `task generate:sqlc` or `sqlc generate`
- Creates type-safe Go code in `internal/store/queries/`

**CSS files** (`web/ui/styles/**/*.css`):
- Generate with `task generate:css`
- Creates type-safe constants in `internal/ui/styles*.gen.go`
- Use constants in templ: `class={ ui.Btn }`

**Migrations** (`internal/store/migrations/*.sql`):
- Name format: `NNN_description.sql`
- Include `-- +goose Up` and `-- +goose Down`
- Automatically run on app startup

### Environment-specific code:
Use build tags for dev/prod separation:
```go
//go:build dev
// or
//go:build !dev
```

## Logging

**Use structured logging with `log/slog`:**
```go
logger.Info("server started", "port", cfg.Port)
logger.Error("failed to get todos", "error", err)
logger.Debug("processing request", "session_id", sessionID)
```

**Never use:**
- `log` package (forbidden by depguard)
- `logrus`, `zap` (use slog only)

## Configuration

**Global config singleton:**
```go
config.Global.Environment  // Dev or Prod
config.Global.Port
config.Global.DBPath
```

**Environment variables:**
- See `.env.example` for available variables
- Defaults in `internal/config/config.go`
