package model

import "time"

type (
	LoginRequest struct{}

	RegisterRequest struct{}

	GetUserRequest struct{}

	UpdateUserRequest struct{}

	DeleteUserRequest struct{}

	VerifyUserRequest struct {
		Token string `json:"token"`
	}
)

type (
	UserResponse struct {
		ID           int       `json:"id,omitempty"`
		Email        string    `json:"email,omitempty"`
		UserName     string    `json:"username,omitempty"`
		Role         string    `json:"role,omitempty"`
		Token        string    `json:"token,omitempty"`
		RefreshToken string    `json:"refresh_token,omitempty"`
		CreatedAt    time.Time `json:"created_at,omitempty"`
		UpdatedAt    time.Time `json:"updated_at,omitempty"`
	}
)

