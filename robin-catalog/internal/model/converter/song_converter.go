package converter

import (
	"github.com/Dokito555/robin/robin-catalog/internal/entity"
	"github.com/Dokito555/robin/robin-catalog/internal/model"
)


func SongToResponse(s *entity.Song) *model.SongResponse {
	return &model.SongResponse{
		ID:        s.ID,
		ArtistID:  s.ArtistID,
		Name:      s.Name,
		Duration:  s.Duration,
		Link:      s.Link,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}
