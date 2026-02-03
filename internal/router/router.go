package router

import (
	"context"
	"net/http"
	"sync"

	"github.com/yourusername/datastar-go-starter-kit/internal/app/handlers"
	"github.com/yourusername/datastar-go-starter-kit/internal/app/services"
	"github.com/yourusername/datastar-go-starter-kit/internal/config"
	"github.com/yourusername/datastar-go-starter-kit/internal/store/queries"
	"github.com/yourusername/datastar-go-starter-kit/web/resources"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/nats-io/nats.go"
	"github.com/starfederation/datastar-go/datastar"
)

func SetupRoutes(ctx context.Context, router chi.Router, sessionStore *sessions.CookieStore, q *queries.Queries, nc *nats.Conn) (err error) {

	if config.Global.Environment == config.Dev {
		setupReload(router)
	}

	router.Handle("/static/*", resources.Handler())

	// Create TODO service and handlers
	todoService, err := services.NewTodoService(q, sessionStore)
	if err != nil {
		return err
	}
	todoHandlers := handlers.NewHandlers(todoService, nc, sessionStore)

	// Setup index/TODO routes
	router.Get("/", todoHandlers.IndexPage)
	router.Route("/api/todos", func(r chi.Router) {
		r.Get("/updates", todoHandlers.TodosUpdates)
		r.Put("/reset", todoHandlers.ResetTodos)
		r.Put("/cancel", todoHandlers.CancelEdit)
		r.Put("/mode/{mode}", todoHandlers.SetMode)
		r.Post("/{idx}/toggle", todoHandlers.ToggleTodo)
		r.Post("/{idx}/start-edit", todoHandlers.StartEdit)
		r.Post("/{idx}/save-edit", todoHandlers.SaveEdit)
		r.Delete("/{idx}", todoHandlers.DeleteTodo)
	})

	return nil
}

func setupReload(router chi.Router) {
	reloadChan := make(chan struct{}, 1)
	var hotReloadOnce sync.Once

	router.Get("/reload", func(w http.ResponseWriter, r *http.Request) {
		sse := datastar.NewSSE(w, r)
		reload := func() { sse.ExecuteScript("window.location.reload()") }
		hotReloadOnce.Do(reload)
		select {
		case <-reloadChan:
			reload()
		case <-r.Context().Done():
		}
	})

	router.Get("/hotreload", func(w http.ResponseWriter, r *http.Request) {
		select {
		case reloadChan <- struct{}{}:
		default:
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

}
