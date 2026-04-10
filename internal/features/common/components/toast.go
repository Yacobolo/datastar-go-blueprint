package components

import (
	"time"

	"github.com/yacobolo/datastar-go-blueprint/internal/ui"

	g "maragu.dev/gomponents"
	data "maragu.dev/gomponents-datastar"
	h "maragu.dev/gomponents/html"
)

// ToastType identifies the visual style of a toast notification.
type ToastType string

const (
	// ToastSuccess styles a successful toast notification.
	ToastSuccess ToastType = "success"
	// ToastError styles an error toast notification.
	ToastError ToastType = "error"
	// ToastInfo styles an informational toast notification.
	ToastInfo ToastType = "info"
	// ToastWarning styles a warning toast notification.
	ToastWarning ToastType = "warning"
)

// toastTypeClass returns the type-safe CSS constant for the toast type.
func toastTypeClass(t ToastType) string {
	switch t {
	case ToastSuccess:
		return ui.ToastSuccess
	case ToastError:
		return ui.ToastError
	case ToastInfo:
		return ui.ToastInfo
	default:
		return ui.ToastInfo
	}
}

// Toast renders a transient toast message.
func Toast(message string, toastType ToastType) g.Node {
	return h.Div(
		h.Class(classNames(ui.Toast, toastTypeClass(toastType))),
		data.Init(
			"el.style.opacity = '0'; setTimeout(() => el.remove(), 300)",
			data.ModifierDelay,
			data.Duration(3*time.Second),
		),
		h.Span(g.Text(message)),
	)
}
