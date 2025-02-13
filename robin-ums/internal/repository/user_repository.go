package repository

import (
	"errors"

	"github.com/Dokito555/robin-ums/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[entity.User]
	Log *logrus.Logger
}

func NewUserRepository(log *logrus.Logger, db *gorm.DB) *UserRepository {
	return &UserRepository{
		Repository: Repository[entity.User]{DB: db},
		Log:        log,
	}
}

func (r *UserRepository) FindByToken(db *gorm.DB, user *entity.User, token string) (*entity.User, error) {
	err := db.Where("token = ?", token).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByEmail(db *gorm.DB, user *entity.User, email string) (*entity.User, error) {
	err := db.Where("email = ?", email).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
