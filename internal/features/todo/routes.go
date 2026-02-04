package todo

import (
	"github.com/yacobolo/datastar-go-blueprint/internal/app"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(router chi.Router, application *app.App) error {
	// Extract specific dependencies from App and pass to handlers
	handlers := NewHandlers(
		application.Services.Todo,
		application.NATS,
		application.SessionStore,
	)

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
