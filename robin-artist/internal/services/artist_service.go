package services

import (
	"context"
	"errors"
	"net/http"

	"github.com/Dokito555/robin-artist/constants"
	"github.com/Dokito555/robin-artist/internal/entity"
	"github.com/Dokito555/robin-artist/internal/model/converter"
	"github.com/Dokito555/robin-artist/internal/repository"
	"github.com/Dokito555/robin-artist/utils/errs"
	"github.com/Dokito555/robin-artist/internal/model"
	"github.com/go-playground/validator/v10"
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
		return nil, errs.ERROR_BAD_REQUEST
	}

	if req.Password == "" || req.Email == "" {
		return nil, errs.NewError(http.StatusBadRequest, "password or email required")
	}

	artist, err := s.ArtistRepository.FindByEmail(s.DB, &entity.Artist{}, req.Email)
	if err != nil {
		s.Log.Warnf("database error fetching artist: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	if artist != nil {
		s.Log.Warnf("artist with that email already exists")
		return nil, errs.ERROR_USER_EXIST
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.Log.Warnf("failed to generate bcrypt hash: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
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
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("failed to commit transaction: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
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
		return nil, errs.ERROR_BAD_REQUEST
	}

	if req.Password == "" || req.Email == "" {
		return nil, errs.NewError(http.StatusBadRequest, "passwor or email are required")
	}

	newArtist := new(entity.Artist)
	artist, err := s.ArtistRepository.FindByEmail(s.DB, newArtist, req.Email)
	if err != nil {
		s.Log.Warnf("database error fetching artist: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	if err := bcrypt.CompareHashAndPassword([]byte(artist.Password), []byte(req.Password)); err != nil {
		s.Log.Warnf("failed to compare hashed password: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	token, err := s.TokenService.GenerateToken(ctx, artist.ID, constants.TOKEN_TYPE_TOKEN, artist.Email, artist.Role, artist.UserName)
	if err != nil {
		s.Log.Warnf("failed to genereate token: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	refreshToken, err := s.TokenService.GenerateToken(ctx, artist.ID, constants.TOKEN_TYPE_REFRESH, artist.Email, artist.Role, artist.UserName)
	if err != nil {
		s.Log.Warnf("failed to genereate refresh token: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	newArtist.Token = token
	newArtist.RefreshToken = refreshToken

	err = s.ArtistRepository.Update(s.DB, newArtist)
	if err != nil {
		s.Log.Warnf("failed to update artist token and refresh token: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("failed to commit transaction: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
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
		return nil, errs.ERROR_BAD_REQUEST
	}

	if req.Password == "" {
		return nil, errors.New("password is empty")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.Log.Warnf("failed to generate bcrypt hash: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	var artist entity.Artist
	if err := tx.First(&artist, req.ID).Error; err != nil {
		s.Log.Warnf("artist not found: %+v", err)
		return nil, errs.ERROR_NOT_FOUND
	}

	artist.Password = string(password)
	artist.UserName = req.UserName
	artist.Genre = req.Genre
	artist.Bio = req.Bio

	err = s.ArtistRepository.Update(tx, &artist)
	if err != nil {
		s.Log.Warnf("failed to update: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	artist.RefreshToken = ""
	artist.Token = ""

	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("failed commit transaction : %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	return converter.ArtistToReponse(&artist), nil
}

func (s *ArtistService) LogoutArtist(ctx context.Context, req *model.LogoutArtistRequest) error {
	s.Log.Info("starting Logout function")
	s.Log.Infof("request received: %+v", req)
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := s.Validate.Struct(req)
	if err != nil {
		s.Log.Warnf("invalid request body : %+v", err)
		return errs.ERROR_BAD_REQUEST
	}

	artist := new(entity.Artist)
	_, err = s.ArtistRepository.FindByToken(s.DB, artist, req.Token)
	if err != nil {
		s.Log.Warnf("failed to find artist in database: %+v", err)
		return errs.ERROR_NOT_FOUND
	}

	artist.Token = ""

	err = s.ArtistRepository.Update(tx, artist)
	if err != nil {
		s.Log.Warnf("failed to update artist: %+v", err)
		return errs.ERROR_INTERNAL_SERVER_ERROR
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("failed commit transaction : %+v", err)
		return errs.ERROR_INTERNAL_SERVER_ERROR
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
		return nil, errs.ERROR_BAD_REQUEST
	}

	artist := new(entity.Artist)
	if err := s.ArtistRepository.FindById(s.DB, artist, req.ID); err != nil {
		s.Log.Warnf("failed to find artist in database: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	artist.Token = ""
	artist.RefreshToken = ""

	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("failed commit transaction : %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
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
		return errs.ERROR_BAD_REQUEST
	}

	artist := new(entity.Artist)
	err = s.ArtistRepository.FindById(tx, artist, req.ID)
	if err != nil {
		s.Log.Warnf("failed to find artist by token : %+v", err)
		return errs.ERROR_INTERNAL_SERVER_ERROR
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.Log.Warnf("artist not found: %+v", err)
		return errs.ERROR_NOT_FOUND
	}

	err = s.ArtistRepository.Delete(s.DB, artist)
	if err != nil {
		s.Log.Warnf("failed to delete artist by token : %+v", err)
		return errs.ERROR_INTERNAL_SERVER_ERROR
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("failed commit transaction : %+v", err)
		return errs.ERROR_INTERNAL_SERVER_ERROR
	}

	return nil
}

func (s *ArtistService) Verify(ctx context.Context, req *model.VerifyRequest) (*model.ArtistResponse, error) {
	s.Log.Info("starting Verify function")
	s.Log.Infof("request received: %+v", req)
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := s.Validate.Struct(req)
	if err != nil {
		s.Log.Warnf("invalid request body : %+v", err)
		return nil, errs.ERROR_BAD_REQUEST
	}

	user := new(entity.Artist)
	rsp, err := s.ArtistRepository.FindByToken(tx, user, req.Token)
	if err != nil {
		s.Log.Warnf("failed find user by token : %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("failed commit transaction : %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	return converter.ArtistToReponse(rsp), nil
}

func (s *ArtistService) GetArtistList(ctx context.Context, page int, limit int) ([]model.ArtistResponse, error) {
	s.Log.Info("starting Get List Artist function")

	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if page == 0 || limit == 0 {
		return nil, errs.ERROR_BAD_REQUEST
	}

	artists, err := s.ArtistRepository.GetArtistList(s.DB, page, limit)
	if err != nil {
		s.Log.Warn("failed to get list of artists in database")
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("failed to commit transaction: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	rsps := make([]model.ArtistResponse, len(artists))
	for i, artist := range artists {
		rsps[i] = *converter.ArtistToReponse(&artist)
	}

	return rsps, nil
}