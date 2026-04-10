// Package components provides reusable UI building blocks for the app shell and shared views.
package components

import (
	"github.com/yacobolo/datastar-go-blueprint/internal/ui"

	g "maragu.dev/gomponents"
	data "maragu.dev/gomponents-datastar"
	h "maragu.dev/gomponents/html"
)

// Header renders the app header with theme toggle.
func Header(title string) g.Node {
	return h.Header(
		h.Class(ui.AppHeader),
		h.Button(
			h.Class(ui.HeaderMenuToggle),
			data.On("click", "$sidebarOpen = !$sidebarOpen"),
			g.Attr("aria-label", "Open menu"),
			h.Type("button"),
			IconMenu(),
		),
		h.Div(
			h.Class(ui.HeaderTitle),
			h.H1(g.Text(title)),
		),
		h.Div(h.Class(ui.HeaderSpacer)),
		h.Button(
			h.Class(ui.HeaderIconBtn),
			data.On("click", "$theme = window.themeToggle.cycle()"),
			data.Attr("title", "'Theme: ' + $theme"),
			g.Attr("aria-label", "Toggle theme"),
			h.Type("button"),
			h.Span(
				data.Show("$theme === 'light'"),
				IconSun(),
			),
			h.Span(
				data.Show("$theme === 'dark'"),
				IconMoon(),
			),
			h.Span(
				data.Show("$theme === 'system'"),
				IconMonitor(),
			),
		),
	)
}
