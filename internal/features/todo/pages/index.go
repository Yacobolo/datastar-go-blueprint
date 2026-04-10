// Package pages contains full-page gomponents views for the todo feature.
package pages

import (
	"github.com/yacobolo/datastar-go-blueprint/internal/features/common/layouts"
	appds "github.com/yacobolo/datastar-go-blueprint/internal/platform/ds"

	g "maragu.dev/gomponents"
	data "maragu.dev/gomponents-datastar"
	h "maragu.dev/gomponents/html"
)

// IndexPage renders the todo app shell.
func IndexPage(title string) g.Node {
	return layouts.Base(
		title,
		h.Div(
			h.Class("w-full"),
			h.Div(
				h.ID("todos-container"),
				data.Init(appds.Get("/api/todos/updates", appds.Opt("requestCancellation", "disabled"))),
				h.Div(
					h.Class("flex min-h-[16rem] items-center justify-center"),
					h.Span(h.Class("loading loading-spinner loading-lg text-primary")),
				),
			),
			h.Div(
				h.ID("toast-container"),
				h.Class("toast toast-top toast-end z-50"),
			),
		),
	)
}
