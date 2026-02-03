package index

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/yourusername/datastar-go-starter-kit/features/common/components"
	"github.com/yourusername/datastar-go-starter-kit/features/common/handlers"
	indexcomponents "github.com/yourusername/datastar-go-starter-kit/features/index/components"
	"github.com/yourusername/datastar-go-starter-kit/features/index/pages"
	"github.com/yourusername/datastar-go-starter-kit/features/index/services"
	"github.com/yourusername/datastar-go-starter-kit/internal/pubsub"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/nats-io/nats.go"
	"github.com/starfederation/datastar-go/datastar"
)

type Handlers struct {
	todoService  *services.TodoService
	nc           *nats.Conn
	sessionStore sessions.Store
}

func NewHandlers(todoService *services.TodoService, nc *nats.Conn, store sessions.Store) *Handlers {
	return &Handlers{
		todoService:  todoService,
		nc:           nc,
		sessionStore: store,
	}
}

// subject returns the NATS subject for a session
func subject(sessionID string) string {
	return "todos.updates." + sessionID
}

// IndexPage renders the initial page
func (h *Handlers) IndexPage(w http.ResponseWriter, r *http.Request) {
	if err := pages.IndexPage("Datastar Go Starter Kit").Render(r.Context(), w); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// TodosUpdates is the long-running SSE endpoint that pushes real-time updates
func (h *Handlers) TodosUpdates(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := handlers.RequireSession(h.sessionStore, w, r)
	if !ok {
		return
	}

	sse := datastar.NewSSE(w, r)
	ctx := r.Context()

	// Send initial state
	if err := h.refreshTodos(ctx, sse, sessionID); err != nil {
		handlers.LogConsoleError(sse, err)
		return
	}

	// Subscribe to NATS updates for this session
	msgChan := make(chan *nats.Msg, 64)
	sub, err := h.nc.ChanSubscribe(subject(sessionID), msgChan)
	if err != nil {
		handlers.LogConsoleError(sse, err)
		return
	}
	defer sub.Unsubscribe()

	// Listen for updates
	for {
		select {
		case <-ctx.Done():
			return
		case natsMsg := <-msgChan:
			updateMsg, err := pubsub.ParseUpdateMessage(natsMsg.Data)
			if err != nil {
				slog.Error("failed to parse update message", "error", err)
				continue
			}

			// Refresh TODO list if requested
			if updateMsg.RefreshTodos {
				if err := h.refreshTodos(ctx, sse, sessionID); err != nil {
					handlers.LogConsoleError(sse, err)
					return
				}
			}

			// Send toast if present
			if updateMsg.Toast != nil {
				toastComponent := components.Toast(updateMsg.Toast.Message, updateMsg.Toast.Type)
				if err := sse.PatchElementTempl(
					toastComponent,
					datastar.WithSelectorID("toast-container"),
					datastar.WithModeAppend(),
				); err != nil {
					slog.Error("failed to send toast", "error", err)
				}
			}
		}
	}
}

// refreshTodos fetches current state and sends via SSE
func (h *Handlers) refreshTodos(ctx context.Context, sse *datastar.ServerSentEventGenerator, sessionID string) error {
	// Fetch todos from database
	dbTodos, err := h.todoService.Queries().GetTodosByUser(ctx, sessionID)
	if err != nil {
		return err
	}

	mvc := &indexcomponents.TodoMVC{
		Mode:       indexcomponents.TodoViewModeAll,
		EditingIdx: -1,
	}

	if len(dbTodos) == 0 {
		h.todoService.ResetMVC(mvc)
	} else {
		mvc.Todos = make([]*indexcomponents.Todo, len(dbTodos))
		for i, dbTodo := range dbTodos {
			mvc.Todos[i] = &indexcomponents.Todo{
				Text:      dbTodo.Task,
				Completed: dbTodo.Completed.Int64 == 1,
			}
		}
	}

	return sse.PatchElementTempl(indexcomponents.TodosMVCView(mvc))
}

// notifyUpdate publishes a NATS message to trigger UI refresh
func (h *Handlers) notifyUpdate(sessionID string, opts ...pubsub.NotifyOption) {
	if err := pubsub.Notify(h.nc, subject(sessionID), opts...); err != nil {
		slog.Error("failed to notify update", "error", err)
	}
}

// ResetTodos resets to default todos
func (h *Handlers) ResetTodos(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := handlers.RequireSession(h.sessionStore, w, r)
	if !ok {
		return
	}

	_, mvc, err := h.todoService.GetSessionMVC(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.todoService.ResetMVC(mvc)
	if err := h.todoService.SaveMVC(r.Context(), sessionID, mvc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Notify via NATS (triggers SSE push)
	h.notifyUpdate(sessionID,
		pubsub.WithRefresh(),
		pubsub.WithToast("Todos reset", components.ToastSuccess))

	w.WriteHeader(http.StatusOK)
}

// CancelEdit cancels editing mode
func (h *Handlers) CancelEdit(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := handlers.RequireSession(h.sessionStore, w, r)
	if !ok {
		return
	}

	_, mvc, err := h.todoService.GetSessionMVC(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.todoService.CancelEditing(mvc)
	if err := h.todoService.SaveMVC(r.Context(), sessionID, mvc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.notifyUpdate(sessionID, pubsub.WithRefresh())
	w.WriteHeader(http.StatusOK)
}

// SetMode changes the view filter mode
func (h *Handlers) SetMode(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := handlers.RequireSession(h.sessionStore, w, r)
	if !ok {
		return
	}

	modeStr := chi.URLParam(r, "mode")
	modeRaw, err := strconv.Atoi(modeStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mode := indexcomponents.TodoViewMode(modeRaw)
	if mode < indexcomponents.TodoViewModeAll || mode > indexcomponents.TodoViewModeCompleted {
		http.Error(w, "invalid mode", http.StatusBadRequest)
		return
	}

	_, mvc, err := h.todoService.GetSessionMVC(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.todoService.SetMode(mvc, mode)
	if err := h.todoService.SaveMVC(r.Context(), sessionID, mvc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.notifyUpdate(sessionID, pubsub.WithRefresh())
	w.WriteHeader(http.StatusOK)
}

// ToggleTodo toggles completion state
func (h *Handlers) ToggleTodo(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := handlers.RequireSession(h.sessionStore, w, r)
	if !ok {
		return
	}

	idx, ok := handlers.RequireIntParam(w, r, "idx")
	if !ok {
		return
	}

	_, mvc, err := h.todoService.GetSessionMVC(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.todoService.ToggleTodo(mvc, idx)
	if err := h.todoService.SaveMVC(r.Context(), sessionID, mvc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.notifyUpdate(sessionID, pubsub.WithRefresh())
	w.WriteHeader(http.StatusOK)
}

// StartEdit enters edit mode for a todo
func (h *Handlers) StartEdit(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := handlers.RequireSession(h.sessionStore, w, r)
	if !ok {
		return
	}

	idx, ok := handlers.RequireIntParam(w, r, "idx")
	if !ok {
		return
	}

	_, mvc, err := h.todoService.GetSessionMVC(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.todoService.StartEditing(mvc, idx)
	if err := h.todoService.SaveMVC(r.Context(), sessionID, mvc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.notifyUpdate(sessionID, pubsub.WithRefresh())
	w.WriteHeader(http.StatusOK)
}

// SaveEdit creates or updates a todo
func (h *Handlers) SaveEdit(w http.ResponseWriter, r *http.Request) {
	type Store struct {
		Input string `json:"input"`
	}
	store := &Store{}

	if err := datastar.ReadSignals(r, store); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if store.Input == "" {
		return
	}

	sessionID, ok := handlers.RequireSession(h.sessionStore, w, r)
	if !ok {
		return
	}

	idx, ok := handlers.RequireIntParam(w, r, "idx")
	if !ok {
		return
	}

	_, mvc, err := h.todoService.GetSessionMVC(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.todoService.EditTodo(mvc, idx, store.Input)
	if err := h.todoService.SaveMVC(r.Context(), sessionID, mvc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Notify via NATS
	toastMsg := "Todo updated"
	if idx < 0 {
		toastMsg = "Todo created"
	}
	h.notifyUpdate(sessionID,
		pubsub.WithRefresh(),
		pubsub.WithToast(toastMsg, components.ToastSuccess))

	w.WriteHeader(http.StatusOK)
}

// DeleteTodo removes a todo
func (h *Handlers) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := handlers.RequireSession(h.sessionStore, w, r)
	if !ok {
		return
	}

	idx, ok := handlers.RequireIntParam(w, r, "idx")
	if !ok {
		return
	}

	_, mvc, err := h.todoService.GetSessionMVC(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.todoService.DeleteTodo(mvc, idx)
	if err := h.todoService.SaveMVC(r.Context(), sessionID, mvc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.notifyUpdate(sessionID,
		pubsub.WithRefresh(),
		pubsub.WithToast("Todo deleted", components.ToastSuccess))

	w.WriteHeader(http.StatusOK)
}
