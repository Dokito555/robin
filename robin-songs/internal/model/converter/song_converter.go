package converter

import (
	"github.com/Dokito555/robin-songs/internal/entity"
	"github.com/Dokito555/robin-songs/internal/model"
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
