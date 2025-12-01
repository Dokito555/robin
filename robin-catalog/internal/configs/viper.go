package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

func NewViper() *viper.Viper {
	config := viper.New()

	config.SetConfigName(".env") 
	config.SetConfigType("env")  
	config.AddConfigPath("./")
	config.AddConfigPath("./../") 
	
	config.AutomaticEnv()

	err := config.ReadInConfig()
	if err != nil {
		fmt.Printf("No .env file found, relying on environment variables: %v\n", err)
	}

	return config
}
