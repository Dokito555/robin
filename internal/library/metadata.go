package library

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bogem/id3v2/v2"
)

type TrackMeta struct {
	Title     string
	Artist    string
	Album     string
	Genre     string
	Duration  time.Duration
	Path      string
	CoverData []byte
}

// TODO: check for corrupted metadata
func ReadMeta(path string) (TrackMeta, error) {
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".mp3":
			return readMP3(path)
	case ".flac", ".ogg", ".wav", ".m4a", ".aac":
		return readGeneric(path)
	default:
		return TrackMeta{}, fmt.Errorf("unsupported format: %s", ext)

	}
}

func ReadMetaSafe(path string) (meta TrackMeta, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic reading tags from %s: %v", path, err)
		}
	}()
	return ReadMeta(path)
}

func readMP3(path string) (TrackMeta, error) {
	tag, err := id3v2.Open(path, id3v2.Options{Parse: true})
	if err != nil {
		return fallback(path), nil
	}
	defer tag.Close()

	meta := TrackMeta{
		Path: path,
		Title: strings.TrimSpace(tag.Title()),
		Artist: strings.TrimSpace(tag.Artist()),
		Album: strings.TrimSpace(tag.Album()),
		Genre: strings.TrimSpace(tag.Genre()),
	}

	pictures := tag.GetFrames(tag.CommonID("Attached picture"))
	if len(pictures) > 0 {
		if pic, ok := pictures[0].(id3v2.PictureFrame); ok {
			meta.CoverData = pic.Picture
		}
	}
	
	meta.Duration = estimateDuration(path)
	
	return sanitize(meta), nil
}

// TODO: support multiple audio file formats
func readGeneric(path string) (TrackMeta, error) {
	meta := fallback(path)
	meta.Duration = estimateDuration(path)
	meta.CoverData = findFolderArt(filepath.Dir(path))

	return meta, nil
}

func fallback(path string) TrackMeta {
	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	parts := strings.SplitN(name, " - ", 2)
	meta := TrackMeta{Path: path}
	if len(parts) == 2 {
		meta.Artist = strings.TrimSpace(parts[0])
		meta.Title = strings.TrimSpace(parts[1])
	} else {
		meta.Title = strings.TrimSpace(name)
	}
	return meta
}

func sanitize(m TrackMeta) TrackMeta {
	if m.Title == "" {
		m.Title = strings.TrimSuffix(filepath.Base(m.Path), filepath.Ext(m.Path))
	}
	if m.Artist == "" {
		m.Artist = "Unknown Artist"
	}
	if m.Album == "" {
		m.Album = "Unknown Album"
	}

	if m.CoverData == nil {
		m.CoverData = findFolderArt(filepath.Dir(m.Path))
	}
	return m
}

func findFolderArt(dir string) []byte {
	candidates := []string {
		"cover.jpg", "cover.jpeg", "cover.png",
		"folder.jpg", "folder.jpeg", "folder.png",
		"front.jpg", "front.jpeg", "front.png",
		"album.jpg", "album.jpeg", "album.png",
	}
	for _, name := range candidates {
		p := filepath.Join(dir, name)
		data, err := os.ReadFile(p)
		if err == nil {
			return data
		}
	}

	return nil
}

func estimateDuration(path string) time.Duration {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}

	const avgBitsPerSec = 128_000
	secs := float64(info.Size()*8) / float64(avgBitsPerSec)
	return time.Duration(secs) * time.Second
}

func (m TrackMeta) validate() error {
	if m.Path == "" {
		return fmt.Errorf("empty path")
	}

	if m.Title == "" {
		return fmt.Errorf("could not determine title for %s", m.Path)
	}

	// large duration might be corrupted
	if m.Duration > 24*time.Hour {
		m.Duration = 0
	}

	if len(m.CoverData) > 0 && len(m.CoverData) < 128 {
		m.CoverData = nil
	}

	return nil
}