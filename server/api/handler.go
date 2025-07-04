package api

import (
	"encoding/json"
	"net/http"
	"path"

	"github.com/supurazako/remote-sofa/server/stream"
)

type StreamManager interface {
	GetStatus(sessionID string) (*stream.Status, error)
}

type Handler struct {
	manager StreamManager
}

func NewHandler(manager StreamManager) http.Handler {
	return &Handler{manager: manager}
}

func(h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// extract sessionID from path
	sessionID := path.Base(path.Dir(r.URL.Path))

	status, err := h.manager.GetStatus(sessionID)
	if err != nil {
		http.Error(w, `{"error": "session not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	switch status.State {
	case stream.StateProcessing:
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{
			"status": status.State.String(),
		})

	case stream.StateCompleted:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":	status.State.String(),
			"url":		status.PlaylistPath,
		})
	default:
		// unexpected case
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "unknown session state",
		})
	}
}
