package events

import "sync"

type EventType string

const (
	TrackStarted  EventType = "track.started"
	TrackPaused   EventType = "track.paused"
	TrackResumed  EventType = "track.resumed"
	TrackStopped  EventType = "track.stopped"
	TrackEnded    EventType = "track.ended"
	SeekPerformed EventType = "track.seeked"

	QueueChanged EventType = "queue.changed"
	QueueCleared EventType = "queue.cleared"

	VolumeChanged EventType = "volume.changed"
	MuteToggled   EventType = "volume.muted"

	ScanStarted   EventType = "library.scan.started"
	ScanProgress  EventType = "library.scan.progress"
	ScanDone      EventType = "library.scan.done"
	TrackAdded    EventType = "library.track.added"
	TrackRemoved  EventType = "library.track.removed"
	FolderAdded   EventType = "library.folder.added"
	FolderRemoved EventType = "library.folder.removed"
)

type Event struct {
	Type    EventType
	Payload any
}

type TrackStartedPayload struct {
	TrackID int64
	Title   string
	Artist  string
	Album   string
}

type TrackAddedPayload struct {
	TrackID int64
	Path    string
}

type TrackEndedPayload struct {
	TrackID int64
}

type SeekPayload struct {
	PositionMS int64
}

type VolumePayload struct {
	Volume int
	Muted  bool
}

type ScanProgressPayload struct {
	FolderID int64
	Current  int
	Total    int
}

type ScanDonePayload struct {
	FolderID    int64
	TracksAdded int
}

type TrackRemovedPayload struct {
	Path string
}

type FolderPayload struct {
	FolderID int64
	Path     string
}

type Handler func(Event)

type Bus struct {
	mu       sync.RWMutex
	handlers map[EventType][]Handler
}

func NewBus() *Bus {
	return &Bus{
		handlers: make(map[EventType][]Handler),
	}
}

func (b *Bus) Subscribe(eventType EventType, h Handler) func() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers[eventType] = append(b.handlers[eventType], h)

	idx := len(b.handlers[eventType]) - 1

	return func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		handlers := b.handlers[eventType]
		if idx < len(handlers) {
			b.handlers[eventType] = append(handlers[:idx], handlers[idx+1:]...)
		}
	}
}

func (b *Bus) SubscribeAll(h Handler) func() {
	b.mu.Lock()
	defer b.mu.Unlock()

	const wildcard EventType = "*"
	b.handlers[wildcard] = append(b.handlers[wildcard], h)
	idx := len(b.handlers[wildcard]) - 1

	return func() {
		b.mu.Unlock()
		handlers := b.handlers[wildcard]
		if idx < len(handlers) {
			b.handlers[wildcard] = append(handlers[:idx], handlers[idx+1:]...)
		}
	}
}

func (b *Bus) Publish(eventType EventType, payload any) {
	b.mu.RLock()
	specific := make([]Handler, len(b.handlers[eventType]))
	copy(specific, b.handlers[eventType])
	wildcard := make([]Handler, len(b.handlers["*"]))
	copy(wildcard, b.handlers["*"])
	b.mu.RUnlock()

	e := Event{
		Type:    eventType,
		Payload: payload,
	}

	go func() {
		for _, h := range specific {
			h(e)
		}
		for _, h := range wildcard {
			h(e)
		}
	}()
}
