package services

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/yacobolo/datastar-go-starter-kit/internal/domain"
	todocomponents "github.com/yacobolo/datastar-go-starter-kit/internal/features/todo/components"
	"github.com/yacobolo/datastar-go-starter-kit/internal/store/queries"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/samber/lo"
)

type TodoService struct {
	todoRepo    domain.TodoRepository
	sessionRepo domain.SessionRepository
	store       sessions.Store
}

func NewTodoService(todoRepo domain.TodoRepository, sessionRepo domain.SessionRepository, store sessions.Store) *TodoService {
	return &TodoService{
		todoRepo:    todoRepo,
		sessionRepo: sessionRepo,
		store:       store,
	}
}

func (s *TodoService) GetSessionMVC(w http.ResponseWriter, r *http.Request) (string, *todocomponents.TodoMVC, error) {
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
func (s *TodoService) GetMVCBySessionID(ctx context.Context, sessionID string) (*todocomponents.TodoMVC, error) {
	// Get todos from database
	dbTodos, err := s.todoRepo.GetTodosByUser(ctx, sessionID)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get todos: %w", err)
	}

	// Get session to load UI state
	session, err := s.sessionRepo.GetSession(ctx, sessionID)
	mode := todocomponents.TodoViewModeAll
	editingIdx := -1

	if err == nil {
		// Session exists, load UI state
		if session.Mode.Valid {
			mode = todocomponents.TodoViewMode(session.Mode.Int64)
		}
		if session.EditingIdx.Valid {
			editingIdx = int(session.EditingIdx.Int64)
		}
	}

	mvc := &todocomponents.TodoMVC{
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
		mvc.Todos = make([]*todocomponents.Todo, len(dbTodos))
		for i, dbTodo := range dbTodos {
			mvc.Todos[i] = &todocomponents.Todo{
				Text:      dbTodo.Task,
				Completed: dbTodo.Completed.Int64 == 1,
			}
		}
	}

	return mvc, nil
}

func (s *TodoService) SaveMVC(ctx context.Context, sessionID string, mvc *todocomponents.TodoMVC) error {
	return s.saveMVCToDB(ctx, sessionID, mvc)
}

func (s *TodoService) ResetMVC(mvc *todocomponents.TodoMVC) {
	s.resetMVC(mvc)
}

func (s *TodoService) ToggleTodo(mvc *todocomponents.TodoMVC, index int) {
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

func (s *TodoService) EditTodo(mvc *todocomponents.TodoMVC, index int, text string) {
	if index >= 0 && index < len(mvc.Todos) {
		mvc.Todos[index].Text = text
	} else if index < 0 {
		mvc.Todos = append(mvc.Todos, &todocomponents.Todo{
			Text:      text,
			Completed: false,
		})
	}
	mvc.EditingIdx = -1
}

func (s *TodoService) DeleteTodo(mvc *todocomponents.TodoMVC, index int) {
	if index >= 0 && index < len(mvc.Todos) {
		mvc.Todos = append(mvc.Todos[:index], mvc.Todos[index+1:]...)
	} else if index < 0 {
		mvc.Todos = lo.Filter(mvc.Todos, func(todo *todocomponents.Todo, i int) bool {
			return !todo.Completed
		})
	}
}

func (s *TodoService) SetMode(mvc *todocomponents.TodoMVC, mode todocomponents.TodoViewMode) {
	mvc.Mode = mode
}

func (s *TodoService) StartEditing(mvc *todocomponents.TodoMVC, index int) {
	mvc.EditingIdx = index
}

func (s *TodoService) CancelEditing(mvc *todocomponents.TodoMVC) {
	mvc.EditingIdx = -1
}

func (s *TodoService) saveMVCToDB(ctx context.Context, sessionID string, mvc *todocomponents.TodoMVC) error {
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

func (s *TodoService) resetMVC(mvc *todocomponents.TodoMVC) {
	mvc.Mode = todocomponents.TodoViewModeAll
	mvc.Todos = []*todocomponents.Todo{
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
