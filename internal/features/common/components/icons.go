package components

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func icon(children ...g.Node) g.Node {
	nodes := g.Group{
		h.Class("h-5 w-5 shrink-0"),
		g.Attr("viewBox", "0 0 24 24"),
		g.Attr("fill", "none"),
		g.Attr("stroke", "currentColor"),
		g.Attr("stroke-width", "2"),
		g.Attr("stroke-linecap", "round"),
		g.Attr("stroke-linejoin", "round"),
	}
	nodes = append(nodes, children...)
	return h.SVG(nodes...)
}

// IconSun renders a sun icon for light mode.
func IconSun() g.Node {
	return icon(
		g.El("circle", g.Attr("cx", "12"), g.Attr("cy", "12"), g.Attr("r", "4")),
		g.El("path", g.Attr("d", "M12 2v2")),
		g.El("path", g.Attr("d", "M12 20v2")),
		g.El("path", g.Attr("d", "m4.93 4.93 1.41 1.41")),
		g.El("path", g.Attr("d", "m17.66 17.66 1.41 1.41")),
		g.El("path", g.Attr("d", "M2 12h2")),
		g.El("path", g.Attr("d", "M20 12h2")),
		g.El("path", g.Attr("d", "m6.34 17.66-1.41 1.41")),
		g.El("path", g.Attr("d", "m19.07 4.93-1.41 1.41")),
	)
}

// IconMoon renders a moon icon for dark mode.
func IconMoon() g.Node {
	return icon(
		g.El("path", g.Attr("d", "M12 3a6 6 0 0 0 9 9 9 9 0 1 1-9-9Z")),
	)
}

// IconMonitor renders a monitor/computer icon for system theme.
func IconMonitor() g.Node {
	return icon(
		g.El("rect", g.Attr("width", "20"), g.Attr("height", "14"), g.Attr("x", "2"), g.Attr("y", "3"), g.Attr("rx", "2")),
		g.El("line", g.Attr("x1", "8"), g.Attr("x2", "16"), g.Attr("y1", "21"), g.Attr("y2", "21")),
		g.El("line", g.Attr("x1", "12"), g.Attr("x2", "12"), g.Attr("y1", "17"), g.Attr("y2", "21")),
	)
}

// IconMenu renders a hamburger menu icon for mobile navigation.
func IconMenu() g.Node {
	return icon(
		g.El("line", g.Attr("x1", "4"), g.Attr("x2", "20"), g.Attr("y1", "12"), g.Attr("y2", "12")),
		g.El("line", g.Attr("x1", "4"), g.Attr("x2", "20"), g.Attr("y1", "6"), g.Attr("y2", "6")),
		g.El("line", g.Attr("x1", "4"), g.Attr("x2", "20"), g.Attr("y1", "18"), g.Attr("y2", "18")),
	)
}

// IconHome renders a home icon.
func IconHome() g.Node {
	return icon(
		g.El("path", g.Attr("d", "m3 9 9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z")),
		g.El("polyline", g.Attr("points", "9 22 9 12 15 12 15 22")),
	)
}

// IconChevronLeft renders a left chevron icon.
func IconChevronLeft() g.Node {
	return icon(
		g.El("path", g.Attr("d", "m15 18-6-6 6-6")),
	)
}

// IconChevronRight renders a right chevron icon.
func IconChevronRight() g.Node {
	return icon(
		g.El("path", g.Attr("d", "m9 18 6-6-6-6")),
	)
}
