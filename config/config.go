package config

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Server       ServerConfig
	Postgresql   PostgresConfig
	UrlShortener UrlShortenerConfig
	Kafka        KafkaConfig
}

type ServerConfig struct {
	Port string
}

type PostgresConfig struct {
	Host              string
	Port              string
	User              string
	Password          string
	DbName            string
	MaxConnection     string
	MaxConnectionIdle string
}

type UrlShortenerConfig struct {
	BaseUrl string
	Length  int
}

type KafkaConfig struct {
	Broker string
}

func LoadConfig() (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev" // Default to development if no environment is set
	}
	viper.SetConfigName(fmt.Sprintf("config.%s", env))
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config/")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	log.Printf("Loaded %s configuration", env)
	return &config, nil
}
