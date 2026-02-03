package todo

import (
	"github.com/yacobolo/datastar-go-starter-kit/internal/features/todo/services"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/nats-io/nats.go"
)

func SetupRoutes(router chi.Router, store sessions.Store, nc *nats.Conn, todoService *services.TodoService) error {
	handlers := NewHandlers(todoService, nc, store)

	router.Get("/", handlers.IndexPage)

	router.Route("/api", func(apiRouter chi.Router) {
		apiRouter.Route("/todos", func(todosRouter chi.Router) {
			todosRouter.Get("/updates", handlers.TodosUpdates)
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
