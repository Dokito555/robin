package model

import "time"

type (
	CreatePlaylistRequest struct {
		UserID      int    `json:"user_id" validate:"required"`
		Name        string `json:"name" validate:"required"`
		Description string `json:"description" validate:"required"`
	}

	GetPlaylistRequest struct {
		ID int `json:"id" validate:"required"`
	}

	UpdatePlaylistRequest struct {
		ID          int    `json:"id" validate:"required"`
		Name        string `json:"user_id" validate:"required"`
		Description string `json:"description" validate:"required"`
	}

	DeletePlaylistRequest struct {
		ID int `json:"id" validate:"required"`
	}
)

type (
	PlaylistResponse struct {
		ID          int       `json:"id,omitempty"`
		UserID      int       `json:"user_id,omitempty"`
		Name        string    `json:"name,omitempty"`
		Description string    `json:"description,omitempty"`
		CreatedAt   time.Time `json:"created_at,omitempty"`
		UpdatedAt   time.Time `json:"updated_at,omitempty"`
	}
)
