package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"time"

	"github.com/Dokito555/robin/robin-catalog/internal/entity"
	"github.com/Dokito555/robin/robin-catalog/internal/model"
	"github.com/Dokito555/robin/robin-catalog/internal/model/converter"
	"github.com/Dokito555/robin/robin-catalog/internal/repository"
	"github.com/Dokito555/robin/robin-catalog/utils/errs"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tcolgate/mp3"
	"gorm.io/gorm"
)

type SongService struct {
	Log             *logrus.Logger
	DB              *gorm.DB
	Validate        *validator.Validate
	Viper           *viper.Viper
	SongRepository  *repository.SongRepository
	AlbumRepository *repository.AlbumRepository
}

func NewSongService(log *logrus.Logger, db *gorm.DB, validate *validator.Validate, viper *viper.Viper, repo *repository.SongRepository, albumRepo *repository.AlbumRepository) *SongService {
	return &SongService{
		Log:            log,
		DB:             db,
		Validate:       validate,
		Viper:          viper,
		SongRepository: repo,
		AlbumRepository: albumRepo,
	}
}

func (s *SongService) CalculateMP3Duration(file multipart.File) (time.Duration, error) {
	var duration time.Duration
	decoder := mp3.NewDecoder(file)
	var frame mp3.Frame
	skipped := 0

	for {
		err := decoder.Decode(&frame, &skipped)
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, fmt.Errorf("error decoding MP3 frame: %w", err)
		}
		duration += frame.Duration()
	}

	return duration, nil
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

	album := new(entity.Album)
	if err := s.AlbumRepository.FindById(s.DB, album, req.AlbumID); err != nil {
    return nil, errs.ERROR_NOT_FOUND
	}

	if album.ArtistID != req.ArtistID {
		return nil, errs.ERROR_FORBIDDEN
	}

	fileName := fmt.Sprintf("songs/%d-%s", time.Now().Unix(), fileReq.FileHeader.Filename)

	tempFilePath := fmt.Sprintf("/tmp/%s", fileReq.FileHeader.Filename)
	// tight coupling
	// err = ctx.Value("ginCotext").(*gin.Context).SaveUploadedFile(fileReq.FileHeader, tempFilePath)
	// if err != nil {
	// 	s.Log.Warnf("failed to save temp file: %v", err)
	// 	return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	// }
	// defer os.Remove(tempFilePath)

	// chance of messing up file streams
	out, err := os.Create(tempFilePath)
	if err != nil {
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}
	defer out.Close()
	defer os.Remove(tempFilePath)

	if _, err := io.Copy(out, fileReq.File); err != nil {
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	objectName, err := s.SongRepository.UploadFileToMinio(tempFilePath, fileName)
	if err != nil {
		s.Log.Warnf("failed to upload file to MinIO: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	// file := &model.File{
	// 	File:     fileReq.File,
	// 	FileName: fileName,
	// }

	// url, err := s.SongRepository.UploadFileToS3(viper.GetString("AWS_SONG_BUCKET"), file)
	// if err != nil {
	// 	s.Log.Println(viper.GetString("AWS_SONG_BUCKET"))
	// 	return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	// }
	// TODO: check if current user already have song

	newSong := &entity.Song{
		ArtistID: req.ArtistID,
		AlbumID:  req.AlbumID,
		Name:     req.Name,
		Link:     objectName,
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

	return converter.SongToResponse(song), nil
}

func (s *SongService) GetSongStream(ctx context.Context, req *model.GetSongRequest) (*model.SongStreamResponse, error) {
	s.Log.Info("starting Get Song Stream URL Function")
	s.Log.Infof("request received: %+v", req)

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

	streamURL, err := s.SongRepository.RetrieveFileFromMinio(song.Link)
	if err != nil {
		s.Log.Warnf("failed to generate presigned URL from MinIO: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	return &model.SongStreamResponse{
		SongID:    song.ID,
		Name:      song.Name,
		StreamURL: streamURL.String(),
		Duration:  song.Duration,
	}, nil
}

func (s *SongService) DeleteSong(ctx context.Context, req *model.DeleteSongRequest) error {
	s.Log.Info("starting Delete Song Function")
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

	// if err := s.SongRepository.Update(s.DB, song); err != nil {
	// 	s.Log.Warnf("failed to update song in database: %+v", err)
	// 	return errs.ERROR_INTERNAL_SERVER_ERROR
	// }

	return nil
}
