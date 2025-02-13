package repository

import (
	"errors"

	"github.com/Dokito555/robin-ums/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ArtistRepository struct {
	Repository[entity.Artist]
	Log *logrus.Logger
}

func NewAristRepository(log *logrus.Logger, db *gorm.DB) *ArtistRepository {
	return &ArtistRepository{
		Repository: Repository[entity.Artist]{DB: db},
		Log:        log,
	}
}

func (r *ArtistRepository) FindByToken(db *gorm.DB, artist *entity.Artist, token string) (*entity.Artist, error) {
	err := db.Where("token = ?", token).First(artist).Error
	if err != nil {
		return nil, err
	}
	return artist, nil
}

func (r *ArtistRepository) FindByEmail(db *gorm.DB, user *entity.Artist, email string) (*entity.Artist, error) {
	err := db.Where("email = ?", email).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
