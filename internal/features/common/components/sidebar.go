package components

import (
	"github.com/yacobolo/datastar-go-blueprint/internal/ui"

	g "maragu.dev/gomponents"
	data "maragu.dev/gomponents-datastar"
	h "maragu.dev/gomponents/html"
)

// Sidebar renders a simple collapsible navigation sidebar.
func Sidebar() g.Node {
	return h.Aside(
		h.Class(ui.AppSidebar),
		h.Data("class", "{'app-sidebar--collapsed': $sidebarCollapsed}"),
		h.Div(
			h.Class(ui.SidebarHeader),
			h.A(
				h.Href("/"),
				h.Class(ui.SidebarLogo),
				h.Span(h.Class(ui.SidebarLogoIcon), g.Text("DS")),
				h.Span(h.Class(ui.SidebarLogoText), g.Text("Datastar")),
			),
		),
		h.Nav(
			h.Class(ui.SidebarNav),
			h.Div(
				h.Class(ui.NavSection),
				h.A(
					h.Href("/"),
					h.Class(classNames(ui.NavItem, ui.NavItemActive)),
					h.Span(h.Class(ui.NavItemIcon), IconHome()),
					h.Span(h.Class(ui.NavItemLabel), g.Text("Todos")),
				),
			),
		),
		h.Div(
			h.Class(ui.SidebarFooter),
			h.Button(
				h.Class(ui.SidebarToggle),
				data.On("click", "$sidebarCollapsed = !$sidebarCollapsed"),
				g.Attr("aria-label", "Toggle sidebar"),
				h.Type("button"),
				h.Span(
					h.Class(ui.NavItemIcon),
					data.Show("!$sidebarCollapsed"),
					IconChevronLeft(),
				),
				h.Span(
					h.Class(ui.NavItemIcon),
					data.Show("$sidebarCollapsed"),
					IconChevronRight(),
				),
				h.Span(h.Class(ui.SidebarToggleText), g.Text("Collapse")),
			),
		),
	)
}

// MobileSidebar renders the mobile overlay sidebar.
func MobileSidebar() g.Node {
	return g.Group{
		h.Div(
			h.Class(ui.SidebarOverlay),
			data.Show("$sidebarOpen"),
			data.On("click", "$sidebarOpen = false"),
		),
		h.Aside(
			h.Class(classNames(ui.AppSidebar, ui.AppSidebarMobile)),
			h.Data("class", "{'app-sidebar--open': $sidebarOpen}"),
			h.Div(
				h.Class(ui.SidebarHeader),
				h.A(
					h.Href("/"),
					h.Class(ui.SidebarLogo),
					data.On("click", "$sidebarOpen = false"),
					h.Span(h.Class(ui.SidebarLogoIcon), g.Text("DS")),
					h.Span(h.Class(ui.SidebarLogoText), g.Text("Datastar")),
				),
			),
			h.Nav(
				h.Class(ui.SidebarNav),
				h.Div(
					h.Class(ui.NavSection),
					h.A(
						h.Href("/"),
						h.Class(classNames(ui.NavItem, ui.NavItemActive)),
						data.On("click", "$sidebarOpen = false"),
						h.Span(h.Class(ui.NavItemIcon), IconHome()),
						h.Span(h.Class(ui.NavItemLabel), g.Text("Todos")),
					),
				),
			),
		),
	}
}
