package components

import (
	"github.com/yacobolo/datastar-go-blueprint/internal/ui"

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
		h.Class(classNames(ui.Flex, ui.JustifyCenter, ui.PMd)),
		h.Ul(
			h.Class(classNames(ui.Flex, ui.GapMd)),
			h.Li(
				h.Class(classNames(
					ui.HoverTextPrimary,
					when(current == pageIndex, ui.TextPrimary),
					when(current == pageIndex, ui.FontBold),
				)),
				h.A(
					h.Href("/"),
					g.Text("TODO App"),
				),
			),
		),
	)
}
