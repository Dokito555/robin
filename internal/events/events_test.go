// events/events_test.go
package events

import (
	"sync"
	"testing"
	"time"
)

func TestSubscribeAndPublish(t *testing.T) {
	bus := NewBus()

	var received Event
	var wg sync.WaitGroup
	wg.Add(1)

	bus.Subscribe(TrackStarted, func(e Event) {
		received = e
		wg.Done()
	})

	bus.Publish(TrackStarted, TrackStartedPayload{
		TrackID: 1,
		Title:   "Numb",
		Artist:  "Linkin Park",
	})

	wg.Wait()

	if received.Type != TrackStarted {
		t.Errorf("expected TrackStarted, got %s", received.Type)
	}

	p, ok := received.Payload.(TrackStartedPayload)
	if !ok {
		t.Fatal("payload type assertion failed")
	}
	if p.Title != "Numb" {
		t.Errorf("expected Numb, got %s", p.Title)
	}
}

func TestUnsubscribe(t *testing.T) {
	bus := NewBus()
	count := 0

	unsub := bus.Subscribe(TrackEnded, func(e Event) {
		count++
	})

	bus.Publish(TrackEnded, TrackEndedPayload{TrackID: 1})
	time.Sleep(20 * time.Millisecond)

	unsub() // remove handler

	bus.Publish(TrackEnded, TrackEndedPayload{TrackID: 1})
	time.Sleep(20 * time.Millisecond)

	if count != 1 {
		t.Errorf("expected handler to fire once, fired %d times", count)
	}
}

func TestSubscribeAll(t *testing.T) {
	bus := NewBus()

	var received []EventType
	var mu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(3)

	bus.SubscribeAll(func(e Event) {
		mu.Lock()
		received = append(received, e.Type)
		mu.Unlock()
		wg.Done()
	})

	bus.Publish(TrackStarted, nil)
	bus.Publish(VolumeChanged, nil)
	bus.Publish(ScanDone, nil)

	wg.Wait()

	if len(received) != 3 {
		t.Errorf("expected 3 events, got %d", len(received))
	}
}