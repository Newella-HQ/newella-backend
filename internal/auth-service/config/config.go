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

func InitAuthServiceConfig() (cfg *AuthServiceConfig, err error) {
	defer func() {
		if recErr := recover(); recErr != nil {
			err = fmt.Errorf("empty config parameter")
		}
	}()

	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("can't get env variables: %w", err)
	}

	return &AuthServiceConfig{
		PostgresConfig: config.PostgresConfig{
			Host:     panicIfEmpty(os.Getenv("POSTGRES_HOST")),
			Port:     panicIfEmpty(os.Getenv("POSTGRES_PORT")),
			Username: panicIfEmpty(os.Getenv("POSTGRES_USERNAME")),
			Password: panicIfEmpty(os.Getenv("POSTGRES_PASSWORD")),
			Name:     panicIfEmpty(os.Getenv("POSTGRES_NAME")),
			SSLMode:  panicIfEmpty(os.Getenv("POSTGRES_SSLMODE")),
		},
		ServerConfig: config.ServerConfig{
			Port: panicIfEmpty(os.Getenv("SERVER_PORT")),
		},
		OAuthConfig: config.OAuthConfig{
			ClientID:     panicIfEmpty(os.Getenv("GOOGLE_CLIENT_ID")),
			ClientSecret: panicIfEmpty(os.Getenv("GOOGLE_CLIENT_SECRET")),
		},
		JWTConfig: config.JWTConfig{
			SigningKey: panicIfEmpty(os.Getenv("JWT_SIGNING_KEY")),
		},
	}, nil
}

func panicIfEmpty(s string) string {
	if s == "" {
		panic("empty config parameter")
	}
	return s
}
