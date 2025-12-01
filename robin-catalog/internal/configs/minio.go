package configs

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewMinioClient(viper *viper.Viper, log *logrus.Logger) *minio.Client {
	minioClient, err := minio.New(viper.GetString("MINIO_ENDPOINT"), &minio.Options{
		Creds:  credentials.NewStaticV4(viper.GetString("MINIO_ACCESS_KEY_ID"), viper.GetString("MINIO_SECRET_ACCESS_KEY"), ""),
		Secure: true,
	})

	if err != nil {
		log.Fatalln(err)
	}

	err = minioClient.MakeBucket(context.Background(), viper.GetString("MINIO_BUCKET_NAME"), minio.MakeBucketOptions{})

	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(context.Background(), viper.GetString("MINIO_BUCKET"))
		if errBucketExists != nil && exists {
			log.Printf("We already own %s\n", viper.GetString("MINIO_BUCKET_NAME"))
		} else {
			log.Fatalln(err)
		}
	}
	
	log.Info("MINIO BUCKET CREATED")

	return minioClient
}