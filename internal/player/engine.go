package player

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/flac"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/wav"
)

const (
	sampleRate = beep.SampleRate(44100)
	bufferSize = 4096
)

// owns audio output device and decodes audio files into PCM
type Engine struct {
	format  beep.Format
	ctrl    *beep.Ctrl
	current beep.StreamSeekCloser
}

func NewEngine() (*Engine, error) {
	if err := speaker.Init(sampleRate, bufferSize); err != nil {
		return nil, fmt.Errorf("speaker init: %w", err)
	}

	return &Engine{}, nil
}

// opens an audio file, decodes it, and returns the stream + format
// caller is responsible for closing the stream
func (e *Engine) Load(path string) (beep.StreamSeekCloser, beep.Format, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, beep.Format{}, fmt.Errorf("open file: %w", err)
	}

	// TODO: add support for ogg, m4a, aac formats
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".mp3":
		stream, format, err := mp3.Decode(f)
		if err != nil {
			f.Close()
			return nil, beep.Format{}, fmt.Errorf("decode mp3: %w", err)
		}
		return stream, format, nil
	case ".flac":
		stream, format, err := flac.Decode(f)
		if err != nil {
			f.Close()
			return nil, beep.Format{}, fmt.Errorf("decode flac: %w", err)
		}
		return stream, format, nil
	case ".wav":
		stream, format, err := wav.Decode(f)
		if err != nil {
			f.Close()
			return nil, beep.Format{}, fmt.Errorf("decode wav: %w", err)
		}
		return stream, format, nil
	default:
		f.Close()
		return nil, beep.Format{}, fmt.Errorf("unssuported format: %s", ext)
	}
}

// sends a streamer to the speaker
// onDone is called when the stream finishes naturally
func (e *Engine) Play(stream beep.Streamer, onDone func()) {
	speaker.Clear()

	ctrl := &beep.Ctrl{
		Streamer: stream,
		Paused:   false,
	}
	e.ctrl = ctrl

	speaker.Play(beep.Seq(ctrl, beep.Callback(onDone)))
}

func (e *Engine) Pause() {
	speaker.Lock()
	if e.ctrl != nil {
		e.ctrl.Paused = true
	}
	speaker.Unlock()
}

func (e *Engine) Resume() {
	speaker.Lock()
	if e.ctrl != nil {
		e.ctrl.Paused = false
	}
	speaker.Unlock()
}

func (e *Engine) Stop() {
	speaker.Clear()
	e.ctrl = nil
	if e.current != nil {
		e.current.Close()
		e.current = nil
	}
}

// moves playback to a position in the current stream
func (e *Engine) Seek(stream beep.StreamSeekCloser, format beep.Format, pos time.Duration) error {
	samplePos := format.SampleRate.N(pos)
	speaker.Lock()
	err := stream.Seek(samplePos)
	speaker.Unlock()
	return err
}

// returns the current playback position
func (e *Engine) Position(stream beep.StreamSeekCloser, format beep.Format) time.Duration {
	speaker.Lock()
	pos := format.SampleRate.D(stream.Position())
	speaker.Unlock()
	return pos
}

// wraps a stream if it's sample rate differes from the speaker rate
// TODO: should return beep.streamseekcloser
func (e *Engine) resample(stream beep.StreamSeekCloser, from beep.Format) beep.Streamer {
	if from.SampleRate == sampleRate {
		return stream
	}
	return beep.Resample(4, from.SampleRate, sampleRate, stream)
}
