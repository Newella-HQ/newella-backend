package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"github.com/Newella-HQ/newella-backend/internal/config"
)

type AuthServiceConfig struct {
	PostgresConfig config.PostgresConfig
	ServerConfig   config.ServerConfig
	OAuthConfig    config.OAuthConfig
	JWTConfig      config.JWTConfig
}

func InitAuthServiceConfig() (*AuthServiceConfig, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("can't get env variables: %w", err)
	}

	return &AuthServiceConfig{
		PostgresConfig: config.PostgresConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			Username: os.Getenv("POSTGRES_USERNAME"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Name:     os.Getenv("POSTGRES_NAME"),
			SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
		},
		ServerConfig: config.ServerConfig{
			Port: os.Getenv("SERVER_PORT"),
		},
		OAuthConfig: config.OAuthConfig{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		},
		JWTConfig: config.JWTConfig{
			SigningKey: os.Getenv("JWT_SIGNING_KEY"),
		},
	}, nil
}
