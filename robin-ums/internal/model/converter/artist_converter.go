package converter

import (
	"github.com/Dokito555/robin-ums/internal/entity"
	model "github.com/Dokito555/robin-ums/internal/model"
)

func ArtistToReponse(artist *entity.Artist) *model.UserResponse {
	return &model.UserResponse{
		ID:           artist.ID,
		Email:        artist.Email,
		UserName:     artist.UserName,
		Role:         artist.Role,
		Token:        artist.Token,
		RefreshToken: artist.RefreshToken,
		CreatedAt:    artist.CreatedAt,
		UpdatedAt:    artist.UpdatedAt,
	}
}
