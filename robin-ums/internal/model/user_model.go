package model

import "time"

type (
	LoginUserRequest struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	RegisterUserRequest struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
		Username string `json:"username" validate:"username"`
		Role     string `json:"role" validate:"required"`
	}

	GetUserRequest struct {
		Id string `json:"id" validate:"required"`
	}

	UpdateUserRequest struct {
		Password string `json:"password" validate:"required"`
		Username string `json:"username" validate:"username"`
	}

	LogoutUserRequest struct {
		Token string `json:"token" validate:"required"`
	}

	DeleteUserRequest struct {
		Id string `json:"id" validate:"required"`
	}

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
