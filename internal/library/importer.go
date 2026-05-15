package library

import (
	"fmt"
	"log"
	"sync"

	"github.com/Dokito555/robin/internal/events"
	"github.com/Dokito555/robin/internal/storage"
)

type Importer struct {
	db          *storage.DB
	covers      *CoverCache
	bus         *events.Bus
	workerCount int
}

func NewImporter(db *storage.DB, covers *CoverCache, bus *events.Bus) *Importer {
	return &Importer{
		db:          db,
		covers:      covers,
		bus:         bus,
		workerCount: 4,
	}
}

func (imp *Importer) Import(folderID int64, rootPath string) {
	go func() {
		imp.bus.Publish(events.ScanStarted, events.FolderPayload {
			FolderID: folderID,
			Path: rootPath,
		})

		results, err := Scan(rootPath)
		if err != nil {
			log.Printf("importer: scan error %s: %v", rootPath, err)
			return
		}

		// collect path first to know total for progress reporting
		var paths []string
		for r := range results {
			paths = append(paths, r.Path)
		}

		total := len(paths)
		var (
			mu sync.Mutex
			added int
			wg sync.WaitGroup
		)

		sem := make(chan struct{}, imp.workerCount)

		for i, path := range paths {
			wg.Add(1)
			sem <- struct{}{}

			go func (idx int, p string)  {
				defer wg.Done()
				defer func() {<-sem}()

				if err := imp.importOne(folderID, p); err != nil {
					log.Printf("importer: skip %s: %v", p, err)
					imp.db.LogScanError(path, err.Error())
					return
				}

				mu.Lock()
				added++
				current := added
				mu.Unlock()

				imp.bus.Publish(events.ScanProgress, events.ScanProgressPayload{
					FolderID: folderID,
					Current: current,
					Total: total,
				})
			}(i, path)
		}

		wg.Wait()

		imp.bus.Publish(events.ScanDone, events.ScanDonePayload{
			FolderID: folderID,
			TracksAdded: added,
		})

		log.Printf("importer: done folder=%d added=%d total=%d", folderID, added, total)
	}()
}

func (imp *Importer) ImportOne(folderID int64, path string) error {
	return imp.importOne(folderID, path)
}

func (imp *Importer) importOne(folderID int64, path string) error {
	meta, err := ReadMeta(path)
	if err != nil {
		return fmt.Errorf("read meta: %w", err)
	}

	if err := meta.validate(); err != nil {
		return fmt.Errorf("invalid metadata: %w", err)
	}

	coverPath, err := imp.covers.Save(meta.CoverData)
	if err != nil {
		log.Printf("importer: cover save failed for %s: %v", path, err)
	}

	track := storage.Track{
		Title:     meta.Title,
		Artist:    meta.Artist,
		Album:     meta.Album,
		Genre:     meta.Genre,
		Duration:  meta.Duration,
		Path:      meta.Path,
		CoverPath: coverPath,
		FolderID:  folderID,
	}

	id, err := imp.db.UpsertTrack(track)
	if err != nil {
		return fmt.Errorf("upser track: %w",err)
	}

	imp.bus.Publish(events.TrackAdded, events.TrackAddedPayload{
		TrackID: id,
		Path: path,
	})

	return nil
}