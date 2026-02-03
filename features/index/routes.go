package index

import (
	dbstore "github.com/yourusername/datastar-go-starter-kit/db"
	"github.com/yourusername/datastar-go-starter-kit/features/index/services"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/nats-io/nats.go"
)

func SetupRoutes(router chi.Router, sessStore sessions.Store, queries *dbstore.Queries, nc *nats.Conn) error {
	todoService, err := services.NewTodoService(queries, sessStore)
	if err != nil {
		return err
	}

	handlers := NewHandlers(todoService, nc, sessStore)

	router.Get("/", handlers.IndexPage)

	router.Route("/api", func(apiRouter chi.Router) {
		apiRouter.Route("/todos", func(todosRouter chi.Router) {
			todosRouter.Get("/updates", handlers.TodosUpdates) // Changed from "/" to "/updates"
			todosRouter.Put("/reset", handlers.ResetTodos)
			todosRouter.Put("/cancel", handlers.CancelEdit)
			todosRouter.Put("/mode/{mode}", handlers.SetMode)

			todosRouter.Route("/{idx}", func(todoRouter chi.Router) {
				todoRouter.Post("/toggle", handlers.ToggleTodo)
				todoRouter.Route("/edit", func(editRouter chi.Router) {
					editRouter.Get("/", handlers.StartEdit)
					editRouter.Put("/", handlers.SaveEdit)
				})
				todoRouter.Delete("/", handlers.DeleteTodo)
			})
		})
	})

	return nil
}
