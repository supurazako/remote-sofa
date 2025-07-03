package stream

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/supurazako/remote-sofa/server/video"
)

type State int

const (
	StateIdle State = iota 	// 0: waiting (default)
	StateProcessing			// 1: processing
	StateCompleted			// 2: compeleted
	StateFailed				// 3: failed
)

func (s State) String() string {
	return [...]string{"Idle", "Processing", "Completed", "Failed"}[s]
}

type Status struct {
	State			State
	PlaylistPath	string
	Error			error
}

// for DI
type ConverterFunc func(inputFile, outputDir string) error

type Manager struct {
	mu			sync.Mutex
	sessions	map[string]*Status
	converter	ConverterFunc
}

func NewManager(converter ConverterFunc) *Manager {
	return &Manager{
		sessions: make(map[string]*Status),
		converter: converter,
	}
}

func (m *Manager) Startconversion(sessionID, inputFile string) {
	// write session status may conflict, so we need to lock
	m.mu.Lock()
	m.sessions[sessionID] = &Status{State: StateProcessing}
	m.mu.Unlock()

	go func() {
		// NOTE: now we use a local temp directory, but change it to S3 storage in the future
		outputDir, err := os.MkdirTemp("", "hls_"+sessionID+"_+")
		if err != nil {
			m.updateStatusToFailed(sessionID, fmt.Errorf("failed to create temp dir: %w", err))
			return
		}

		// in test, we use a mock converter function, in production, we use ffmpeg
		err = m.converter(inputFile, outputDir)

		if err != nil {
			m.updateStatusToFailed(sessionID, err)
		} else {
			playlistPath := filepath.Join(outputDir, video.PlaylistFilename)
			m.updateStatusToCompleted(sessionID, playlistPath)
		}
	}()
}

func (m *Manager) GetStatus(sessionID string) (*Status, error) {
	// read session status may conflict, so we need to lock
	m.mu.Lock()
	defer m.mu.Unlock()

	status, ok := m.sessions[sessionID]
	if !ok {
		return nil, fmt.Errorf("session not found %s", sessionID)
	}
	return status, nil
}

func (m *Manager) updateStatusToCompleted(sessionID, playlistPath string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if status, ok := m.sessions[sessionID]; ok {
		status.State = StateCompleted
		status.PlaylistPath = playlistPath
	}
}

func (m *Manager) updateStatusToFailed(sessionID string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if status, ok := m.sessions[sessionID]; ok {
		status.State = StateFailed
		status.Error = err
	}
}
