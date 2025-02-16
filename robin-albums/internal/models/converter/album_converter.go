package converter

import (
	"github.com/Dokito555/robin-albums/internal/entity"
	model "github.com/Dokito555/robin-albums/internal/models"
)

func AlbumToResponse(album *entity.Album) *model.AlbumResponse {
	return &model.AlbumResponse{
		ID:        album.ID,
		Name:      album.Name,
		ArtistID:  album.ArtistID,
		Type:      album.Type,
		CreatedAt: album.CreatedAt,
		UpdatedAt: album.UpdatedAt,
	}
}
