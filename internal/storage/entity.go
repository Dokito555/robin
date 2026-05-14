package storage

import "time"

type Track struct {
	ID        int64
	Title     string
	Artist    string
	Album     string
	Genre     string
	Duration  time.Duration
	Path      string
	CoverPath string
	FolderID  int64
	CreatedAt time.Time
}

type Folder struct {
	ID        int64
	Path      string
	Type      string
	CreatedAt time.Time
}

type Playlist struct {
	ID        int64
	Name      string
	CreatedAt time.Time
}
