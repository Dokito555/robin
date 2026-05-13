package main

import (
	"fmt"
	"time"

	"github.com/Dokito555/robin/internal/events"
	"github.com/Dokito555/robin/internal/library"
	"github.com/Dokito555/robin/internal/storage"
)

func main() {
	db, _ := storage.Open("./dev.db")
	bus := events.NewBus()
	cache, _ := library.NewCoverCache("./cache/covers")
	lib, _ := library.New(db, cache, bus)
	defer lib.Close()

	bus.Subscribe(events.ScanDone, func(e events.Event) {
		p := e.Payload.(events.ScanDonePayload)
		fmt.Printf("scan done — %d tracks\n", p.TracksAdded)

		tracks, _ := lib.AllTracks()
		for _, t := range tracks {
			fmt.Printf("  %s — %s (%s)\n", t.Artist, t.Title, t.Album)
		}
	})

	bus.Subscribe(events.TrackAdded, func(e events.Event) {
		p := e.Payload.(events.TrackAddedPayload)
		fmt.Printf("live: new file detected %s\n", p.Path)
	})

	lib.AddFolder("C:/Users/Zwacht/OneDrive/Documents/friendsure/music", "local")

	// drop a new mp3 into the folder while this runs — watcher should pick it up
	time.Sleep(60 * time.Second)
}
