package services

import (
	"context"
	"errors"

	"github.com/Dokito555/robin-ums/constants"
	"github.com/Dokito555/robin-ums/internal/entity"
	model "github.com/Dokito555/robin-ums/internal/models"
	"github.com/Dokito555/robin-ums/internal/models/converter"
	"github.com/Dokito555/robin-ums/internal/repository"
	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserService struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *repository.UserRepository
	TokenService   *TokenService
}

func NewUserService(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,
	userRepository *repository.UserRepository, tokenSvc *TokenService) *UserService {
	return &UserService{
		DB:             db,
		Log:            logger,
		Validate:       validate,
		UserRepository: userRepository,
		TokenService: 	tokenSvc,
	}
}

func (s *UserService) Verify(ctx context.Context, req *model.VerifyUserRequest) (*model.UserResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := s.Validate.Struct(req)
	if err != nil {
		s.Log.Warnf("invalid request body : %+v", err)
		return nil, errors.New(constants.BAD_REQUEST)
	}

	user := new(entity.User)
	rsp, err := s.UserRepository.FindByToken(tx, user, req.Token)
	if err != nil {
		s.Log.Warnf("failed find user by token : %+v", err)
		return nil, errors.New(constants.INTERNAL_SERVER_ERROR)
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("failed commit transaction : %+v", err)
		return nil, errors.New(constants.INTERNAL_SERVER_ERROR)
	}

	return converter.UserToResponse(rsp), nil
}