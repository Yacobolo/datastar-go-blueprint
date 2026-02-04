package todo

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	commoncomponents "github.com/yacobolo/datastar-go-blueprint/internal/features/common/components"
	todocomponents "github.com/yacobolo/datastar-go-blueprint/internal/features/todo/components"
	"github.com/yacobolo/datastar-go-blueprint/internal/features/todo/pages"
	"github.com/yacobolo/datastar-go-blueprint/internal/features/todo/services"
	"github.com/yacobolo/datastar-go-blueprint/internal/platform/pubsub"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/nats-io/nats.go"
	"github.com/starfederation/datastar-go/datastar"
)

// Helper functions
func RequireSession(store sessions.Store, w http.ResponseWriter, r *http.Request) (string, bool) {
	sess, err := store.Get(r, "connections")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", false
	}

	id, ok := sess.Values["id"].(string)
	if !ok {
		id = uuid.New().String()
		sess.Values["id"] = id
		if err := sess.Save(r, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return "", false
		}
	}
	return id, true
}

func RequireIntParam(w http.ResponseWriter, r *http.Request, param string) (int, bool) {
	val, err := strconv.Atoi(chi.URLParam(r, param))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 0, false
	}
	return val, true
}

func LogConsoleError(sse *datastar.ServerSentEventGenerator, err error) {
	if err := sse.ConsoleError(err); err != nil {
		slog.Error("failed to send console error", "error", err)
	}
}

type Handlers struct {
	todoService  *services.TodoService
	nats         *nats.Conn
	sessionStore sessions.Store
}

func NewHandlers(todoService *services.TodoService, nats *nats.Conn, sessionStore sessions.Store) *Handlers {
	return &Handlers{
		todoService:  todoService,
		nats:         nats,
		sessionStore: sessionStore,
	}
}

// subject returns the NATS subject for a session
func subject(sessionID string) string {
	return "todos.updates." + sessionID
}

// IndexPage renders the initial page
func (h *Handlers) IndexPage(w http.ResponseWriter, r *http.Request) {
	if err := pages.IndexPage("Datastar Go Blueprint").Render(r.Context(), w); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// TodosUpdates is the long-running SSE endpoint that pushes real-time updates
func (h *Handlers) TodosUpdates(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := RequireSession(h.sessionStore, w, r)
	if !ok {
		return
	}

	sse := datastar.NewSSE(w, r)
	ctx := r.Context()

	// Send initial state
	if err := h.refreshTodos(ctx, sse, sessionID); err != nil {
		LogConsoleError(sse, err)
		return
	}

	// Subscribe to NATS updates for this session
	msgChan := make(chan *nats.Msg, 64)
	sub, err := h.nats.ChanSubscribe(subject(sessionID), msgChan)
	if err != nil {
		LogConsoleError(sse, err)
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
					LogConsoleError(sse, err)
					return
				}
			}

			// Send toast if present
			if updateMsg.Toast != nil {
				toastComponent := commoncomponents.Toast(updateMsg.Toast.Message, updateMsg.Toast.Type)
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
	// Get MVC state from service
	mvc, err := h.todoService.GetMVCBySessionID(ctx, sessionID)
	if err != nil {
		return err
	}

	return sse.PatchElementTempl(todocomponents.TodosMVCView(mvc))
}

// notifyUpdate publishes a NATS message to trigger UI refresh
func (h *Handlers) notifyUpdate(sessionID string, opts ...pubsub.NotifyOption) {
	if err := pubsub.Notify(h.nats, subject(sessionID), opts...); err != nil {
		slog.Error("failed to notify update", "error", err)
	}
}

// ResetTodos resets to default todos
func (h *Handlers) ResetTodos(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := RequireSession(h.sessionStore, w, r)
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
		pubsub.WithToast("Todos reset", commoncomponents.ToastSuccess))

	w.WriteHeader(http.StatusOK)
}

// CancelEdit cancels editing mode
func (h *Handlers) CancelEdit(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := RequireSession(h.sessionStore, w, r)
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
	sessionID, ok := RequireSession(h.sessionStore, w, r)
	if !ok {
		return
	}

	modeStr := chi.URLParam(r, "mode")
	modeRaw, err := strconv.Atoi(modeStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mode := todocomponents.TodoViewMode(modeRaw)
	if mode < todocomponents.TodoViewModeAll || mode > todocomponents.TodoViewModeCompleted {
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
	sessionID, ok := RequireSession(h.sessionStore, w, r)
	if !ok {
		return
	}

	idx, ok := RequireIntParam(w, r, "idx")
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
	sessionID, ok := RequireSession(h.sessionStore, w, r)
	if !ok {
		return
	}

	idx, ok := RequireIntParam(w, r, "idx")
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

	sessionID, ok := RequireSession(h.sessionStore, w, r)
	if !ok {
		return
	}

	idx, ok := RequireIntParam(w, r, "idx")
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
		pubsub.WithToast(toastMsg, commoncomponents.ToastSuccess))

	w.WriteHeader(http.StatusOK)
}

// DeleteTodo removes a todo
func (h *Handlers) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := RequireSession(h.sessionStore, w, r)
	if !ok {
		return
	}

	idx, ok := RequireIntParam(w, r, "idx")
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
		pubsub.WithToast("Todo deleted", commoncomponents.ToastSuccess))

	w.WriteHeader(http.StatusOK)
}
