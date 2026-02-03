// Package handlers provides shared HTTP handler utilities.
package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/starfederation/datastar-go/datastar"
)

// RequireSession extracts the session ID from the request, creating one if needed.
// Returns the session ID and true on success.
// On failure, writes an HTTP 500 error and returns ("", false).
func RequireSession(store sessions.Store, w http.ResponseWriter, r *http.Request) (string, bool) {
	sess, err := store.Get(r, "connections")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", false
	}

	id, ok := sess.Values["id"].(string)
	if !ok {
		id = uuid.New().String()
		sess.Values["id"] = id
		if err := sess.Save(r, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return "", false
		}
	}
	return id, true
}

// RequireIDParam parses an int64 ID from the given chi URL parameter.
// Returns the ID and true on success.
// On failure, writes an HTTP 400 error and returns (0, false).
func RequireIDParam(w http.ResponseWriter, r *http.Request, param string) (int64, bool) {
	id, err := strconv.ParseInt(chi.URLParam(r, param), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

// RequireIntParam parses an int from the given chi URL parameter.
// Returns the value and true on success.
// On failure, writes an HTTP 400 error and returns (0, false).
func RequireIntParam(w http.ResponseWriter, r *http.Request, param string) (int, bool) {
	val, err := strconv.Atoi(chi.URLParam(r, param))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 0, false
	}
	return val, true
}

// LogConsoleError sends an error to the browser console via SSE.
// If sending fails, logs the failure.
func LogConsoleError(sse *datastar.ServerSentEventGenerator, err error) {
	if err := sse.ConsoleError(err); err != nil {
		slog.Error("failed to send console error", "error", err)
	}
}

// NewSSEWithSignals reads the request body, creates an SSE generator, and unmarshals
// signals from the body. Use this for POST/PATCH/DELETE handlers that need both
// SSE responses and signal data.
//
// Returns the SSE generator and true on success.
// On failure, sends error to client (http.Error or SSE console) and returns nil, false.
func NewSSEWithSignals(w http.ResponseWriter, r *http.Request, signals any) (*datastar.ServerSentEventGenerator, bool) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, false
	}

	sse := datastar.NewSSE(w, r)

	if err := json.Unmarshal(body, signals); err != nil {
		bodyPreview := string(body)
		if len(bodyPreview) > 500 {
			bodyPreview = bodyPreview[:500] + "..."
		}
		LogConsoleError(sse, fmt.Errorf("failed to unmarshal signals: %w, body: %s", err, bodyPreview))
		return nil, false
	}

	return sse, true
}

// PatchSignals sends signal patches via SSE.
// Logs an error if the patch fails.
func PatchSignals(sse *datastar.ServerSentEventGenerator, signals map[string]any) {
	data, err := json.Marshal(signals)
	if err != nil {
		slog.Error("failed to marshal signals", "error", err)
		return
	}
	if err := sse.PatchSignals(data); err != nil {
		slog.Error("failed to patch signals", "error", err)
	}
}
