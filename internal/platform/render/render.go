// Package render provides gomponents rendering helpers for HTTP and SSE responses.
package render

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/starfederation/datastar-go/datastar"
	g "maragu.dev/gomponents"
)

// HTML renders a gomponents node as an HTML response.
func HTML(w http.ResponseWriter, node g.Node) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return node.Render(w)
}

// String renders a gomponents node to a string.
func String(node g.Node) (string, error) {
	var buf bytes.Buffer
	if err := node.Render(&buf); err != nil {
		return "", fmt.Errorf("render node: %w", err)
	}
	return buf.String(), nil
}

// PatchSSE renders a gomponents node and patches it through Datastar SSE.
func PatchSSE(sse *datastar.ServerSentEventGenerator, node g.Node, opts ...datastar.PatchElementOption) error {
	out, err := String(node)
	if err != nil {
		return err
	}

	return sse.PatchElements(out, opts...)
}
