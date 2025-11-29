package model

import "time"

type (
	CreateAlbumRequest struct {
		ArtistID int    `json:"artist_id" validate:"required"`
		Name     string `json:"name" validate:"required"`
		Type     string `json:"type" validate:"required"`
	}

	UpdateAlbumRequest struct {
		ID   int    `json:"id" validate:"required"`
		Name string `json:"name" validate:"required"`
		Type string `json:"type" validate:"required"`
	}

	DeleteAlbumRequest struct {
		ID int `json:"id" validate:"required"`
	}

	GetAlbumRequest struct {
		ID int `json:"id" validate:"required"`
	}
)

type (
	AlbumResponse struct {
		ID        int       `json:"id,omitempty"`
		Name      string    `json:"name,omitempty"`
		ArtistID  int       `json:"artist_id,omitempty"`
		Type      string    `json:"type,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty"`
		UpdatedAt time.Time `json:"updated_at,omitempty"`
	}
)
