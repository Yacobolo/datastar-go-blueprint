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
go test ./internal/ui/todo
```

### Linting
```bash
# Run Go linters (golangci-lint)
task lint

# Run all checks (tests + lint + asset builds)
task check
```

### Building
```bash
# Full production build (with code generation + asset bundling)
task build

# Development mode (hot reload with Air, esbuild, and Tailwind)
task dev

# Stop development processes
task dev:stop

# Code generation only
task generate:all       # All generators (sqlc)
task generate:sqlc      # Generate type-safe DB queries
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
- Dot imports are acceptable for gomponents packages in UI files:
  - `. "maragu.dev/gomponents"`
  - `. "maragu.dev/gomponents/html"`
- `data "maragu.dev/gomponents-datastar"` - Datastar attributes in gomponents views
- Prefer short aliases only when they add clarity, e.g. `core` for `internal/ui/core`

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
- `internal/services/` - Business logic and feature UI state models

**Adapters:**
- **Driving** (Primary): `internal/ui/*/handlers.go`, `routes.go`
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
    sessionName string
}

func NewTodoService(todoRepo domain.TodoRepository, sessionRepo domain.SessionRepository, store sessions.Store, sessionName string) *TodoService {
    return &TodoService{todoRepo: todoRepo, ...}
}
```

**Handlers** (HTTP Layer):
```go
type Handlers struct {
    logger      *slog.Logger
    todoService *services.TodoService
    sessionName string
}

func NewHandlers(logger *slog.Logger, svc *services.TodoService, nats *nats.Conn, sessionStore sessions.Store, sessionName string) *Handlers {
    return &Handlers{logger: logger, todoService: svc, ...}
}

// Handler signature
func (h *Handlers) IndexPage(w http.ResponseWriter, r *http.Request) {
    // Extract params → Get session → Call service → Respond
}
```

### Flat UI Organization
Organize shared UI and feature UI with a flat structure:

```
internal/ui/core/
├── component_header.go
├── component_sidebar.go
├── component_toast.go
└── layout_base.go

internal/ui/todo/
├── component_todo.go
├── page_index.go
├── handlers.go
└── routes.go

internal/services/
└── todo_service.go
```

### Parallel dev worktrees
`task dev` computes worktree-specific development settings so multiple local worktrees do not collide:
- `PORT`
- `NATS_PORT`
- `DB_PATH`
- `SESSION_NAME`
- `SESSION_SECRET`

Use the URL printed by `task dev` instead of assuming `http://localhost:8080`.

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
    sessionID, ok := RequireSession(h.sessionStore, h.sessionName, w, r)
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

### Gomponents Component Patterns

**Tailwind/DaisyUI classes are written inline in views:**
```go
Button(
    Class("btn btn-primary btn-lg"),
    Text("Save"),
)
```

**Datastar attributes (import as `data`):**
```go
import (
    . "maragu.dev/gomponents"
    . "maragu.dev/gomponents/html"
    data "maragu.dev/gomponents-datastar"
)

// Signals
data.Signals(map[string]any{"input": input})
data.Bind("input")

// Events
data.On("click", "@post('/api/todos/1/toggle')")
data.On("keydown", "if (evt.key === 'Enter') { ... }")

// HTTP methods
data.Init("@get('/api/todos/updates')")
data.On("click", "@post('/api/todos/1')")
data.On("click", "@put('/api/todos/reset')")
data.On("click", "@delete('/api/todos/1')")

// Conditional rendering
data.Show("$theme === 'light'")
```

**Component composition:**
```go
func Base(title string, children ...Node) Node {
    return Body(
        Group(children),
    )
}

// Usage
layouts.Base("My Page", Div(Text("Content here")))
```

## Development Workflow

### When modifying files:

**Gomponents view files** (`.go`):
- Author views directly in Go; there is no template code generation step
- Auto-rebuilt in dev mode via `task dev` and Air
- Shared rendering helpers live in `internal/platform/render`

`task dev` uses worktree-specific ports, NATS state, SQLite DB files, and cookie names so parallel worktrees can run side by side.

**SQL queries** (`internal/store/queries/*.sql`):
- Generate with `task generate:sqlc` or `sqlc generate`
- Creates type-safe Go code in `internal/store/queries/`

**Tailwind CSS entry file** (`web/ui/styles/main.css`):
- Keep this file minimal: Tailwind import, `@source`, and DaisyUI plugin config only
- Build CSS with `task build:web:css`
- Prefer DaisyUI defaults and Tailwind utility classes directly in gomponents/Lit markup

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
