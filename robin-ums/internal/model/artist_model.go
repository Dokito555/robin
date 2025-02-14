package model

import "time"

type (
	RegisterArtistRequest struct {
		UserName string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
		Genre    string `json:"genre" validate:"required"`
		Bio      string `json:"bio" validate:"required"`
		Role     string `json:"role" validate:"required"`
	}

	LoginArtistRequest struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	GetArtistRequest struct {
		ID int `json:"id" validate:"required"`
	}

	DeleteArtistRequest struct {
		ID int `json:"id" validate:"required"`
	}

	LogoutArtistRequest struct {
		ID int `json:"id" validate:"required"`
	}

	UpdateArtistRequest struct {
		Password string `json:"password" validate:"required"`
		Genre    string `json:"genre" validate:"required"`
		Bio      string `json:"bio" validate:"required"`
	}
)

type (
	ArtistResponse struct {
		ID           int       `json:"id,omitempty"`
		Email        string    `json:"email,omitempty"`
		UserName     string    `json:"username,omitempty"`
		Genre        string    `json:"genre,omitempty"`
		Bio          string    `json:"bio,omitempty"`
		Role         string    `json:"role,omitempty"`
		Token        string    `json:"token,omitempty"`
		RefreshToken string    `json:"refresh_token,omitempty"`
		CreatedAt    time.Time `json:"created_at,omitempty"`
		UpdatedAt    time.Time `json:"updated_at,omitempty"`
	}
)
