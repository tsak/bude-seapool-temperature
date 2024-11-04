package main

import (
	"encoding/gob"
	"errors"
	"log/slog"
	"os"
	"sync"
	"time"
)

// State represents the state of the application
type State struct {
	LastRequest   time.Time
	ImageRedraws  int
	ImageRequests int
}

func (s State) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Time("last_request", s.LastRequest),
		slog.Int("image_redraws", s.ImageRedraws),
		slog.Int("image_requests", s.ImageRequests),
	)
}

// StateManager handles synchronization and persistence of application state.
type StateManager struct {
	sync.Mutex
	state    *State
	filename string
}

func NewStateManager(filename string, interval time.Duration) (*StateManager, error) {
	sm := StateManager{
		state:    &State{},
		filename: filename,
	}

	go sm.autoSave(interval)

	if err := sm.Load(); err != nil {
		return &sm, err
	}

	return &sm, nil
}

// autoSave periodically saves the current state to a file at the given interval.
func (sm *StateManager) autoSave(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		err := sm.Save()
		if err != nil {
			slog.Error("failed to save state", "error", err)
		}
		slog.Debug("automatically saved state", "state", sm.state, "interval", interval)
	}
}

func (sm *StateManager) Load() error {
	sm.Lock()
	defer sm.Unlock()

	f, err := os.Open(sm.filename)
	// If the state file doesn't exist, save empty state
	if errors.Is(err, os.ErrNotExist) {
		return sm.save()
	}

	if err = gob.NewDecoder(f).Decode(&sm.state); err != nil {
		slog.Error("failed to decode state", "error", err)
		err = os.Remove(sm.filename)
		if err != nil {
			return err
		}
		return sm.save()
	}

	return nil
}

func (sm *StateManager) save() error {
	f, err := os.Create(sm.filename)
	if err != nil {
		return err
	}

	return gob.NewEncoder(f).Encode(sm.state)
}

func (sm *StateManager) Save() error {
	sm.Lock()
	defer sm.Unlock()

	return sm.save()
}

func (sm *StateManager) SetLastRequest(dt time.Time) {
	sm.Lock()
	defer sm.Unlock()
	sm.state.LastRequest = dt
}

func (sm *StateManager) IncrementImageRedraws() {
	sm.Lock()
	defer sm.Unlock()
	sm.state.ImageRedraws++
}

func (sm *StateManager) IncrementImageRequests() {
	sm.Lock()
	defer sm.Unlock()
	sm.state.ImageRequests++
}
