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
		Secure: false, 
	})

	if err != nil {
		log.Fatalln(err)
	}

	bucketName := viper.GetString("MINIO_BUCKET_NAME")
	err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})

	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			log.Printf("Bucket %s already exists\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Infof("MinIO bucket '%s' created successfully", bucketName)
	}
	
	log.Info("MINIO BUCKET CREATED")

	return minioClient
}