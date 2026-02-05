// Package router configures HTTP routes and middleware.
package router

import (
	"context"
	"net/http"
	"sync"

	"github.com/yacobolo/datastar-go-blueprint/internal/app"
	"github.com/yacobolo/datastar-go-blueprint/internal/config"
	"github.com/yacobolo/datastar-go-blueprint/internal/features/todo"
	"github.com/yacobolo/datastar-go-blueprint/web/resources"

	"github.com/go-chi/chi/v5"
	"github.com/starfederation/datastar-go/datastar"
)

// SetupRoutes configures all HTTP routes for the application.
func SetupRoutes(_ context.Context, router chi.Router, application *app.App) error {

	if config.Global.Environment == config.Dev {
		setupReload(router)
	}

	router.Handle("/static/*", resources.Handler())

	// Setup feature routes
	if err := todo.SetupRoutes(router, application); err != nil {
		return err
	}

	return nil
}

func setupReload(router chi.Router) {
	reloadChan := make(chan struct{}, 1)
	var hotReloadOnce sync.Once

	router.Get("/reload", func(w http.ResponseWriter, r *http.Request) {
		sse := datastar.NewSSE(w, r)
		reload := func() { _ = sse.ExecuteScript("window.location.reload()") }
		hotReloadOnce.Do(reload)
		select {
		case <-reloadChan:
			reload()
		case <-r.Context().Done():
		}
	})

	router.Get("/hotreload", func(w http.ResponseWriter, _ *http.Request) {
		select {
		case reloadChan <- struct{}{}:
		default:
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

}
