package router

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/yourusername/datastar-go-starter-kit/config"
	dbstore "github.com/yourusername/datastar-go-starter-kit/db"
	indexFeature "github.com/yourusername/datastar-go-starter-kit/features/index"
	"github.com/yourusername/datastar-go-starter-kit/web/resources"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/nats-io/nats.go"
	"github.com/starfederation/datastar-go/datastar"
)

func SetupRoutes(ctx context.Context, router chi.Router, sessionStore *sessions.CookieStore, queries *dbstore.Queries, nc *nats.Conn) (err error) {

	if config.Global.Environment == config.Dev {
		setupReload(router)
	}

	router.Handle("/static/*", resources.Handler())

	if err := indexFeature.SetupRoutes(router, sessionStore, queries, nc); err != nil {
		return fmt.Errorf("error setting up routes: %w", err)
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
