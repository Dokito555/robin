package library

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
)

type CoverCache struct {
	dir string
}

func NewCoverCache(dir string) (*CoverCache, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create cover cache dir: %w", err)
	}
	return &CoverCache{dir: dir}, nil
}

func (c *CoverCache) Save(data []byte) (string, error) {
	if len(data) == 0 {
		return "", nil
	}

	hash := sha256.Sum256(data)
	// use first 16 hex chars, collision risk negligable for cover art
	name := fmt.Sprintf("%x.jpg", hash[:8])
	dest := filepath.Join(c.dir, name)

	if _, err := os.Stat(dest); err == nil {
		return dest, nil
	}

	if err := os.WriteFile(dest, data, 0644); err != nil {
		return "", fmt.Errorf("write cover: %w", err) 
	}

	return dest, nil
}

func (c *CoverCache) Delete(path string) error {
	if path == "" {
		return nil
	}

	err := os.Remove(path)
	if os.IsNotExist(err) {
		return nil
	}

	return err
}

// removes cover files that are no longer referenced by any track
func (c *CoverCache) Prune(activePaths map[string]struct{}) error {
	entries, err := os.ReadDir(c.dir)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		full := filepath.Join(c.dir, e.Name())
		// TODO: could be better
		if _, active := activePaths[full]; !active {
			os.Remove(full)
		}
	}

	return nil
}