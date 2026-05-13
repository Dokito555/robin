package library

import (
	"log"
	"path/filepath"
	"sync"

	"github.com/Dokito555/robin/internal/events"
	"github.com/Dokito555/robin/internal/storage"
	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	fsw      *fsnotify.Watcher
	db       *storage.DB
	imp      *Importer
	bus      *events.Bus
	folders  map[string]int64
	mu       sync.RWMutex
	stopOnce sync.Once
	done     chan struct{}
}

func NewWatcher(db *storage.DB, imp *Importer, bus *events.Bus) (*Watcher, error) {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	w := &Watcher{
		fsw: fsw,
		db: db,
		imp: imp,
		bus: bus,
		folders: make(map[string]int64),
		done: make(chan struct{}),
	}
	go w.loop()
	return w, nil
}

func (w *Watcher) Watch(folderID int64, path string) error {
	if err := w.fsw.Add(path); err != nil {
		return err
	}

	w.mu.Lock()
	w.folders[path] = folderID
	w.mu.Unlock()
	log.Printf("watcher: watching folder=%d path=%s", folderID, path)
	return nil
}

func (w *Watcher) Unwatch(path string) {
	w.fsw.Remove(path)
	w.mu.Lock()
	delete(w.folders, path)
	w.mu.Unlock()
	log.Printf("watcher: unwatched path=%s", path)
}

func (w *Watcher) Close() {
	w.stopOnce.Do(func() {
		close(w.done)
		w.fsw.Close()
	})
}

func (w *Watcher) loop() {
	for {
		select {
		case <-w.done:
			return
			
		case event, ok := <-w.fsw.Events:
			if !ok {
				return
			}
			w.handleEvent(event)
		
		case err, ok := <-w.fsw.Errors:
			if !ok {
				return
			}
			log.Printf("watcher: fsnotify error: %v", err)
		}
	}
}

func (w *Watcher) handleEvent(event fsnotify.Event) {
	path := event.Name

	if !IsAudioFile(path) {
		return
	}

	folderID := w.folderIDFor(path)
	if folderID == 0 {
		return
	}

	switch {
	case event.Has(fsnotify.Create):
		log.Printf("watcher: new file detected: %s", path)
		if err := w.imp.ImportOne(folderID, path); err != nil {
			log.Printf("wather: import failed %s: %v", path, err)
		}
		
	case event.Has(fsnotify.Remove), event.Has(fsnotify.Rename):
		log.Printf("watcher: file removed: %s", path)
		if err := w.db.DeleteTrackByPath(path); err != nil {
			log.Printf("watcher: delete failed %s: %v", path, err)
		}
		w.bus.Publish(events.TrackRemoved, events.TrackRemovedPayload{
			Path: path,
		})

	case event.Has(fsnotify.Write):
		log.Printf("watcher: file modifiedL %s", path)
		if err := w.imp.ImportOne(folderID, path); err != nil {
			log.Printf("watcher: re-import failed %s: %v", path, err)
		}
	}
}

func (w *Watcher) folderIDFor(path string) int64 {
	w.mu.RLock()
	defer w.mu.RUnlock()

	dir := filepath.Dir(path)
	for watchedPath, id := range w.folders {
		rel, err := filepath.Rel(watchedPath, dir)
		if err == nil && rel != ".." {
			return id
		}
	}

	return 0
}

