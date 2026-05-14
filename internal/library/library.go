package library

import (
	"fmt"
	"log"

	"github.com/Dokito555/robin/internal/events"
	"github.com/Dokito555/robin/internal/storage"
)

type Library struct {
	db      *storage.DB
	covers  *CoverCache
	imp     *Importer
	watcher *Watcher
	bus     *events.Bus
}

// TODO: support massive libraries
func New(db *storage.DB, covers *CoverCache, bus *events.Bus) (*Library, error) {
	imp := NewImporter(db, covers, bus)

	watcher, err := NewWatcher(db, imp, bus)
	if err != nil {
		return nil, fmt.Errorf("create watcher: %w", err)
	}

	lib := &Library{
		db: db,
		covers: covers,
		imp: imp,
		watcher: watcher,
		bus: bus,
	}

	if err := lib.resumeFolders(); err != nil {
		log.Printf("library: resume folders error: %v", err)
	}
	
	return lib, nil
}

// registers a new root folder, triggers a scan, and starts watching
func (lib *Library) AddFolder(path, folderType string) error {
	folderID, err := lib.db.InsertFolder(path, folderType)
	if err != nil {
		return fmt.Errorf("insert folder: %w", err)
	}

	lib.bus.Publish(events.FolderAdded, events.FolderPayload{
		FolderID: folderID,
		Path: path,
	})

	if err := lib.watcher.Watch(folderID, path); err != nil {
		log.Printf("library: watch failed for %s: %v", path, err)
	}

	lib.imp.Import(folderID, path)

	return nil
}

func (lib *Library) RemoveFolder(folderID int64) error {
	folders, err := lib.db.GetFolders()
	if err != nil {
		return err
	}

	var path string
	for _, f := range folders {
		if f.ID == folderID {
			path = f.Path
			break
		}
	}

	lib.watcher.Unwatch(path)

	if err := lib.db.DeleteFolder(folderID); err != nil {
		return fmt.Errorf("delete folder: %w", err)
	}

	lib.bus.Publish(events.FolderRemoved, events.FolderPayload{
		FolderID: folderID,
		Path: path,
	})

	return nil
}

func (lib *Library) RescanFolder(folderID int64, path string) {
	lib.imp.Import(folderID, path)
}

func (lib *Library) AllTracks() ([]storage.Track, error) {
	return lib.db.GetAllTracks()
}

func (lib *Library) Search(query string) ([]storage.Track, error) {
	return lib.db.SearchTracks(query)
}

func (lib *Library) TracksIn(folderID int64) ([]storage.Track, error) {
	return lib.db.GetTracksByFolder(folderID)
}

func (lib *Library) Folder() ([]storage.Folder, error) {
	return lib.db.GetFolders()
}

func (lib *Library) Close() {
	lib.watcher.Close()
}

// re-attaches watchers for all folders saved in the DB
func (lib *Library) resumeFolders() error {
	folders, err := lib.db.GetFolders()
	if err != nil {
		return err
	}

	for _, f := range folders {
		if err := lib.watcher.Watch(f.ID, f.Path); err != nil {
			log.Printf("library: resume watch failed %s: %v", f.Path, err)
		}
	}

	log.Printf("library: resumed %d folder watchers", len(folders))
	return nil
}
