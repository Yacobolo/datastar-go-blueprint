// Package components provides reusable UI building blocks for the app shell and shared views.
package components

import (
	g "maragu.dev/gomponents"
	data "maragu.dev/gomponents-datastar"
	h "maragu.dev/gomponents/html"
)

// Header renders the app header with theme toggle.
func Header(title string) g.Node {
	return h.Header(
		h.Class("navbar sticky top-0 z-20 border-b border-base-300 bg-base-100/90 px-4 shadow-sm backdrop-blur"),
		h.Button(
			h.Class("btn btn-square btn-ghost lg:hidden"),
			data.On("click", "$sidebarOpen = !$sidebarOpen"),
			g.Attr("aria-label", "Open menu"),
			h.Type("button"),
			IconMenu(),
		),
		h.Div(
			h.Class("flex-1"),
			h.H1(
				h.Class("text-lg font-semibold"),
				g.Text(title),
			),
		),
		h.Button(
			h.Class("btn btn-circle btn-ghost"),
			data.On("click", "window.themeToggle.cycle()"),
			data.Attr("title", "$theme === 'system' ? 'Theme: system (' + $resolvedTheme + ')' : 'Theme: ' + $theme"),
			g.Attr("aria-label", "Toggle theme"),
			h.Type("button"),
			h.Span(
				data.Show("$resolvedTheme === 'light'"),
				IconSun(),
			),
			h.Span(
				data.Show("$resolvedTheme === 'dark'"),
				IconMoon(),
			),
		),
	)
}
