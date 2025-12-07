package main

import "github.com/Dokito555/robin-ums/internal/config"

func main() {
	// init configs
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	validate := config.NewValidator(viperConfig)
	app := config.NewGin(viperConfig)
	// kafkaProducer, err := config.NewKafkaProducer(viperConfig, log)
	// if err != nil {
	// 	log.Fatalf("Failed to create Kafka producer: %v", err)
	// }
	// defer kafkaProducer.Close()
	// should be grpc

	// inject configs to app
	config.Bootstrap(&config.BootstrapConfig{
		DB:            db,
		App:           app,
		Log:           log,
		Validate:      validate,
		Config:        viperConfig,
		// KafkaProducer: kafkaProducer,
	})

	// run app
	port := viperConfig.GetString("APP_PORT")
	err := app.Run(":" + port)
	log.Info("Listening to port: " + port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
