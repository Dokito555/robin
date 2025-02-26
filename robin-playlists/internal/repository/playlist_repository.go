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

func (r *PlaylistRepository) GetPlaylistList(db *gorm.DB, id int) ([]entity.Playlist, error) {
	var playlists []entity.Playlist
	if err := db.Where("user_id == ?", id).Find(&playlists).Error; err != nil {
		return nil, err
	}
	return playlists, nil
}