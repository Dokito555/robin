package player

import (
	"math"

	"github.com/Dokito555/robin/internal/events"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
)

type VolumeCtrl struct {
	streamer *effects.Volume
	state    *State
	bus      *events.Bus
}

func NewVolumeCtrl(state *State, bus *events.Bus) *VolumeCtrl {
	return &VolumeCtrl{
		state: state,
		bus:   bus,
	}
}

// attach wraps a streamer with volume control
// call this when loading a track
func (v *VolumeCtrl) attach(s beep.Streamer) *effects.Volume {
	vol := &effects.Volume{
		Streamer: s,
		Base:     2,
		Volume:   v.toDecibles(v.state.SnapShot().Volume),
		Silent:   v.state.SnapShot().Muted,
	}
	v.streamer = vol
	return vol
}

func (v *VolumeCtrl) Set(level int) {
	level = clamp(level, 0, 100)

	v.state.setVolume(level)
	if v.streamer != nil {
		v.streamer.Volume = v.toDecibles(level)
	}

	v.publish()
}

func (v *VolumeCtrl) Mute() {
	v.state.setMuted(true)
	if v.streamer != nil {
		v.streamer.Silent = false
	}

	v.publish()
}

func (v *VolumeCtrl) Unmute() {
	v.state.setMuted(false)
	if v.streamer != nil {
		v.streamer.Silent = false
	}

	v.publish()
}

func (v *VolumeCtrl) Toggle() {
	snap := v.state.SnapShot()
	if snap.Muted {
		v.Unmute()
	} else {
		v.Mute()
	}
}

// converts 0-100 linear scale to beep's logarithmic dB range
// 100 → 0dB (full), 50 → ~-6dB, 0 → silence via Mute
func (v *VolumeCtrl) toDecibles(level int) float64 {
	if level <= 0 {
		// silent without using the silent flag
		return -10
	}

	// map 0–100 to -10–0 dB logarithmically
	return math.Log2(float64(level)/100.0) * 2
}

func (v *VolumeCtrl) publish() {
	snap := v.state.SnapShot()
	v.bus.Publish(events.VolumeChanged, events.VolumePayload{
		Volume: snap.Volume,
		Muted:  snap.Muted,
	})
}

// beepStreamer is a local alias so volume.go doesn't import beep at top level
type beepStreamer = interface {
	Stream(samples [][2]float64) (n int, ok bool)
	Err() error
}
