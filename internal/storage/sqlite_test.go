package storage

import (
	"testing"
	"time"
)

func TestStorage(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer db.Close()

	// Insert folder
	folderID, err := db.InsertFolder("/home/user/Music", "local")
	if err != nil {
		t.Fatalf("insert folder: %v", err)
	}

	// Upsert track
	trackID, err := db.UpsertTrack(Track{
		Title:    "Numb",
		Artist:   "Linkin Park",
		Album:    "Meteora",
		Genre:    "Rock",
		Duration: 3*time.Minute + 5*time.Second,
		Path:     "/home/user/Music/Linkin Park/Numb.mp3",
		FolderID: folderID,
	})
	if err != nil {
		t.Fatalf("upsert track: %v", err)
	}
	t.Logf("inserted track id=%d", trackID)

	// Get all tracks
	tracks, err := db.GetAllTracks()
	if err != nil {
		t.Fatalf("get tracks: %v", err)
	}
	if len(tracks) != 1 {
		t.Fatalf("expected 1 track, got %d", len(tracks))
	}
	if tracks[0].Artist != "Linkin Park" {
		t.Errorf("expected Linkin Park, got %s", tracks[0].Artist)
	}

	// Search
	results, err := db.SearchTracks("numb")
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 search result, got %d", len(results))
	}

	// Playlist
	plID, _ := db.CreatePlaylist("My Mix")
	db.AddTrackToPlaylist(plID, tracks[0].ID, 0)
	plTracks, _ := db.GetPlaylistTracks(plID)
	if len(plTracks) != 1 {
		t.Errorf("expected 1 playlist track, got %d", len(plTracks))
	}

	// Delete folder cascades to tracks
	db.DeleteFolder(folderID)
	tracks, _ = db.GetAllTracks()
	if len(tracks) != 0 {
		t.Errorf("expected cascade delete, still got %d tracks", len(tracks))
	}
}