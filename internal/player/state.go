package player

import (
	"sync"
	"time"

	"github.com/Dokito555/robin/internal/storage"
)

type Status int

const (
	StatusStopped Status = iota
	StatusPlaying
	StatusPaused
)

func (s Status) String() string {
	switch s {
	case StatusPlaying:
		return "playing"
	case StatusPaused:
		return "paused"
	default:
		return "stopped"
	}
}

type State struct {
	mu           sync.RWMutex
	status       Status
	currentTrack *storage.Track
	position     time.Duration
	duration     time.Duration
	volume       int
	muted        bool
}

func NewState() *State {
	return &State{
		status: StatusStopped,
		volume: 80,
	}
}

func (s *State) setStatus(status Status) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.status = status
}

func (s *State) setTrack(t *storage.Track) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.currentTrack = t
	s.position = 0

	if t != nil {
		s.duration = t.Duration
	} else {
		s.duration = 0
	}
}

func (s *State) setPosition(pos time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.position = pos
}

func (s *State) setVolume(v int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.volume = clamp(v, 0, 100)
}

func (s *State) setMuted(m bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.muted = m
}

// copy of state, safe to read without locks
type Snapshot struct {
	Status       Status
	CurrentTrack *storage.Track
	Position     time.Duration
	Duration     time.Duration
	Volume       int
	Muted        bool
}

func (s *State) SnapShot() Snapshot {
	s.mu.Lock()
	defer s.mu.Unlock()

	return Snapshot{
		Status:       s.status,
		CurrentTrack: s.currentTrack,
		Position:     s.position,
		Duration:     s.duration,
		Volume:       s.volume,
		Muted:        s.muted,
	}
}

func clamp(v, min, max int) int {
	if v < min {
		return min
	} 
	if v > max {
		return max
	}
	return v
}
