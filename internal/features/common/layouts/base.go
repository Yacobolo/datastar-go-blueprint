// Package layouts provides shared page layout nodes for gomponents views.
package layouts

import (
	"github.com/yacobolo/datastar-go-blueprint/internal/config"
	"github.com/yacobolo/datastar-go-blueprint/internal/features/common/components"
	appds "github.com/yacobolo/datastar-go-blueprint/internal/platform/ds"
	"github.com/yacobolo/datastar-go-blueprint/internal/ui"
	"github.com/yacobolo/datastar-go-blueprint/web/resources"

	g "maragu.dev/gomponents"
	data "maragu.dev/gomponents-datastar"
	h "maragu.dev/gomponents/html"
)

func baseSignals() map[string]any {
	return map[string]any{
		"theme":            "system",
		"sidebarCollapsed": false,
		"sidebarOpen":      false,
	}
}

// Base renders the common application shell.
func Base(title string, children ...g.Node) g.Node {
	bodyChildren := g.Group{
		data.Signals(baseSignals()),
	}

	if config.Global.Environment == config.Dev {
		bodyChildren = append(bodyChildren,
			h.Div(
				data.Init(appds.Get(
					"/reload",
					appds.Opt("retryMaxCount", "1000"),
					appds.Opt("retryInterval", "20"),
					appds.Opt("retryMaxWaitMs", "200"),
				)),
			),
			g.El("datastar-inspector"),
		)
	}

	bodyChildren = append(bodyChildren,
		h.Div(
			h.Class(ui.AppShell),
			components.Sidebar(),
			h.Div(
				h.Class(ui.AppContent),
				components.Header(title),
				h.Main(
					h.Class(ui.AppMain),
					g.Group(children),
				),
			),
		),
		components.MobileSidebar(),
		h.Script(
			h.Type("module"),
			h.Src(resources.StaticPath("theme-switcher.js")),
		),
	)

	return h.Doctype(
		h.HTML(
			h.Lang("en"),
			h.Head(
				h.TitleEl(g.Text(title)),
				h.Meta(h.Name("viewport"), h.Content("width=device-width, initial-scale=1, maximum-scale=1, user-scalable=0")),
				h.Script(g.Raw(`
				(function() {
					var stored = localStorage.getItem('theme') || 'system';
					var prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
					var resolved = stored === 'system' ? (prefersDark ? 'dark' : 'light') : stored;
					document.documentElement.dataset.theme = resolved;
				})();
			`)),
				h.Link(h.Rel("preconnect"), h.Href("https://fonts.googleapis.com")),
				h.Link(h.Rel("preconnect"), h.Href("https://fonts.gstatic.com"), g.Attr("crossorigin")),
				h.Link(
					h.Href("https://fonts.googleapis.com/css2?family=Fira+Code:wght@300..700&family=Inter:wght@100..900&family=Gideon+Roman:ital,wght@0,300;0,400;0,700;0,900;1,300;1,400;1,700;1,900&display=swap"),
					h.Rel("stylesheet"),
				),
				h.Script(h.Src("https://code.iconify.design/iconify-icon/2.1.0/iconify-icon.min.js")),
				h.Link(h.Href(resources.StaticPath("styles.css")), h.Rel("stylesheet"), h.Type("text/css")),
				h.Script(h.Type("module"), h.Src("https://cdn.jsdelivr.net/gh/starfederation/datastar@1.0.0-RC.7/bundles/datastar.js")),
				h.Script(h.Type("module"), h.Src(resources.StaticPath("libs/index.js"))),
			),
			h.Body(bodyChildren...),
		),
	)
}
