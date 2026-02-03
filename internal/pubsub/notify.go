package pubsub

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
	"github.com/yourusername/datastar-go-starter-kit/features/common/components"
)

// UpdateMessage is the payload sent over NATS for UI updates
type UpdateMessage struct {
	RefreshTodos bool       `json:"refreshTodos,omitempty"`
	Toast        *ToastData `json:"toast,omitempty"`
}

type ToastData struct {
	Message string               `json:"message"`
	Type    components.ToastType `json:"type"`
}

// NotifyOption is a functional option for building UpdateMessage
type NotifyOption func(*UpdateMessage)

// WithRefresh signals that the TODO list should be refreshed
func WithRefresh() NotifyOption {
	return func(m *UpdateMessage) {
		m.RefreshTodos = true
	}
}

// WithToast adds a toast notification
func WithToast(msg string, toastType components.ToastType) NotifyOption {
	return func(m *UpdateMessage) {
		m.Toast = &ToastData{
			Message: msg,
			Type:    toastType,
		}
	}
}

// Notify publishes an update message to the given NATS subject
func Notify(nc *nats.Conn, subject string, opts ...NotifyOption) error {
	msg := UpdateMessage{}
	for _, opt := range opts {
		opt(&msg)
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return nc.Publish(subject, data)
}

// ParseUpdateMessage unmarshals a NATS message into UpdateMessage
func ParseUpdateMessage(data []byte) (UpdateMessage, error) {
	var msg UpdateMessage
	err := json.Unmarshal(data, &msg)
	return msg, err
}
