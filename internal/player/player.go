package player

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Dokito555/robin/internal/events"
	"github.com/Dokito555/robin/internal/storage"
	"github.com/gopxl/beep/v2"
)

type Player struct {
	mu       sync.Mutex
	engine   *Engine
	state    *State
	volume   *VolumeCtrl
	Queue    *Queue
	bus      *events.Bus
	stream   beep.StreamSeekCloser
	format   beep.Format
	ticker   *time.Ticker
	tickStop chan struct{}
}

func NewPlayer(bus *events.Bus) (*Player, error) {
	engine, err := NewEngine()
	if err != nil {
		return nil, fmt.Errorf("engine: %w", err)
	}

	state := NewState()

	p := &Player{
		engine:   engine,
		state:    state,
		volume:   NewVolumeCtrl(state, bus),
		Queue:    NewQueue(bus),
		bus:      bus,
		tickStop: make(chan struct{}),
	}

	return p, nil
}

func (p *Player) Play(track storage.Track) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.playTrack(&track)
}

func (p *Player) PlayCurrent() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	t := p.Queue.Current()
	if t == nil {
		return fmt.Errorf("queue is empty")
	}

	return p.playTrack(t)
}

func (p *Player) Pause() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state.SnapShot().Status != StatusPlaying {
		return
	}

	p.engine.Pause()
	p.state.setStatus(StatusPaused)
	p.stopTicker()
	p.bus.Publish(events.TrackPaused, nil)
}

func (p *Player) Resume() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state.SnapShot().Status != StatusPaused {
		return
	}

	p.engine.Resume()
	p.state.setStatus(StatusPlaying)
	p.startTicker()
	p.bus.Publish(events.TrackResumed, nil)
}

func (p *Player) Toggle() {
	snap := p.state.SnapShot()
	switch snap.Status {
	case StatusPlaying:
		p.Pause()
	case StatusPaused:
		p.Resume()
	}
}

func (p *Player) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.stopTicker()
	p.engine.Stop()
	p.stream = nil
	p.state.setStatus(StatusStopped)
	p.bus.Publish(events.TrackStopped, nil)
}

func (p *Player) Next() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	t := p.Queue.Next()
	if t == nil {
		p.engine.Stop()
		p.state.setStatus(StatusStopped)
		return nil
	}

	return p.playTrack(t)
}

func (p *Player) Prev() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	t := p.Queue.Prev()
	if t == nil {
		return nil
	}
	return p.playTrack(t)
}

func (p *Player) Seek(pos time.Duration) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.stream == nil {
		return fmt.Errorf("nothing playing")
	}

	if err := p.engine.Seek(p.stream, p.format, pos); err != nil {
		return err
	}

	p.state.setPosition(pos)
	p.bus.Publish(events.SeekPerformed, events.SeekPayload{
		PositionMS: pos.Milliseconds(),
	})
	return nil
}

func (p *Player) SetVolume(level int) {
	p.volume.Set(level)
}

func (p *Player) Mute() {
	p.volume.Mute()
}

func (p *Player) Unmute() {
	p.volume.Unmute()
}

func (p *Player) ToggleMute() {
	p.volume.Toggle()
}

func (p *Player) State() Snapshot {
	return p.state.SnapShot()
}

func (p *Player) playTrack(track *storage.Track) error {
	p.stopTicker()
	p.engine.Stop()

	stream, format, err := p.engine.Load(track.Path)
	if err != nil {
		return fmt.Errorf("load track: %w", err)
	}

	resampled := p.engine.resample(stream, format)

	volStream := p.volume.attach(resampled)

	p.stream = stream
	p.format = format
	p.state.setTrack(track)
	p.state.setStatus(StatusPlaying)

	p.engine.Play(volStream, func() {
		p.bus.Publish(events.TrackEnded, events.TrackEndedPayload{
			TrackID: track.ID,
		})
		log.Printf("player: track ended %s", track.Title)

		go func() {
			// small delay ensures playTrack has fully returned and p.mu is released
			time.Sleep(50 * time.Millisecond)
			if err := p.Next(); err != nil {
				log.Printf("player: auto-next error: %v", err)
			}
		}()
	})

	p.bus.Publish(events.TrackStarted, events.TrackStartedPayload{
		TrackID: track.ID,
		Title:   track.Title,
		Artist:  track.Artist,
		Album:   track.Album,
	})

	p.startTicker()
	return nil
}

// updates position in state every second during playback
func (p *Player) startTicker() {
	// p.tickStop = make(chan struct{})
	// p.ticker = time.NewTicker(time.Second)
	// go func() {
	// 	for {
	// 		select {
	// 		case <-p.ticker.C:
	// 			if p.stream != nil {
	// 				pos := p.engine.Position(p.stream, p.format)
	// 				p.state.setPosition(pos)
	// 			}
	// 		case <-p.tickStop:
	// 			return
	// 		}
	// 	}
	// }()

	// i have no idea what im doing here
	ticker := time.NewTicker(time.Second)
	stop := make(chan struct{})

	p.ticker = ticker
	p.tickStop = stop

	go func(t *time.Ticker, stopCh chan struct{}) {
		for {
			select {
			case <-t.C:
				p.mu.Lock()

				if p.stream != nil {
					pos := p.engine.Position(p.stream, p.format)
					p.state.setPosition(pos)
				}

				p.mu.Unlock()

			case <-stopCh:
				t.Stop()
				return
			}
		}
	}(ticker, stop)
}

func (p *Player) stopTicker() {
	if p.ticker != nil {
		close(p.tickStop)
		p.ticker = nil
		p.tickStop = nil
	}
}
