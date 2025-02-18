package services

import (
	"context"
	"net/http"

	"github.com/Dokito555/robin-playlists/internal/entity"
	model "github.com/Dokito555/robin-playlists/internal/models"
	"github.com/Dokito555/robin-playlists/internal/models/converter"
	"github.com/Dokito555/robin-playlists/internal/repository"
	"github.com/Dokito555/robin-playlists/internal/utils/errs"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PlaylistService struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validate           *validator.Validate
	PlaylistRepository *repository.PlaylistRepository
}

func NewPlaylistService(db *gorm.DB, log *logrus.Logger, validator *validator.Validate, Repo *repository.PlaylistRepository) *PlaylistService {
	return &PlaylistService{
		DB:                 db,
		Log:                log,
		Validate:           validator,
		PlaylistRepository: Repo,
	}
}

func (s *PlaylistService) CreatePlaylist(ctx context.Context, req *model.CreatePlaylistRequest) (*model.PlaylistResponse, error) {
	s.Log.Info("starting Create Playlist function")
	s.Log.Infof("request received: %+v", req)

	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := s.Validate.Struct(req)
	if err != nil {
		s.Log.Warnf("invalid request body: %+v", err)
		return nil, errs.ERROR_BAD_REQUEST
	}

	if req.Name == "" {
		return nil, errs.NewError(http.StatusBadRequest, "name is required")
	}

	newPlaylist := &entity.Playlist{
		UserID:      req.UserID,
		Name:        req.Name,
		Description: req.Description,
	}

	err = s.PlaylistRepository.Create(s.DB, newPlaylist)
	if err != nil {
		s.Log.Warnf("failed to create new playlist in database: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("failed to commit transaction: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	return converter.PlaylistToResponse(newPlaylist), nil
}

func (s *PlaylistService) GetPlaylist(ctx context.Context, req *model.GetPlaylistRequest) (*model.PlaylistResponse, error) {
	s.Log.Info("starting Get Playlist function")
	s.Log.Infof("request received: %+v", req)

	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := s.Validate.Struct(req)
	if err != nil {
		s.Log.Warnf("invalid request body: %+v", err)
		return nil, errs.ERROR_BAD_REQUEST
	}

	playlist := new(entity.Playlist)
	err = s.PlaylistRepository.FindById(s.DB, playlist, req.ID)
	if err != nil {
		s.Log.Warnf("failed to find playlist in db: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("failed to commit transaction: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	return converter.PlaylistToResponse(playlist), nil
}

func (s *PlaylistService) UpdatePlaylist(ctx context.Context, req *model.UpdatePlaylistRequest) (*model.PlaylistResponse, error) {
	s.Log.Info("starting Update Playlist function")
	s.Log.Infof("request received: %+v", req)

	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := s.Validate.Struct(req)
	if err != nil {
		s.Log.Warnf("invalid request body: %+v", err)
		return nil, errs.ERROR_BAD_REQUEST
	}

	playlist := new(entity.Playlist)
	if err := s.PlaylistRepository.FindById(s.DB, playlist, req.ID); err != nil {
		return nil, errs.ERROR_NOT_FOUND
	}

	playlist.Name = req.Name
	playlist.Description = req.Description

	err = s.PlaylistRepository.Update(s.DB, playlist)
	if err != nil {
		s.Log.Warnf("failed to update playlist in db: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("failed to commit transaction: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	return converter.PlaylistToResponse(playlist), nil
}

func (s *PlaylistService) DeletePlaylist(ctx context.Context, req *model.DeletePlaylistRequest) error {
	s.Log.Info("starting Delete Album function")
	s.Log.Infof("request received: %+v", req)

	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := s.Validate.Struct(req)
	if err != nil {
		s.Log.Warnf("invalid request body: %+v", err)
		return errs.ERROR_BAD_REQUEST
	}

	playlist := new(entity.Playlist)
	err = s.PlaylistRepository.FindById(s.DB, playlist, req.ID)
	if err != nil {
		s.Log.Warnf("failed to find playlist in db: %+v", err)
		return errs.ERROR_INTERNAL_SERVER_ERROR
	}

	if playlist.ID == 0 {
		s.Log.Warnf("playlist not found: %+v", err)
		return errs.ERROR_NOT_FOUND
	}

	err = s.PlaylistRepository.Delete(s.DB, playlist)
	if err != nil {
		s.Log.Warnf("failed to delete playlist in db: %+v", err)
		return errs.ERROR_INTERNAL_SERVER_ERROR
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("failed to commit transaction: %+v", err)
		return errs.ERROR_INTERNAL_SERVER_ERROR
	}

	return nil
}
