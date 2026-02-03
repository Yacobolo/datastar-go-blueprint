package router

import (
	"context"
	"net/http"
	"sync"

	"github.com/yacobolo/datastar-go-starter-kit/internal/config"
	"github.com/yacobolo/datastar-go-starter-kit/internal/features/todo"
	"github.com/yacobolo/datastar-go-starter-kit/internal/features/todo/services"
	"github.com/yacobolo/datastar-go-starter-kit/internal/store/queries"
	"github.com/yacobolo/datastar-go-starter-kit/web/resources"

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

	// Create TODO service
	todoService, err := services.NewTodoService(q, sessionStore)
	if err != nil {
		return err
	}

	// Setup feature routes
	if err := todo.SetupRoutes(router, sessionStore, nc, todoService); err != nil {
		return err
	}

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
