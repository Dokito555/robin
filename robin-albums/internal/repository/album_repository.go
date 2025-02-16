package repository

import (
	"github.com/Dokito555/robin-albums/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AlbumRepository struct {
	Repository[entity.Album]
	Log *logrus.Logger
}

func NewAlbumRepository(log *logrus.Logger, db *gorm.DB) *AlbumRepository {
	return &AlbumRepository{
		Repository: Repository[entity.Album]{DB: db},
		Log:        log,
	}
}