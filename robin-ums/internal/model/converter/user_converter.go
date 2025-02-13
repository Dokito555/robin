package converter

import (
	"github.com/Dokito555/robin-ums/internal/entity"
	model "github.com/Dokito555/robin-ums/internal/model"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		ID:           user.ID,
		Email:        user.Email,
		UserName:     user.UserName,
		Role:         user.Role,
		Token:        user.Token,
		RefreshToken: user.RefreshToken,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}