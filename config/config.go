package config

import (
	"errors"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	Env        string           `env:"ENV" envDefault:"development"`
	Host       string           `env:"HOST" envDefault:"localhost"`
	Port       string           `env:"PORT" envDefault:"5000"`
	Postgres   PostgresConfig   `envPrefix:"POSTGRES_"`
	JWT        JWTConfig        `envPrefix:"JWT_"`
	Cloudinary CloudinaryConfig `envPrefix:"CLOUDINARY"`
}

type PostgresConfig struct {
	Host     string `env:"HOST" envDefault:"localhost"`
	Port     string `env:"PORT" envDefault:"5432"`
	User     string `env:"USER" envDefault:"postgres"`
	Password string `env:"PASSWORD" envDefault:"postgres"`
	Database string `env:"DATABASE" envDefault:"postgres"`
}

type JWTConfig struct {
	SecretKey string `env:"SECRET_KEY" envDefault:"secret"`
	ExpiresAt int    `env:"EXPIRES_AT" envDefault:"24"`
}

type CloudinaryConfig struct {
}

func NewConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, errors.New("ERROR LOADING .ENV FILE")
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, errors.New("ERROR PARSING ENVIRONMENT VARIABLES")
	}

	return &cfg, nil
}
