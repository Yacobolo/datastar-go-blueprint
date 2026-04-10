package components

import (
	"time"

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

func toastTypeClass(t ToastType) string {
	switch t {
	case ToastSuccess:
		return "alert-success"
	case ToastError:
		return "alert-error"
	case ToastInfo:
		return "alert-info"
	case ToastWarning:
		return "alert-warning"
	default:
		return "alert-info"
	}
}

// Toast renders a transient toast message.
func Toast(message string, toastType ToastType) g.Node {
	return h.Div(
		h.Class(classNames("alert shadow-lg transition-opacity duration-300", toastTypeClass(toastType))),
		g.Attr("role", "alert"),
		data.Init(
			"el.classList.add('opacity-0'); setTimeout(() => el.remove(), 300)",
			data.ModifierDelay,
			data.Duration(3*time.Second),
		),
		h.Span(g.Text(message)),
	)
}
