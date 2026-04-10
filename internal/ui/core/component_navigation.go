package core

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type page int

const (
	pageIndex page = iota
)

// Navigation renders the shared top navigation.
func Navigation(current page) g.Node {
	return h.Nav(
		h.Class("flex justify-center p-4"),
		h.Ul(
			h.Class("menu menu-horizontal rounded-box bg-base-100 shadow"),
			h.Li(
				h.Class(classNames(when(current == pageIndex, "menu-active"))),
				h.A(
					h.Href("/"),
					g.Text("TODO App"),
				),
			),
		),
	)
}
