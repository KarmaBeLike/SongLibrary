package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port           int    `mapstructure:"PORT"`
	DBHost         string `mapstructure:"DB_HOST"`
	DBPort         int    `mapstructure:"DB_PORT"`
	DBName         string `mapstructure:"DB_NAME"`
	DBUser         string `mapstructure:"DB_USER"`
	DBPassword     string `mapstructure:"DB_PASSWORD"`
	ExternalAPIURL string `mapstructure:"API_URL"`
}

func Load() (*Config, error) {
	config := &Config{}

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.AutomaticEnv() // Загружать переменные окружения

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	err = viper.Unmarshal(config)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	return config, nil
}
