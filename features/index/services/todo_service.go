package services

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	dbstore "github.com/yourusername/datastar-go-starter-kit/db"
	"github.com/yourusername/datastar-go-starter-kit/features/index/components"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/samber/lo"
)

type TodoService struct {
	queries *dbstore.Queries
	store   sessions.Store
}

func NewTodoService(queries *dbstore.Queries, store sessions.Store) (*TodoService, error) {
	return &TodoService{
		queries: queries,
		store:   store,
	}, nil
}

// Queries returns the sqlc queries instance for direct access
func (s *TodoService) Queries() *dbstore.Queries {
	return s.queries
}

func (s *TodoService) GetSessionMVC(w http.ResponseWriter, r *http.Request) (string, *components.TodoMVC, error) {
	ctx := r.Context()
	sessionID, err := s.upsertSessionID(r, w)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get session id: %w", err)
	}

	// Get todos from database
	dbTodos, err := s.queries.GetTodosByUser(ctx, sessionID)
	if err != nil && err != sql.ErrNoRows {
		return "", nil, fmt.Errorf("failed to get todos: %w", err)
	}

	mvc := &components.TodoMVC{
		Mode:       components.TodoViewModeAll,
		EditingIdx: -1,
	}

	// Convert database todos to component todos
	if len(dbTodos) == 0 {
		// Initialize with default todos
		s.resetMVC(mvc)
		// Save defaults to database
		if err := s.saveMVCToDB(ctx, sessionID, mvc); err != nil {
			return "", nil, fmt.Errorf("failed to save default todos: %w", err)
		}
	} else {
		mvc.Todos = make([]*components.Todo, len(dbTodos))
		for i, dbTodo := range dbTodos {
			mvc.Todos[i] = &components.Todo{
				Text:      dbTodo.Task,
				Completed: dbTodo.Completed.Int64 == 1,
			}
		}
	}

	return sessionID, mvc, nil
}

func (s *TodoService) SaveMVC(ctx context.Context, sessionID string, mvc *components.TodoMVC) error {
	return s.saveMVCToDB(ctx, sessionID, mvc)
}

func (s *TodoService) ResetMVC(mvc *components.TodoMVC) {
	s.resetMVC(mvc)
}

func (s *TodoService) ToggleTodo(mvc *components.TodoMVC, index int) {
	if index < 0 {
		setCompletedTo := false
		for _, todo := range mvc.Todos {
			if !todo.Completed {
				setCompletedTo = true
				break
			}
		}
		for _, todo := range mvc.Todos {
			todo.Completed = setCompletedTo
		}
	} else if index < len(mvc.Todos) {
		todo := mvc.Todos[index]
		todo.Completed = !todo.Completed
	}
}

func (s *TodoService) EditTodo(mvc *components.TodoMVC, index int, text string) {
	if index >= 0 && index < len(mvc.Todos) {
		mvc.Todos[index].Text = text
	} else if index < 0 {
		mvc.Todos = append(mvc.Todos, &components.Todo{
			Text:      text,
			Completed: false,
		})
	}
	mvc.EditingIdx = -1
}

func (s *TodoService) DeleteTodo(mvc *components.TodoMVC, index int) {
	if index >= 0 && index < len(mvc.Todos) {
		mvc.Todos = append(mvc.Todos[:index], mvc.Todos[index+1:]...)
	} else if index < 0 {
		mvc.Todos = lo.Filter(mvc.Todos, func(todo *components.Todo, i int) bool {
			return !todo.Completed
		})
	}
}

func (s *TodoService) SetMode(mvc *components.TodoMVC, mode components.TodoViewMode) {
	mvc.Mode = mode
}

func (s *TodoService) StartEditing(mvc *components.TodoMVC, index int) {
	mvc.EditingIdx = index
}

func (s *TodoService) CancelEditing(mvc *components.TodoMVC) {
	mvc.EditingIdx = -1
}

func (s *TodoService) saveMVCToDB(ctx context.Context, sessionID string, mvc *components.TodoMVC) error {
	// Delete all existing todos for this user
	if err := s.queries.DeleteAllTodosByUser(ctx, sessionID); err != nil {
		return fmt.Errorf("failed to delete existing todos: %w", err)
	}

	// Insert all todos
	for _, todo := range mvc.Todos {
		completed := int64(0)
		if todo.Completed {
			completed = 1
		}

		todoID := uuid.New().String()
		if err := s.queries.CreateTodo(ctx, dbstore.CreateTodoParams{
			ID:        todoID,
			UserID:    sessionID,
			Task:      todo.Text,
			Completed: sql.NullInt64{Int64: completed, Valid: true},
		}); err != nil {
			return fmt.Errorf("failed to create todo: %w", err)
		}
	}

	return nil
}

func (s *TodoService) resetMVC(mvc *components.TodoMVC) {
	mvc.Mode = components.TodoViewModeAll
	mvc.Todos = []*components.Todo{
		{Text: "Learn any backend language", Completed: true},
		{Text: "Learn Datastar", Completed: false},
		{Text: "Create Hypermedia", Completed: false},
		{Text: "???", Completed: false},
		{Text: "Profit", Completed: false},
	}
	mvc.EditingIdx = -1
}

func (s *TodoService) upsertSessionID(r *http.Request, w http.ResponseWriter) (string, error) {
	sess, err := s.store.Get(r, "connections")
	if err != nil {
		return "", fmt.Errorf("failed to get session: %w", err)
	}

	id, ok := sess.Values["id"].(string)

	if !ok {
		id = uuid.New().String()
		sess.Values["id"] = id
		if err := sess.Save(r, w); err != nil {
			return "", fmt.Errorf("failed to save session: %w", err)
		}
	}

	return id, nil
}
