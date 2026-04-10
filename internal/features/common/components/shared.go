package components

import (
	"fmt"
	"strings"

	"github.com/yacobolo/datastar-go-blueprint/internal/ui"

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

// KVPairsAttrs converts key/value pairs into gomponents attributes.
func KVPairsAttrs(kvPairs ...string) g.Group {
	if len(kvPairs)%2 != 0 {
		panic("kvPairs must be a multiple of 2")
	}

	attrs := make(g.Group, 0, len(kvPairs)/2)
	for i := 0; i < len(kvPairs); i += 2 {
		attrs = append(attrs, g.Attr(kvPairs[i], kvPairs[i+1]))
	}
	return attrs
}

// Icon renders an iconify icon.
func Icon(icon string, attrs ...string) g.Node {
	nodes := g.Group{
		g.Attr("icon", icon),
		g.Attr("noobserver"),
	}
	nodes = append(nodes, KVPairsAttrs(attrs...)...)
	return g.El("iconify-icon", nodes...)
}

// SseIndicator renders a shared loading indicator bound to a Datastar signal.
func SseIndicator(signalName string) g.Node {
	return h.Div(
		h.Class(ui.TextPrimary),
		h.Data("class", fmt.Sprintf("{'%s': $%s}", classNames(ui.Loading, ui.LoadingDots, ui.MlMd), signalName)),
	)
}
