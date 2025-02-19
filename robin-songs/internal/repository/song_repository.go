package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Dokito555/robin-songs/internal/entity"
	"github.com/Dokito555/robin-songs/internal/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SongRepository struct {
	Repository[entity.Song]
	S3  *s3.Client
	Log *logrus.Logger
}

func NewSongRepository(log *logrus.Logger, db *gorm.DB) *SongRepository {
	return &SongRepository{
		Repository: Repository[entity.Song]{DB: db},
		Log:        log,
	}
}

func (r *SongRepository) GetFileFromS3(bucketName, fileName string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	obj := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	}

	_, err := r.S3.GetObject(ctx, obj)
	if err != nil {
		return "", err
	}


	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, fileName)
	return url, nil
}

func (r *SongRepository) UploadFileToS3(bucketName string, model model.File) (string, error) {
	// Timemout after 30 seconds
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// Read file (might be not necessary)
	fileStat, _ := model.File.Seek(0, 2) // Get file size
	model.File.Seek(0, 0)                // Reset the file pointer
	buffer := make([]byte, fileStat)
	model.File.Read(buffer)

	input := &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(model.FileName),
		Body:        model.File,
		ContentType: aws.String("audio/mpeg"),
	}

	_, err := r.S3.PutObject(ctx, input)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, model.FileName)
	return url, nil
}
