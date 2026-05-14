package library

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var audioExtensions = map[string]struct{}{
	".mp3":  {},
	".flac": {},
	".ogg":  {},
	".wav":  {},
	".m4a":  {},
	".aac":  {},
}

// holds discovered audio file path
type ScanResult struct {
	Path string
}

func Scan(rootPath string) (<-chan ScanResult, error) {
	info, err := os.Stat(rootPath)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", rootPath)
	}

	ch := make(chan ScanResult, 64)

	go func ()  {
		defer close(ch)
		filepath.WalkDir(rootPath, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				// unreadable dir
				return nil
			}
			if d.IsDir() {
				name := d.Name()
				if strings.HasPrefix(name, ".") || strings.HasPrefix(name, "@") {
					return filepath.SkipDir
				}
				return nil
			}
			ext := strings.ToLower(filepath.Ext(path))
			if _, ok := audioExtensions[ext]; ok {
				ch <- ScanResult{Path: path}
			}
			return nil
		})
	}()

	return ch, nil
}

func IsAudioFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	_, ok := audioExtensions[ext]
	return ok
}


