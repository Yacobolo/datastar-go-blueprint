// Package pages contains full-page gomponents views for the todo feature.
package pages

import (
	"github.com/yacobolo/datastar-go-blueprint/internal/features/common/layouts"
	appds "github.com/yacobolo/datastar-go-blueprint/internal/platform/ds"
	"github.com/yacobolo/datastar-go-blueprint/internal/ui"

	g "maragu.dev/gomponents"
	data "maragu.dev/gomponents-datastar"
	h "maragu.dev/gomponents/html"
)

// IndexPage renders the todo app shell.
func IndexPage(title string) g.Node {
	return layouts.Base(
		title,
		h.Div(
			h.Class(ui.Page),
			h.Div(
				h.ID("todos-container"),
				data.Init(appds.Get("/api/todos/updates", appds.Opt("requestCancellation", "disabled"))),
				h.Div(
					h.Class(ui.TodoLoading),
					h.P(g.Text("Loading todos...")),
				),
			),
			h.Div(
				h.ID("toast-container"),
				h.Class(ui.ToastContainer),
			),
		),
	)
}
