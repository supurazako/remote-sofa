package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/supurazako/remote-sofa/server/stream"
)

type mockStreamManager struct {
	status	*stream.Status
	err		error
}

func (m *mockStreamManager) GetStatus(sessionID string) (*stream.Status, error) {
	return m.status, m.err
}

// In test, do not use StartConversion. so blank
func (m *mockStreamManager) StartConversion(sessionID, inputFile string) {}

func TestStreamStatusHandler(t *testing.T) {
	t.Run("should return 202 when status is processing", func(t *testing.T) {
		// SET UP
		mockManager := &mockStreamManager{
			status: &stream.Status{State: stream.StateProcessing},
		}
		handler := NewHandler(mockManager)

		req := httptest.NewRequest("GET", "/api/sessions/test-session/stream", nil)
		rr := httptest.NewRecorder()

		// EXECUTE
		handler.ServeHTTP(rr, req)

		// ASSERT
		if rr.Code != http.StatusAccepted {
			t.Errorf("expected status code %d, but got %d", http.StatusAccepted, rr.Code)
		}

		var respBody map[string]string
		if err := json.Unmarshal(rr.Body.Bytes(), &respBody); err != nil {
			t.Fatalf("failed to unmarshal response body: %v", err)
		}
		if respBody["status"] != "Processing" {
			t.Errorf("expected status 'Processing', but got '%s'", respBody["status"])
		}
	})

	t.Run("should return 200 with URL when status is completed", func(t *testing.T) {
		// SET UP
		mockManager := &mockStreamManager{
			status: &stream.Status{
				State:			stream.StateCompleted,
				PlaylistPath:	"/streams/test-session/playlist.m3u8",
			},
		}
		handler := NewHandler(mockManager)

		req := httptest.NewRequest("GET", "/api/sessions/test-session/stream", nil)
		rr := httptest.NewRecorder()

		// EXECUTE
		handler.ServeHTTP(rr, req)

		// ASSERT
		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, but got %d", http.StatusOK, rr.Code)
		}

		var respBody map[string]string
		if err := json.Unmarshal(rr.Body.Bytes(), &respBody); err != nil {
			t.Fatalf("failed to unmarshal response body: %v", err)
		}
		if respBody["status"] != "Completed" {
			t.Errorf("expected status 'Completed', but got '%s'", respBody["status"])
		}
		url := "/streams/test-session/playlist.m3u8"
		if respBody["url"] != url {
			t.Errorf("expected url '%s', but got '%s'", url, respBody["url"])
		}
	})

	t.Run("should return 404 when session is not found", func(t *testing.T) {
		// SET UP
		// expect: streamManager return error
		mockManager := &mockStreamManager{
			status:	nil,
			err:	http.ErrMissingFile,
		}
		handler := NewHandler(mockManager)

		req := httptest.NewRequest("GET", "/api/sessions/not-found/stream", nil)
		rr := httptest.NewRecorder()

		// EXECUTE
		handler.ServeHTTP(rr, req)

		// ASSERT
		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status code %d, but got %d", http.StatusNotFound, rr.Code)
		}
	})
}
