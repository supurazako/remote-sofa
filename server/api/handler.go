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
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "session not found"})
		return
	}

	switch status.State {
	case stream.StateProcessing:
		respondJSON(w, http.StatusAccepted, map[string]string{
			"status": status.State.String(),
		})
	case stream.StateCompleted:
		respondJSON(w, http.StatusOK, map[string]string {
			"status":	status.State.String(),
			"url":		status.PlaylistPath,
		})
	case stream.StateFailed:
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"status": status.State.String(),
			"error":  status.Error.Error(),
		})
	default:
		// unexpected case
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "unknown session state",
		})
	}
}

// helper func for JSON response
func respondJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}
