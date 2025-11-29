package services

import (
	"context"
	"net/http"

	"github.com/Dokito555/robin/robin-catalog/internal/entity"
	"github.com/Dokito555/robin/robin-catalog/internal/model"
	"github.com/Dokito555/robin/robin-catalog/internal/model/converter"
	"github.com/Dokito555/robin/robin-catalog/internal/repository"
	"github.com/Dokito555/robin/robin-catalog/utils/errs"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AlbumService struct {
	DB              *gorm.DB
	Log             *logrus.Logger
	Validate       *validator.Validate
	AlbumRepository *repository.AlbumRepository
}

func NewAlbumService(db *gorm.DB, log *logrus.Logger, validator *validator.Validate, AlbumRepository *repository.AlbumRepository) *AlbumService {
	return &AlbumService{
		DB:              db,
		Log:             log,
		Validate:       validator,
		AlbumRepository: AlbumRepository,
	}
}

func (s *AlbumService) CreateAlbum(ctx context.Context, req *model.CreateAlbumRequest) (*model.AlbumResponse, error) {
	s.Log.Info("starting Create Album function")
	s.Log.Infof("request received: %+v", req)

	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := s.Validate.Struct(req)
	if err != nil {
		s.Log.Warnf("invalid request body: %+v", err)
		return nil, errs.ERROR_BAD_REQUEST
	}

	if req.Name == "" || req.Type == "" {
		return nil, errs.NewError(http.StatusBadRequest, "name or type are required")
	}

	newAlbum := &entity.Album{
		ArtistID: req.ArtistID,
		Name: req.Name,
		Type: req.Type,
	}
	err = s.AlbumRepository.Create(s.DB, newAlbum)
	if err != nil {
		s.Log.Warnf("failed to create new album in database: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	
	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("failed to commit transaction: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	return converter.AlbumToResponse(newAlbum), nil
}

func (s *AlbumService) GetAlbum(ctx context.Context, req *model.GetAlbumRequest) (*model.AlbumResponse, error) {
	s.Log.Info("starting Get Album function")
	s.Log.Infof("request received: %+v", req)

	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := s.Validate.Struct(req)
	if err != nil {
		s.Log.Warnf("invalid request body: %+v", err)
		return nil, errs.ERROR_BAD_REQUEST
	}

	album := new(entity.Album)
	err = s.AlbumRepository.FindById(s.DB, album, req.ID)
	if err != nil {
		s.Log.Warnf("failed to find album in db: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("failed to commit transaction: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	return converter.AlbumToResponse(album), nil
}

func (s *AlbumService) UpdateAlbum(ctx context.Context, req *model.UpdateAlbumRequest) (*model.AlbumResponse, error) {
	s.Log.Info("starting Update Album function")
	s.Log.Infof("request received: %+v", req)

	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := s.Validate.Struct(req)
	if err != nil {
		s.Log.Warnf("invalid request body: %+v", err)
		return nil, errs.ERROR_BAD_REQUEST
	}

	album := new(entity.Album)
	if err := s.AlbumRepository.FindById(s.DB, album, req.ID); err != nil {
		s.Log.Warnf("album not found: %+v", err)
		return nil, errs.ERROR_NOT_FOUND
	}

	album.Name = req.Name
	album.Type = req.Type

	err = s.AlbumRepository.Update(s.DB, album)
	if err != nil {
		s.Log.Warnf("failed to update album in db: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("failed to commit transaction: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	return converter.AlbumToResponse(album), nil
}

func (s *AlbumService) DeleteAlbum(ctx context.Context, req *model.DeleteAlbumRequest) error {
	s.Log.Info("starting Delete Album function")
	s.Log.Infof("request received: %+v", req)

	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := s.Validate.Struct(req)
	if err != nil {
		s.Log.Warnf("invalid request body: %+v", err)
		return errs.ERROR_BAD_REQUEST
	}

	album := new(entity.Album)
	err = s.AlbumRepository.FindById(s.DB, album, req.ID)
	if err != nil {
		s.Log.Warnf("failed to find album in db: %+v", err)
		return errs.ERROR_INTERNAL_SERVER_ERROR
	}

	if album.ID == 0 {
		s.Log.Warnf("album not found: %+v", err)
		return errs.ERROR_NOT_FOUND
	}

	err = s.AlbumRepository.Delete(s.DB, album)
	if err != nil {
		s.Log.Warnf("failed to delete album in db: %+v", err)
		return errs.ERROR_INTERNAL_SERVER_ERROR
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("failed to commit transaction: %+v", err)
		return errs.ERROR_INTERNAL_SERVER_ERROR
	}

	return nil
}

func (s *AlbumService) GetAlbumListByArtistID(ctx context.Context, req *model.GetAlbumRequest) ([]model.AlbumResponse, error) {
	s.Log.Info("starting Get List Albums function")
	s.Log.Infof("request received: %+v", req)

	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := s.Validate.Struct(req)
	if err != nil {
		s.Log.Warnf("invalid request body: %+v", err)
		return nil, errs.ERROR_BAD_REQUEST
	}

	if req.ID == 0 {
		return nil, errs.ERROR_BAD_REQUEST
	}

	albums, err := s.AlbumRepository.AlbumListByArtistID(s.DB, req.ID)
	if err != nil {
		s.Log.Warn("failed to get list of albums in database")
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	rsps := make([]model.AlbumResponse, len(albums))
	for i, album := range albums {
		rsps[i] = *converter.AlbumToResponse(&album)
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("failed to commit transaction: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	return rsps, nil
}