package main

import (
	"fmt"
	"time"

	"github.com/Dokito555/robin/internal/events"
	"github.com/Dokito555/robin/internal/library"
	"github.com/Dokito555/robin/internal/storage"
)

func main() {
	// main.go — temporary smoke test, delete after Phase 5
	db, _ := storage.Open("./dev.db")
	bus := events.NewBus()
	cache, _ := library.NewCoverCache("./cache/covers")
	imp := library.NewImporter(db, cache, bus)

	folderID, _ := db.InsertFolder("C:/Users/Zwacht/OneDrive/Documents/friendsure/music", "local")

	bus.Subscribe(events.ScanDone, func(e events.Event) {
		p := e.Payload.(events.ScanDonePayload)
		fmt.Printf("scan done — %d tracks added\n", p.TracksAdded)
	})

	imp.Import(folderID, "C:/Users/Zwacht/OneDrive/Documents/friendsure/music")
	time.Sleep(30 * time.Second) // wait for scan
}
