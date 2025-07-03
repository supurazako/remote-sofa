package stream

import (
	"fmt"
	"sync"
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

// NOTE: logic is not implement yet
func (m *Manager) Startconversion(sessionID, inputFile string) {
	// some code
}

// NOTE: logic is not implement yet
func (m *Manager) GetStatus(sessionID string) (*Status, error) {
	return nil, fmt.Errorf("session not found: %s", sessionID)
}
