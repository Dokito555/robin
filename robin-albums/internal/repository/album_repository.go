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

func (r *AlbumRepository) AlbumListByArtistID(db *gorm.DB, id int) ([]entity.Album, error) {
	var albums []entity.Album
	if err := db.Where("artist_id == ?", id).Find(&albums).Error; err != nil {
		return nil, err
	}

	return albums, nil
}