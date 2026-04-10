// Package services contains business logic for the todo feature.
package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/yacobolo/datastar-go-blueprint/internal/domain"
	"github.com/yacobolo/datastar-go-blueprint/internal/store/queries"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/samber/lo"
)

// TodoService provides business logic for managing todos.
type TodoService struct {
	todoRepo    domain.TodoRepository
	sessionRepo domain.SessionRepository
	sessionName string
	store       sessions.Store
}

// NewTodoService creates a new TodoService with the given repositories.
func NewTodoService(todoRepo domain.TodoRepository, sessionRepo domain.SessionRepository, store sessions.Store, sessionName string) *TodoService {
	return &TodoService{
		todoRepo:    todoRepo,
		sessionRepo: sessionRepo,
		sessionName: sessionName,
		store:       store,
	}
}

// GetSessionMVC retrieves the TodoMVC state for the current session.
func (s *TodoService) GetSessionMVC(w http.ResponseWriter, r *http.Request) (string, *TodoMVC, error) {
	ctx := r.Context()
	sessionID, err := s.upsertSessionID(r, w)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get session id: %w", err)
	}

	mvc, err := s.GetMVCBySessionID(ctx, sessionID)
	if err != nil {
		return "", nil, err
	}

	return sessionID, mvc, nil
}

// GetMVCBySessionID gets the TodoMVC state for a given session ID.
// This is used by SSE handlers that already have the session ID.
func (s *TodoService) GetMVCBySessionID(ctx context.Context, sessionID string) (*TodoMVC, error) {
	// Get todos from database
	dbTodos, err := s.todoRepo.GetTodosByUser(ctx, sessionID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to get todos: %w", err)
	}

	// Get session to load UI state
	session, err := s.sessionRepo.GetSession(ctx, sessionID)
	mode := TodoViewModeAll
	editingIdx := -1

	if err == nil {
		// Session exists, load UI state
		if session.Mode.Valid {
			mode = TodoViewMode(session.Mode.Int64)
		}
		if session.EditingIdx.Valid {
			editingIdx = int(session.EditingIdx.Int64)
		}
	}

	mvc := &TodoMVC{
		Mode:       mode,
		EditingIdx: editingIdx,
	}

	// Convert database todos to component todos
	if len(dbTodos) == 0 {
		// Initialize with default todos
		s.resetMVC(mvc)
		// Save defaults to database
		if err := s.saveMVCToDB(ctx, sessionID, mvc); err != nil {
			return nil, fmt.Errorf("failed to save default todos: %w", err)
		}
	} else {
		mvc.Todos = make([]*Todo, len(dbTodos))
		for i, dbTodo := range dbTodos {
			mvc.Todos[i] = &Todo{
				Text:      dbTodo.Task,
				Completed: dbTodo.Completed.Int64 == 1,
			}
		}
	}

	return mvc, nil
}

// SaveMVC persists the TodoMVC state to the database.
func (s *TodoService) SaveMVC(ctx context.Context, sessionID string, mvc *TodoMVC) error {
	return s.saveMVCToDB(ctx, sessionID, mvc)
}

// ResetMVC resets the TodoMVC to its initial state.
func (s *TodoService) ResetMVC(mvc *TodoMVC) {
	s.resetMVC(mvc)
}

// ToggleTodo toggles the completion state of a todo by index.
func (s *TodoService) ToggleTodo(mvc *TodoMVC, index int) {
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

// EditTodo updates or creates a todo with the given text.
func (s *TodoService) EditTodo(mvc *TodoMVC, index int, text string) {
	if index >= 0 && index < len(mvc.Todos) {
		mvc.Todos[index].Text = text
	} else if index < 0 {
		mvc.Todos = append(mvc.Todos, &Todo{
			Text:      text,
			Completed: false,
		})
	}
	mvc.EditingIdx = -1
}

// DeleteTodo removes a todo by index or clears completed todos if index is -1.
func (s *TodoService) DeleteTodo(mvc *TodoMVC, index int) {
	if index >= 0 && index < len(mvc.Todos) {
		mvc.Todos = append(mvc.Todos[:index], mvc.Todos[index+1:]...)
	} else if index < 0 {
		mvc.Todos = lo.Filter(mvc.Todos, func(item *Todo, _ int) bool {
			return !item.Completed
		})
	}
}

// SetMode changes the view filter mode for todos.
func (s *TodoService) SetMode(mvc *TodoMVC, mode TodoViewMode) {
	mvc.Mode = mode
}

// StartEditing puts a todo into edit mode.
func (s *TodoService) StartEditing(mvc *TodoMVC, index int) {
	mvc.EditingIdx = index
}

// CancelEditing exits edit mode without saving.
func (s *TodoService) CancelEditing(mvc *TodoMVC) {
	mvc.EditingIdx = -1
}

func (s *TodoService) saveMVCToDB(ctx context.Context, sessionID string, mvc *TodoMVC) error {
	// Delete all existing todos for this user
	if err := s.todoRepo.DeleteAllTodosByUser(ctx, sessionID); err != nil {
		return fmt.Errorf("failed to delete existing todos: %w", err)
	}

	// Insert all todos
	for _, todo := range mvc.Todos {
		completed := int64(0)
		if todo.Completed {
			completed = 1
		}

		todoID := uuid.New().String()
		if err := s.todoRepo.CreateTodo(ctx, queries.CreateTodoParams{
			ID:        todoID,
			UserID:    sessionID,
			Task:      todo.Text,
			Completed: sql.NullInt64{Int64: completed, Valid: true},
		}); err != nil {
			return fmt.Errorf("failed to create todo: %w", err)
		}
	}

	// Save UI state to session
	if err := s.sessionRepo.UpsertSession(ctx, queries.UpsertSessionParams{
		ID:         sessionID,
		Data:       "",
		Mode:       sql.NullInt64{Int64: int64(mvc.Mode), Valid: true},
		EditingIdx: sql.NullInt64{Int64: int64(mvc.EditingIdx), Valid: true},
	}); err != nil {
		return fmt.Errorf("failed to save session state: %w", err)
	}

	return nil
}

func (s *TodoService) resetMVC(mvc *TodoMVC) {
	mvc.Mode = TodoViewModeAll
	mvc.Todos = []*Todo{
		{Text: "Learn any backend language", Completed: true},
		{Text: "Learn Datastar", Completed: false},
		{Text: "Create Hypermedia", Completed: false},
		{Text: "???", Completed: false},
		{Text: "Profit", Completed: false},
	}
	mvc.EditingIdx = -1
}

func (s *TodoService) upsertSessionID(r *http.Request, w http.ResponseWriter) (string, error) {
	sess, err := s.store.Get(r, s.sessionName)
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
