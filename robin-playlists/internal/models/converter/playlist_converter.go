package converter

import (
	"github.com/Dokito555/robin-playlists/internal/entity"
	model "github.com/Dokito555/robin-playlists/internal/models"
)

func PlaylistToResponse(p *entity.Playlist) *model.PlaylistResponse {
	return &model.PlaylistResponse{
		ID: p.ID,
		UserID: p.ID,
		Name: p.Name,
		Description: p.Description,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}