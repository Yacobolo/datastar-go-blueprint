// Package layouts provides shared page layout nodes for gomponents views.
package layouts

import (
	"github.com/yacobolo/datastar-go-blueprint/internal/config"
	"github.com/yacobolo/datastar-go-blueprint/internal/features/common/components"
	appds "github.com/yacobolo/datastar-go-blueprint/internal/platform/ds"
	"github.com/yacobolo/datastar-go-blueprint/web/resources"

	g "maragu.dev/gomponents"
	data "maragu.dev/gomponents-datastar"
	h "maragu.dev/gomponents/html"
)

func baseSignals() map[string]any {
	return map[string]any{
		"theme":            "system",
		"resolvedTheme":    "light",
		"sidebarCollapsed": false,
		"sidebarOpen":      false,
	}
}

// Base renders the common application shell.
func Base(title string, children ...g.Node) g.Node {
	bodyChildren := g.Group{
		data.Signals(baseSignals()),
		data.Init("$theme = window.themeToggle.get(); $resolvedTheme = window.themeToggle.resolved()"),
		data.On("themechange", "$theme = evt.detail.mode; $resolvedTheme = evt.detail.resolved", data.ModifierWindow),
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
			h.Class("drawer min-h-screen bg-base-200 lg:drawer-open"),
			h.Input(
				h.ID("app-drawer"),
				h.Type("checkbox"),
				h.Class("drawer-toggle"),
				data.Bind("sidebarOpen"),
				g.Attr("aria-label", "Toggle navigation"),
			),
			h.Div(
				h.Class("drawer-content flex min-h-screen flex-col"),
				components.Header(title),
				h.Main(
					h.Class("mx-auto flex w-full max-w-6xl flex-1 flex-col gap-6 p-4 sm:p-6"),
					g.Group(children),
				),
			),
			components.Sidebar(),
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
				h.Script(h.Src("https://code.iconify.design/iconify-icon/2.1.0/iconify-icon.min.js")),
				h.Link(h.Href(resources.StaticPath("styles.css")), h.Rel("stylesheet"), h.Type("text/css")),
				h.Script(h.Type("module"), h.Src(resources.StaticPath("theme-switcher.js"))),
				h.Script(h.Type("module"), h.Src("https://cdn.jsdelivr.net/gh/starfederation/datastar@1.0.0-RC.7/bundles/datastar.js")),
				h.Script(h.Type("module"), h.Src(resources.StaticPath("libs/index.js"))),
			),
			h.Body(append(g.Group{
				h.Class("min-h-screen bg-base-200 text-base-content"),
			}, bodyChildren...)...),
		),
	)
}
