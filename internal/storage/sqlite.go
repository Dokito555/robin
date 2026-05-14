package storage

import (
	"time"
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
)

type DB struct {
	Conn *sql.DB
}

func Open(dsn string) (*DB, error) {
	conn, err := sql.Open("sqlite", dsn+"?_journal_mode=WAL&_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	conn.SetMaxOpenConns(1)
	db := &DB{Conn: conn}
	// ?
	if err := db.migrate(); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return db, nil
}

func (db *DB) Close() error {
	return db.Conn.Close()
}

func (db *DB) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS folders (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		path       TEXT    NOT NULL UNIQUE,
		type       TEXT    NOT NULL DEFAULT 'local',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS tracks (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		title      TEXT,
		artist     TEXT,
		album      TEXT,
		genre      TEXT,
		duration   INTEGER, -- stored as milliseconds
		path       TEXT    NOT NULL UNIQUE,
		cover_path TEXT,
		folder_id  INTEGER REFERENCES folders(id) ON DELETE CASCADE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_tracks_artist  ON tracks(artist);
	CREATE INDEX IF NOT EXISTS idx_tracks_album   ON tracks(album);
	CREATE INDEX IF NOT EXISTS idx_tracks_folder  ON tracks(folder_id);

	CREATE TABLE IF NOT EXISTS playlists (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		name       TEXT NOT NULL UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS playlist_tracks (
		playlist_id INTEGER NOT NULL REFERENCES playlists(id) ON DELETE CASCADE,
		track_id    INTEGER NOT NULL REFERENCES tracks(id)    ON DELETE CASCADE,
		position    INTEGER NOT NULL,
		PRIMARY KEY (playlist_id, track_id)
	);`

	_, err := db.Conn.Exec(schema)
	return err
}

func (db *DB) InsertFolder(path, folderType string) (int64, error) {
	var id int64
	err := db.Conn.QueryRow(`
		INSERT INTO folders (path, type)
		VALUES (?, ?)
		ON CONFLICT(path) DO UPDATE SET type = excluded.type
		RETURNING id
	`, path, folderType).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert folder: %w", err)
	}
	return id, nil
}

func (db *DB) GetFolders() ([]Folder, error) {
	rows, err := db.Conn.Query(`SELECT id, path, type, created_at FROM folders`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var folders []Folder
	for rows.Next() {
		var f Folder
		if err := rows.Scan(&f.ID, &f.Path, &f.Type, &f.CreatedAt); err != nil {
			return nil, err
		}
		folders = append(folders, f)
	}
	return folders, rows.Err()
}

func (db *DB) DeleteFolder(id int64) error {
	_, err := db.Conn.Exec(`DELETE FROM folders WHERE id = ?`, id)
	return err
}

func (db *DB) UpsertTrack(t Track) (int64, error) {
	res, err := db.Conn.Exec(`
		INSERT INTO tracks (title, artist, album, genre, duration, path, cover_path, folder_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(path) DO UPDATE SET
			title      = excluded.title,
			artist     = excluded.artist,
			album      = excluded.album,
			genre      = excluded.genre,
			duration   = excluded.duration,
			cover_path = excluded.cover_path,
			folder_id  = excluded.folder_id`,
		t.Title, t.Artist, t.Album, t.Genre,
		t.Duration.Milliseconds(), t.Path, t.CoverPath, t.FolderID,
	)

	if err != nil {
		return 0, fmt.Errorf("upsert track: %w", err)
	}

	return res.LastInsertId()
}

func (db *DB) GetAllTracks() ([]Track, error) {
	return db.queryTracks(`SELECT id, title, artist, album, genre, duration, path, cover_path, folder_id, created_at FROM tracks ORDER BY artist, album, title`)
}

func (db *DB) SearchTracks(query string) ([]Track, error) {
	like := "%" + query + "%"
	return db.queryTracks(`
		SELECT id, title, artist, album, genre, duration, path, cover_path, folder_id, created_at
		FROM tracks
		WHERE title LIKE ? OR artist LIKE ? OR album LIKE ?
		ORDER BY artist, album, title`,
		like, like, like,
	)
}

func (db *DB) GetTracksByFolder(folderID int64) ([]Track, error) {
	return db.queryTracks(
		`SELECT id, title, artist, album, genre, duration, path, cover_path, folder_id, created_at FROM tracks WHERE folder_id = ? ORDER BY title`,
		folderID,
	)
}

func (db *DB) DeleteTrackByPath(path string) error {
	_, err := db.Conn.Exec(`DELETE FROM tracks WHERE path = ?`, path)
	return err
}

func (db *DB) queryTracks(query string, args ...any) ([]Track, error) {
	rows, err := db.Conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []Track
	for rows.Next() {
		var t Track
		var durMS int64
		if err := rows.Scan(
			&t.ID, &t.Title, &t.Artist, &t.Album, &t.Genre,
			&durMS, &t.Path, &t.CoverPath, &t.FolderID, &t.CreatedAt,
		); err != nil {
			return nil, err
		}
		t.Duration = time.Duration(durMS) * time.Millisecond
		tracks = append(tracks, t)
	}
	return tracks, rows.Err()
}

func (db *DB) CreatePlaylist(name string) (int64, error) {
	res, err := db.Conn.Exec(`INSERT INTO playlists (name) VALUES (?)`, name)
	if err != nil {
		return 0, fmt.Errorf("create playlist: %w", err)
	}
	return res.LastInsertId()
}

func (db *DB) GetPlaylist() ([]Playlist, error) {
	rows, err := db.Conn.Query(`SELECT id, name, created_at FROM playlists`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var playlists []Playlist
	for rows.Next() {
		var p Playlist
		if err := rows.Scan(&p.ID, &p.Name, &p.CreatedAt); err != nil {
			return nil, err
		}
		playlists = append(playlists, p)
	}
	return playlists, rows.Err()
}

func (db *DB) AddTrackToPlaylist(playlistID, trackID int64, position int) error {
	_, err := db.Conn.Exec(
		`INSERT INTO playlist_tracks (playlist_id, track_id, position) VALUES (?, ?, ?)
		 ON CONFLICT DO NOTHING`,
		playlistID, trackID, position,
	)
	return err
}

func (db *DB) GetPlaylistTracks(playlistID int64) ([]Track, error) {
	return db.queryTracks(`
		SELECT t.id, t.title, t.artist, t.album, t.genre, t.duration, t.path, t.cover_path, t.folder_id, t.created_at
		FROM tracks t
		JOIN playlist_tracks pt ON pt.track_id = t.id
		WHERE pt.playlist_id = ?
		ORDER BY pt.position`,
		playlistID,
	)
}

func (db *DB) DeletePlaylist(id int64) error {
	_, err := db.Conn.Exec(`DELETE FROM playlists WHERE id = ?`, id)
	return err
}