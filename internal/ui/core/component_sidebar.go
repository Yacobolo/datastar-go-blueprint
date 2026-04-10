package core

import (
	g "maragu.dev/gomponents"
	data "maragu.dev/gomponents-datastar"
	h "maragu.dev/gomponents/html"
)

// Sidebar renders a simple collapsible navigation sidebar.
func Sidebar() g.Node {
	return h.Div(
		h.Class("drawer-side z-30"),
		h.Label(
			h.For("app-drawer"),
			h.Class("drawer-overlay"),
			data.On("click", "$sidebarOpen = false"),
		),
		h.Aside(
			h.Class("flex min-h-full w-72 flex-col overflow-x-hidden border-r border-base-300 bg-base-100 transition-[width] duration-200"),
			h.Data("class", "{'lg:w-20': $sidebarCollapsed, 'lg:w-72': !$sidebarCollapsed}"),
			h.Div(
				h.Class("flex items-center gap-3 border-b border-base-300 px-4 py-4"),
				h.A(
					h.Href("/"),
					h.Class("btn btn-ghost justify-start px-2 text-lg normal-case"),
					h.Span(h.Class("text-primary"), g.Text("DS")),
					h.Span(
						h.Class("truncate"),
						h.Data("class", "{'lg:hidden': $sidebarCollapsed}"),
						g.Text("Datastar"),
					),
				),
			),
			h.Ul(
				h.Class("menu flex-1 gap-1 p-4"),
				h.Li(
					h.A(
						h.Href("/"),
						h.Class("menu-active"),
						h.Data("class", "{'lg:justify-center': $sidebarCollapsed}"),
						IconHome(),
						h.Span(
							h.Class("truncate"),
							h.Data("class", "{'lg:hidden': $sidebarCollapsed}"),
							g.Text("Todos"),
						),
					),
				),
			),
			h.Div(
				h.Class("border-t border-base-300 p-4"),
				h.Button(
					h.Class("btn btn-ghost w-full justify-start"),
					h.Data("class", "{'lg:justify-center': $sidebarCollapsed}"),
					data.Attr("title", "$sidebarCollapsed ? 'Expand sidebar' : 'Collapse sidebar'"),
					data.On("click", "$sidebarCollapsed = !$sidebarCollapsed"),
					g.Attr("aria-label", "Toggle sidebar"),
					h.Type("button"),
					h.Span(
						data.Show("!$sidebarCollapsed"),
						IconChevronLeft(),
					),
					h.Span(
						data.Show("$sidebarCollapsed"),
						IconChevronRight(),
					),
					h.Span(
						h.Class("truncate"),
						h.Data("class", "{'lg:hidden': $sidebarCollapsed}"),
						g.Text("Collapse"),
					),
				),
			),
		),
	)
}

// MobileSidebar renders the mobile overlay sidebar.
func MobileSidebar() g.Node {
	return nil
}
