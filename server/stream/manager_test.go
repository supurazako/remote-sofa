package stream

import (
	"errors"
	"sync"
	"testing"
	"time"
)

func TestStreamManger(t *testing.T) {
	t.Run("Successful conversion", func(t *testing.T) {
		// SET UP
		var wg sync.WaitGroup
		wg.Add(1) // wait for the goroutine complete

		// success case
		mockConverter := func(inputFile, outputDir string) error {
			t.Logf("Mock conversion called with: %s, %s", inputFile, outputDir)
			wg.Done() // signal that the conversion is done
			return nil
		}

		manager := NewManager(mockConverter)
		sessionID := "session-success"
		inputFile := "/path/to/video.mp4"

		// EXECUTE
		manager.Startconversion(sessionID, inputFile)

		// ASSERT
		status, err := manager.GetStatus(sessionID)
		if err != nil {
			t.Fatalf("GetStatus should not return error right after start, but got: %v", err)
		}
		if status.State != StateProcessing {
			t.Errorf("expected state to be %v, but got %v", StateProcessing, status.State)
		}

		if waitTimeout(&wg, 1*time.Second) {
			t.Fatal("timed out waiting for conversion to complete")
		}

		status, err = manager.GetStatus(sessionID)
		if err != nil {
			t.Fatalf("GetStatus should not return error after completion, but got: %v", err)
		}
		if status.State != StateCompleted {
			t.Errorf("expected state to be %v, but got %v", StateCompleted, status.State)
		}

		if status.PlaylistPath == "" {
			t.Error("expected playlist path to be set on completion, but it was empty")
		}
	})

	t.Run("Failed conversion", func(t *testing.T) {
		// SET UP
		var wg sync.WaitGroup
		wg.Add(1)
		expectedErr := errors.New("ffmpeg failed")

		// failure case
		mockConverter := func(inputFile, outputDir string) error {
			wg.Done()
			return expectedErr
		}

		manager := NewManager(mockConverter)
		sessionID := "session-fail"
		inputFile := "/path/to/invalid.mp4"

		// EXECUTE
		manager.Startconversion(sessionID, inputFile)

		// ASSERT
		// wait for async conversion to complete
		if (waitTimeout(&wg, 1*time.Second)) {
			t.Fatal("timed out waiting for conversion to complete")
		}

		status, err := manager.GetStatus(sessionID)
		if err != nil {
			t.Fatalf("GetStatus should not return error after failure, but got: %v", err)
		}
		if status.State != StateFailed {
			t.Errorf("expected state to be %v, but got %v", StateFailed, status.State)
		}
		if status.Error == nil {
			t.Error("expected error to be set on failure, but it was nil")
		}
		if status.Error.Error() != expectedErr.Error() {
			t.Errorf("expected error message '%v', but got '%v'", expectedErr, status.Error)
		}
	})
}

func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // positive case
	case <-time.After(timeout):
		return true // negative case, timed out
	}
}
