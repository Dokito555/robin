package main

import config "github.com/Dokito555/robin/robin-catalog/internal/configs"

func main() {
	// init configs
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	validate := config.NewValidator(viperConfig)
	minioClient := config.NewMinioClient(viperConfig, log)
	app := config.NewGin(viperConfig)
	// TODO: switch to minio
	// s3 := config.NewS3Client(viperConfig, log)
	// should be grpc

	// inject configs to app
	config.Bootstrap(&config.BootstrapConfig{
		DB:          db,
		App:         app,
		Log:         log,
		Validate:    validate,
		Config:      viperConfig,
		MinioClient: minioClient,
		// S3Client: s3,
	})

	// run app
	port := viperConfig.GetString("APP_PORT")
	err := app.Run(":" + port)
	log.Info("Listening to port: " + port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
