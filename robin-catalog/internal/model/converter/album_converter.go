package converter

import (
	"github.com/Dokito555/robin/robin-catalog/internal/entity"
	"github.com/Dokito555/robin/robin-catalog/internal/model"
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
