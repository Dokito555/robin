package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Dokito555/robin/internal/events"
	"github.com/Dokito555/robin/internal/library"
	"github.com/Dokito555/robin/internal/storage"
)

func main() {
	db, err := storage.Open("robin.db")
	if err != nil {
		log.Fatal(err)
	}

	bus := events.NewBus()

	cache, err := library.NewCoverCache("./cache/covers")
	if err != nil {
		log.Fatal(err)
	}

	lib, err := library.New(db, cache, bus)
	if err != nil {
		log.Fatal(err)
	}
	
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

	lib.AddFolder("E:/Music", "local")

	// drop a new mp3 into the folder while this runs — watcher should pick it up
	time.Sleep(60 * time.Second)
}
