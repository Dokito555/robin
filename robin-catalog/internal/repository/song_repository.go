package repository

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/Dokito555/robin/robin-catalog/internal/entity"
	"github.com/Dokito555/robin/robin-catalog/internal/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type SongRepository struct {
	Repository[entity.Song]
	S3  *s3.Client
	MinioClient *minio.Client
	Viper *viper.Viper
	Log *logrus.Logger
}

func NewSongRepository(log *logrus.Logger, db *gorm.DB, m *minio.Client, v *viper.Viper) *SongRepository {
	return &SongRepository{
		Repository: Repository[entity.Song]{DB: db},
		Log:        log,
		MinioClient: m,
		Viper: v,
	}
}

func (r *SongRepository) UploadFileToMinio(filePath, fileName string) (string, error) {
	_, err := r.MinioClient.FPutObject(
		context.Background(),
		r.Viper.GetString("MINIO_BUCKET_NAME"),
		fileName,
		filePath,
		minio.PutObjectOptions{
			ContentType: "audio/mpeg",
		},
	)

	if err != nil {
		return "", err
	}

	return fileName, nil
}

func (r *SongRepository) RetrieveFileFromMinio(fileName string) (*url.URL, error) {
	url, err := r.MinioClient.PresignedGetObject(
		context.Background(),
		r.Viper.GetString("MINIO_BUCKET_NAME"),
		fileName,
		time.Hour,
		nil,
	)

	if err != nil {
		return nil, err
	}

	return url, nil
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

func (r *SongRepository) UploadFileToS3(bucketName string, model *model.File) (string, error) {
	// timemout after 30 seconds
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// read file (might be not necessary)
	fileStat, _ := model.File.Seek(0, 2) // get file size
	model.File.Seek(0, 0)                // reset the file pointer
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
