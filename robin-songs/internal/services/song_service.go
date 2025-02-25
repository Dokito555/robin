package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Dokito555/robin-songs/internal/entity"
	"github.com/Dokito555/robin-songs/internal/model"
	"github.com/Dokito555/robin-songs/internal/model/converter"
	"github.com/Dokito555/robin-songs/internal/repository"
	"github.com/Dokito555/robin-songs/utils/errs"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type SongService struct {
	Log            *logrus.Logger
	DB             *gorm.DB
	Validate       *validator.Validate
	Viper          *viper.Viper
	SongRepository *repository.SongRepository
}

func NewSongService(log *logrus.Logger, db *gorm.DB, validate *validator.Validate, viper *viper.Viper, repo *repository.SongRepository) *SongService {
	return &SongService{
		Log:            log,
		DB:             db,
		Validate:       validate,
		Viper:          viper,
		SongRepository: repo,
	}
}

func (s *SongService) CreateNewSong(ctx context.Context, req *model.CreateSongRequest, fileReq *model.UploadFileRequest) (*model.SongResponse, error) {
	s.Log.Info("starting Create New Song Function")
	s.Log.Infof("request received: %+v", req)

	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := s.Validate.Struct(req)
	if err != nil {
		s.Log.Warnf("invalid request body: %+v", err)
		return nil, errs.ERROR_BAD_REQUEST
	}

	// call upload file to s3 and return link
	fileName := fmt.Sprintf("uploads/%d-%s", time.Now().Unix(), fileReq.FileHeader.Filename)
	file := &model.File{
		File:     fileReq.File,
		FileName: fileName,
	}

	url, err := s.SongRepository.UploadFileToS3(viper.GetString("AWS_SONG_BUCKET"), file)
	if err != nil {
		s.Log.Println(viper.GetString("AWS_SONG_BUCKET"))
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}
	// TODO: check if current user already have song

	newSong := &entity.Song{
		ArtistID: req.ArtistID,
		AlbumID:  req.AlbumID,
		Name:     req.Name,
		Link:     url,
		Duration: req.Duration,
	}

	if err := s.SongRepository.Create(s.DB, newSong); err != nil {
		s.Log.Warnf("failed to create song in database: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("failed to commit transaction: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	return converter.SongToResponse(newSong), nil
}

func (s *SongService) GetSong(ctx context.Context, req *model.GetSongRequest) (*model.SongResponse, error) {
	s.Log.Info("starting Get Song Function")
	s.Log.Infof("request received: %+v", req)

	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := s.Validate.Struct(req)
	if err != nil {
		s.Log.Warnf("invalid request body: %+v", err)
		return nil, errs.ERROR_BAD_REQUEST
	}

	song := new(entity.Song)
	if err := s.SongRepository.FindById(s.DB, song, req.ID); err != nil {
		s.Log.Warnf("failed to get song in database: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("failed to commit transaction: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	return converter.SongToResponse(song), nil
}

func (s *SongService) UpdateSong(ctx context.Context, req *model.UpdateSongRequest) (*model.SongResponse, error) {
	s.Log.Info("starting Update Song Function")
	s.Log.Infof("request received: %+v", req)

	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := s.Validate.Struct(req)
	if err != nil {
		s.Log.Warnf("invalid request body: %+v", err)
		return nil, errs.ERROR_BAD_REQUEST
	}

	song := new(entity.Song)
	if err := s.SongRepository.FindById(s.DB, song, req.ID); err != nil {
		s.Log.Warnf("failed to get song in database: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	song.Name = req.Name
	song.Link = req.Link
	song.Duration = req.Duration

	if err := s.SongRepository.Update(s.DB, song); err != nil {
		s.Log.Warnf("failed to update song in database: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	return converter.SongToResponse(song), nil
}

func (s *SongService) DeleteSong(ctx context.Context, req *model.DeleteSongRequest) error {
	s.Log.Info("starting Update Song Function")
	s.Log.Infof("request received: %+v", req)

	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := s.Validate.Struct(req)
	if err != nil {
		s.Log.Warnf("invalid request body: %+v", err)
		return errs.ERROR_BAD_REQUEST
	}

	song := new(entity.Song)
	err = s.SongRepository.FindById(tx, song, req.ID)
	if err != nil {
		s.Log.Warnf("failed to find song in database : %+v", err)
		return errs.ERROR_INTERNAL_SERVER_ERROR
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.Log.Warnf("artist not found: %+v", err)
		return errs.ERROR_NOT_FOUND
	}

	if err := s.SongRepository.Delete(s.DB, song); err != nil {
		s.Log.Warnf("failed to delete song in database : %+v", err)
		return errs.ERROR_INTERNAL_SERVER_ERROR
	}

	if err := s.SongRepository.Update(s.DB, song); err != nil {
		s.Log.Warnf("failed to update song in database: %+v", err)
		return errs.ERROR_INTERNAL_SERVER_ERROR
	}

	return nil
}
