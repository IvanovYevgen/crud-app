package config

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

type Config struct {
	DB Postgres

	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`

	Auth struct {
		TokenTTL time.Duration `mapstructure:"token_ttl"`
	} `mapstructure:"auth"`
}

type Postgres struct {
	Host     string `envconfig:"DB_HOST"`
	Port     int    `envconfig:"DB_PORT"`
	Username string `envconfig:"DB_USER"`
	Name     string `envconfig:"DB_NAME"`
	SSLMode  string `envconfig:"DB_SSL_MODE"`
	Password string `envconfig:"DB_PASSWORD"`
}

func New(folder, filename string) (*Config, error) {
	cfg := new(Config)

	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("ERROR LOADING .ENV FILE")
	}

	viper.AddConfigPath(folder)
	viper.SetConfigName(filename)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	if err := envconfig.Process("", &cfg.DB); err != nil {
		return nil, err
	}

	fmt.Println("cfg:", cfg)

	return cfg, nil
}
