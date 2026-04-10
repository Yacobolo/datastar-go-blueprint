package core

import (
	lucide "github.com/eduardolat/gomponents-lucide"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func icon(classes ...string) g.Node {
	return h.Class(classNames(append([]string{"h-5 w-5 shrink-0"}, classes...)...))
}

// IconSun renders a sun icon for light mode.
func IconSun() g.Node {
	return lucide.Sun(icon())
}

// IconMoon renders a moon icon for dark mode.
func IconMoon() g.Node {
	return lucide.Moon(icon())
}

// IconMenu renders a hamburger menu icon for mobile navigation.
func IconMenu() g.Node {
	return lucide.Menu(icon())
}

// IconHome renders a home icon.
func IconHome() g.Node {
	return lucide.House(icon())
}

// IconChevronLeft renders a left chevron icon.
func IconChevronLeft() g.Node {
	return lucide.PanelLeftClose(icon())
}

// IconChevronRight renders a right chevron icon.
func IconChevronRight() g.Node {
	return lucide.PanelLeftOpen(icon())
}

// IconTrash renders a trash icon.
func IconTrash() g.Node {
	return lucide.Trash2(icon())
}

// IconListChecks renders a list/check management icon.
func IconListChecks() g.Node {
	return lucide.ListChecks(icon())
}

// IconCheckbox renders a todo checkbox icon.
func IconCheckbox(completed bool) g.Node {
	if completed {
		return lucide.SquareCheck(icon())
	}
	return lucide.Square(icon())
}

// IconClose renders a close icon.
func IconClose() g.Node {
	return lucide.X(icon())
}
