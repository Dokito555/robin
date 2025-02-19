package model

import (
	"mime/multipart"
	"time"
)

type (
	CreateSongRequest struct {
		ArtistID int    `json:"artist_id" validate:"required"`
		AlbumID  int    `json:"album_id" validate:"required"`
		Name     string `json:"name" validate:"required"`
		Link     string `json:"link" validate:"required"`
		Duration int    `json:"duration" validate:"required"`
	}

	GetSongRequest struct {
		ID int `json:"id" validate:"required"`
	}

	UpdateSongRequest struct {
		ID       int    `json:"id" validate:"required"`
		Name     string `json:"name" validate:"required"`
		Link     string `json:"link" validate:"required"`
		Duration int    `json:"duration" validate:"required"`
	}

	DeleteSongRequest struct {
		ID int `json:"id" validate:"required"`
	}
)

type (
	File struct {
		File     multipart.File
		FileName string
	}
)

type (
	SongResponse struct {
		ID        int       `json:"id"`
		ArtistID  int       `json:"artist_id"`
		AlbumID   int       `json:"album_id"`
		Name      string    `json:"name"`
		Link      string    `json:"link"`
		Duration  int       `json:"duration"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)
