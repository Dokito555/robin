package configs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type S3Client struct {
	Client     *s3.Client
	Uploader   *manager.Uploader
	BucketName string
}

func NewS3Client(v *viper.Viper, log *logrus.Logger) *S3Client {
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(v.GetString("BUCKET_REGION")),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			v.GetString("AWS_ACCESS_KEY_ID"),
			v.GetString("AWS_SECRET_ACCESS_KEY"),
			"",
		)),
	)
	if err != nil {
		log.Fatalf("unable to load AWS SDK config, %v", err)
	}

	client := s3.NewFromConfig(awsCfg)
	uploader := manager.NewUploader(client)
	bucket := v.GetString("AWS_SONG_BUCKET")

	return &S3Client{
		Client:     client,
		Uploader:   uploader,
		BucketName: bucket,
	}
}