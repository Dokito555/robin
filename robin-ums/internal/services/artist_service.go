package services

import (
	"context"
	"errors"

	"github.com/Dokito555/robin-ums/constants"
	"github.com/Dokito555/robin-ums/internal/entity"
	"github.com/Dokito555/robin-ums/internal/model"
	"github.com/Dokito555/robin-ums/internal/model/converter"
	"github.com/Dokito555/robin-ums/internal/repository"
	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type ArtistService struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	ArtistRepository *repository.ArtistRepository
	TokenService   *TokenService
}

func NewArtistService(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,
	ArtistRepository *repository.ArtistRepository, tokenSvc *TokenService) *ArtistService {
	return &ArtistService{
		DB:             db,
		Log:            logger,
		Validate:       validate,
		ArtistRepository: ArtistRepository,
		TokenService:   tokenSvc,
	}
}

func (s *ArtistService) RegisterArtist(ctx context.Context, req *model.RegisterArtistRequest) (*model.ArtistResponse, error) {
	s.Log.Info("starting Register Artist function")
	s.Log.Infof("request received: %+v", req)

	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := s.Validate.Struct(req)
	if err != nil {
		s.Log.Warnf("invalid request body: %+v", err)
		return nil, errors.New(constants.BAD_REQUEST)
	}

	if req.Password == "" || req.Email == "" {
		return nil, errors.New("password or email is empty")
	}

	artist, err := s.ArtistRepository.FindByEmail(s.DB, &entity.Artist{}, req.Email)
	if err != nil {
		s.Log.Warnf("database error fetching user: %+v", err)
		return nil, errors.New(constants.NOT_FOUND)
	}

	if artist != nil {
		s.Log.Warnf("user with that email already exists")
		return nil, errors.New("user already exists")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.Log.Warnf("failed to generate bcrypt hash: %+v", err)
		return nil, errors.New(constants.INTERNAL_SERVER_ERROR)
	}

	artist = &entity.Artist{
		Password: string(password),
		Email:    req.Email,
		UserName: req.UserName,
		Role:     req.Role,
	}

	if err := s.ArtistRepository.Create(s.DB, artist); err != nil {
		s.Log.Warnf("failed to create user in database: %+v", err)
		return nil, errors.New(constants.INTERNAL_SERVER_ERROR)
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("failed to commit transaction: %+v", err)
		return nil, errors.New(constants.INTERNAL_SERVER_ERROR)
	}

	return converter.ArtistToReponse(artist), nil
}

func (s *ArtistService) LoginArtist(ctx context.Context, req *model.LoginArtistRequest) (*model.ArtistResponse, error) {
	s.Log.Info("starting Login Artist function")
	s.Log.Infof("request received: %+v", req)

	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := s.Validate.Struct(req)
	if err != nil {
		s.Log.Warnf("invalid request body: %+v", err)
		return nil, errors.New(constants.BAD_REQUEST)
	}

	if req.Password == "" || req.Email == "" {
		return nil, errors.New("password or email is empty")
	}

	newArtist := new(entity.Artist)
	artist, err := s.ArtistRepository.FindByEmail(s.DB, newArtist, req.Email)
	if err != nil {
		s.Log.Warnf("database error fetching artist: %+v", err)
		return nil, errors.New(constants.INTERNAL_SERVER_ERROR)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(artist.Password), []byte(req.Password)); err != nil {
		s.Log.Warnf("failed to compare hashed password: %+v", err)
		return nil, errors.New(constants.INTERNAL_SERVER_ERROR)
	}

	token, err := s.TokenService.GenerateToken(ctx, artist.ID, constants.TOKEN_TYPE_TOKEN, artist.Email, artist.Role)
	if err != nil {
		s.Log.Warnf("failed to genereate token: %+v", err)
		return nil, err
	}

	refreshToken, err := s.TokenService.GenerateToken(ctx, artist.ID, constants.TOKEN_TYPE_REFRESH, artist.Email, artist.Role)
	if err != nil {
		s.Log.Warnf("failed to genereate refresh token: %+v", err)
		return nil, errors.New(constants.INTERNAL_SERVER_ERROR)
	}

	newArtist.Token = token
	newArtist.RefreshToken = refreshToken

	err = s.ArtistRepository.Update(s.DB, newArtist)
	if err != nil {
		s.Log.Warnf("failed to update user token and refresh token: %+v", err)
		return nil, errors.New(constants.INTERNAL_SERVER_ERROR)
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("failed to commit transaction: %+v", err)
		return nil, errors.New(constants.INTERNAL_SERVER_ERROR)
	}

	return converter.ArtistToReponse(newArtist), nil
}

func (s *ArtistService) UpdateArtist(ctx context.Context, req *model.UpdateArtistRequest) (*model.ArtistResponse, error) {
	
	return nil, nil
}

func (s *ArtistService) LogoutArtist(ctx context.Context, req *model.LogoutArtistRequest) (error) {
	return nil
}

func (s *ArtistService) GetArtist(ctx context.Context, req *model.GetArtistRequest) (*model.ArtistResponse, error) {
	return nil, nil
}

func (s *ArtistService) DeleteArtist(ctx context.Context, req *model.DeleteArtistRequest) (error) {
	return nil
}