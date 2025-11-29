package converter

import (
	"github.com/Dokito555/robin/robin-catalog/internal/entity"
	"github.com/Dokito555/robin/robin-catalog/internal/model"
)


func ArtistToReponse(artist *entity.Artist) *model.ArtistResponse {
	return &model.ArtistResponse{
		ID:           artist.ID,
		Email:        artist.Email,
		UserName:     artist.UserName,
		Role:         artist.Role,
		Bio:          artist.Bio,
		Genre:        artist.Genre,
		Token:        artist.Token,
		RefreshToken: artist.RefreshToken,
		CreatedAt:    artist.CreatedAt,
		UpdatedAt:    artist.UpdatedAt,
	}
}
