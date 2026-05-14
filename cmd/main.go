package main

import (
	"fmt"
	"time"

	"github.com/Dokito555/robin/internal/events"
	"github.com/Dokito555/robin/internal/library"
	"github.com/Dokito555/robin/internal/player"
	"github.com/Dokito555/robin/internal/storage"
)

func main() {
	// db, err := storage.Open("robin.db")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// bus := events.NewBus()

	// cache, err := library.NewCoverCache("./cache/covers")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// lib, err := library.NewLibrary(db, cache, bus)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer lib.Close()

	// bus.Subscribe(events.ScanDone, func(e events.Event) {
	// 	p := e.Payload.(events.ScanDonePayload)
	// 	fmt.Printf("scan done — %d tracks\n", p.TracksAdded)

	// 	tracks, _ := lib.AllTracks()
	// 	for _, t := range tracks {
	// 		fmt.Printf("  %s — %s (%s)\n", t.Artist, t.Title, t.Album)
	// 	}
	// })

	// bus.Subscribe(events.TrackAdded, func(e events.Event) {
	// 	p := e.Payload.(events.TrackAddedPayload)
	// 	fmt.Printf("live: new file detected %s\n", p.Path)
	// })

	// lib.AddFolder("E:/Music", "local")

	// // drop a new mp3 into the folder while this runs — watcher should pick it up
	// time.Sleep(60 * time.Second)

	db, _ := storage.Open("./robin.db")
	bus := events.NewBus()
	cache, _ := library.NewCoverCache("./cache/covers")
	lib, _ := library.NewLibrary(db, cache, bus)
	defer lib.Close()

	p, _ := player.NewPlayer(bus)

	bus.Subscribe(events.TrackStarted, func(e events.Event) {
		pay := e.Payload.(events.TrackStartedPayload)
		fmt.Printf("▶  now playing: %s — %s\n", pay.Artist, pay.Title)
	})

	bus.Subscribe(events.TrackEnded, func(e events.Event) {
		fmt.Println("⏭  track ended — advancing queue")
	})

	// load library into queue and play
	tracks, _ := lib.AllTracks()
	if len(tracks) > 0 {
		p.Queue.SetTracks(tracks)
		p.PlayCurrent()
	}

	// let it play, test controls
	time.Sleep(15 * time.Second)
	fmt.Println("pausing...")
	p.Pause()

	time.Sleep(2 * time.Second)
	fmt.Println("resuming...")
	p.Resume()

	time.Sleep(5 * time.Second)
	fmt.Println("skipping...")
	p.Next()

	time.Sleep(10 * time.Second)
	p.Stop()
}
