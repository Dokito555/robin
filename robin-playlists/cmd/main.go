package main

import "github.com/Dokito555/robin-playlists/internal/config"

func main() {
	// init configs
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	validate := config.NewValidator(viperConfig)
	app := config.NewGin(viperConfig)
	// should be grpc

	// inject configs to app
	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
		Config:   viperConfig,
	})

	// run app
	port := viperConfig.GetString("APP_PORT")
	err := app.Run(":" + port)
	log.Info("Listening to port: " + port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
