// Package core provides reusable UI building blocks for the app shell and shared views.
package core

import (
	"fmt"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func classNames(classes ...string) string {
	filtered := make([]string, 0, len(classes))
	for _, class := range classes {
		if class != "" {
			filtered = append(filtered, class)
		}
	}
	return strings.Join(filtered, " ")
}

func when(condition bool, value string) string {
	if condition {
		return value
	}
	return ""
}

// SseIndicator renders a shared loading indicator bound to a Datastar signal.
func SseIndicator(signalName string) g.Node {
	return h.Span(
		h.Data(
			"class",
			fmt.Sprintf("{'loading loading-dots loading-xs text-base-content/60': $%s, 'hidden': !$%s}", signalName, signalName),
		),
	)
}
