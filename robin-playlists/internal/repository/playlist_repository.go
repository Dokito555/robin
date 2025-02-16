package repository

import (
	"github.com/Dokito555/robin-playlists/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PlaylistRepository struct {
	Repository[entity.Playlist]
	Log *logrus.Logger
}

func NewPlaylistRepository(log *logrus.Logger, db *gorm.DB) *PlaylistRepository {
	return &PlaylistRepository{
		Repository: Repository[entity.Playlist]{DB: db},
		Log:        log,
	}
}
