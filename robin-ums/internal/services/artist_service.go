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
	DB               *gorm.DB
	Log              *logrus.Logger
	Validate         *validator.Validate
	ArtistRepository *repository.ArtistRepository
	TokenService     *TokenService
}

func NewArtistService(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,
	ArtistRepository *repository.ArtistRepository, tokenSvc *TokenService) *ArtistService {
	return &ArtistService{
		DB:               db,
		Log:              logger,
		Validate:         validate,
		ArtistRepository: ArtistRepository,
		TokenService:     tokenSvc,
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
		s.Log.Warnf("database error fetching artist: %+v", err)
		return nil, errors.New(constants.NOT_FOUND)
	}

	if artist != nil {
		s.Log.Warnf("artist with that email already exists")
		return nil, errors.New("artist already exists")
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
		Genre:    req.Genre,
		Bio:      req.Bio,
	}

	if err := s.ArtistRepository.Create(s.DB, artist); err != nil {
		s.Log.Warnf("failed to create artist in database: %+v", err)
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
		s.Log.Warnf("failed to update artist token and refresh token: %+v", err)
		return nil, errors.New(constants.INTERNAL_SERVER_ERROR)
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("failed to commit transaction: %+v", err)
		return nil, errors.New(constants.INTERNAL_SERVER_ERROR)
	}

	return converter.ArtistToReponse(newArtist), nil
}

func (s *ArtistService) UpdateArtist(ctx context.Context, req *model.UpdateArtistRequest) (*model.ArtistResponse, error) {
	s.Log.Info("starting Update Artist function")
	s.Log.Infof("request received: %+v", req)
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := s.Validate.Struct(req)
	if err != nil {
		s.Log.Warnf("invalid request body : %+v", err)
		return nil, errors.New(constants.BAD_REQUEST)
	}

	if req.Password == "" {
		return nil, errors.New("password is empty")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.Log.Warnf("failed to generate bcrypt hash: %+v", err)
		return nil, errors.New(constants.INTERNAL_SERVER_ERROR)
	}

	artist := &entity.Artist{
		Password: string(password),
		UserName: req.UserName,
		Genre:    req.Genre,
		Bio:      req.Bio,
	}

	err = s.ArtistRepository.Update(s.DB, artist)
	if err != nil {
		s.Log.Warnf("failed to update: %+v", err)
		return nil, errors.New(constants.INTERNAL_SERVER_ERROR)
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("failed commit transaction : %+v", err)
		return nil, errors.New(constants.INTERNAL_SERVER_ERROR)
	}
	return nil, nil
}

func (s *ArtistService) LogoutArtist(ctx context.Context, req *model.LogoutArtistRequest) error {
	s.Log.Info("starting Logout function")
	s.Log.Infof("request received: %+v", req)
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := s.Validate.Struct(req)
	if err != nil {
		s.Log.Warnf("invalid request body : %+v", err)
		return errors.New(constants.BAD_REQUEST)
	}

	artist := new(entity.Artist)
	_, err = s.ArtistRepository.FindByToken(s.DB, artist, req.Token)
	if err != nil {
		s.Log.Warnf("failed to find artist in database: %+v", err)
		return errors.New(constants.NOT_FOUND)
	}

	artist.Token = ""

	err = s.ArtistRepository.Update(tx, artist)
	if err != nil {
		s.Log.Warnf("failed to update artist: %+v", err)
		return errors.New(constants.INTERNAL_SERVER_ERROR)
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("failed commit transaction : %+v", err)
		return errors.New(constants.INTERNAL_SERVER_ERROR)
	}
	return nil
}

func (s *ArtistService) GetArtist(ctx context.Context, req *model.GetArtistRequest) (*model.ArtistResponse, error) {
	s.Log.Info("starting Get artist function")
	s.Log.Infof("request received: %+v", req)
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := s.Validate.Struct(req)
	if err != nil {
		s.Log.Warnf("invalid request body : %+v", err)
		return nil, errors.New(constants.BAD_REQUEST)
	}

	artist := new(entity.Artist)
	if err := s.ArtistRepository.FindById(s.DB, artist, req.ID); err != nil {
		s.Log.Warnf("failed to find artist in database: %+v", err)
		return nil, errors.New(constants.NOT_FOUND)
	}

	artist.Token = ""
	artist.RefreshToken = ""

	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("failed commit transaction : %+v", err)
		return nil, errors.New(constants.INTERNAL_SERVER_ERROR)
	}

	return converter.ArtistToReponse(artist), nil
}

func (s *ArtistService) DeleteArtist(ctx context.Context, req *model.DeleteArtistRequest) error {
	s.Log.Info("starting Delete artist function")
	s.Log.Infof("request received: %+v", req)
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := s.Validate.Struct(req)
	if err != nil {
		s.Log.Warnf("invalid request body : %+v", err)
		return errors.New(constants.BAD_REQUEST)
	}

	artist := new(entity.Artist)
	err = s.ArtistRepository.FindById(tx, artist, req.ID)
	if err != nil {
		s.Log.Warnf("failed to find artist by token : %+v", err)
		return errors.New(constants.INTERNAL_SERVER_ERROR)
	}

	if artist == nil {
		s.Log.Warnf("artist not found: %+v", err)
		return errors.New(constants.NOT_FOUND)
	}

	err = s.ArtistRepository.Delete(s.DB, artist)
	if err != nil {
		s.Log.Warnf("failed to delete artist by token : %+v", err)
		return errors.New(constants.INTERNAL_SERVER_ERROR)
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("failed commit transaction : %+v", err)
		return errors.New(constants.INTERNAL_SERVER_ERROR)
	}

	return nil
}
