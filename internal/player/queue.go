package player

import (
	"math/rand"
	"sync"

	"github.com/Dokito555/robin/internal/events"
	"github.com/Dokito555/robin/internal/storage"
)

type RepeatMode int

const (
	RepeatNone RepeatMode = iota
	RepeatOne
	RepeatAll
)

type Queue struct {
	mu      sync.Mutex
	tracks  []storage.Track
	index   int
	shuffle bool
	repeat  RepeatMode
	bus     *events.Bus
}

func NewQueue(bus *events.Bus) *Queue {
	return &Queue{
		bus: bus,
		index: -1,
	}
}

func (q *Queue) Add(tracks ...storage.Track) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.tracks = append(q.tracks, tracks...)
	q.publish()
}

func (q *Queue) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.tracks = nil
	q.index = -1
	q.bus.Publish(events.QueueCleared, nil)
}

// replace queue entirely
func (q *Queue) SetTracks(tracks []storage.Track) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.tracks = make([]storage.Track, len(tracks))
	copy(q.tracks, tracks)

	q.index = 0
	if q.shuffle {
		q.doShuffle()
	}

	q.publish()
}

func (q *Queue) Current() *storage.Track {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.index < 0 || q.index >= len(q.tracks) {
		return nil
	}
	
	t := q.tracks[q.index]
	return &t
}

func (q *Queue) Next() *storage.Track {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.tracks) == 0 {
		return nil
	}

	switch q.repeat {
	case RepeatOne:
	case RepeatAll:
		q.index = (q.index + 1) % len(q.tracks)
	default:
		if q.index + 1 >= len(q.tracks) {
			return nil
		}
		q.index++
	}

	t := q.tracks[q.index]
	return &t
}

func (q *Queue) Prev() *storage.Track {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.tracks) == 0 {
		return nil
	}

	if q.index > 0 {
		q.index--
	}

	t := q.tracks[q.index]
	return &t
}

func (q *Queue) JumpTo(index int) *storage.Track {
	q.mu.Lock()
	defer q.mu.Unlock()

	if index < 0 || index >= len(q.tracks) {
		return nil
	}

	q.index = index
	t := q.tracks[q.index]
	return &t
}

func (q *Queue) SetShuffle(on bool) {
	q.mu.Lock()

	q.shuffle = on
	if on {
		q.doShuffle()
	}

	q.mu.Unlock()
	
	if on {
		q.publish()
	}
}

func (q *Queue) SetRepeat(mode RepeatMode) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.repeat = mode
}

func (q *Queue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return len(q.tracks)
}

func (q *Queue) doShuffle() {
	rand.Shuffle(len(q.tracks), func(i, j int) {
		q.tracks[i], q.tracks[j] = q.tracks[j], q.tracks[i]
	})

	q.index = 0
}

func (q *Queue) publish() {
	q.bus.Publish(events.QueueChanged, nil)
}