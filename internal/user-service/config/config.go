package config

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"

	"github.com/Newella-HQ/newella-backend/internal/config"
)

type UserServiceConfig struct {
	PostgresConfig  config.PostgresConfig
	ServerConfig    config.ServerConfig
	JWTConfig       config.JWTConfig
	LogLevel        config.LogLevel
	DatabaseTimeout time.Duration
}

func InitUserServiceConfig() (cfg *UserServiceConfig, err error) {
	defer func() {
		if recErr := recover(); recErr != nil {
			err = fmt.Errorf("can't init auth config: %s", recErr)
		}
	}()

	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("can't get env variables: %w", err)
	}

	return &UserServiceConfig{
		PostgresConfig: config.PostgresConfig{
			Host:     config.GetAndValidateEnv("POSTGRES_HOST"),
			Port:     config.GetAndValidateEnv("POSTGRES_PORT"),
			Username: config.GetAndValidateEnv("POSTGRES_USERNAME"),
			Password: config.GetAndValidateEnv("POSTGRES_PASSWORD"),
			Name:     config.GetAndValidateEnv("POSTGRES_NAME"),
			SSLMode:  config.GetAndValidateEnv("POSTGRES_SSLMODE"),
		},
		ServerConfig: config.ServerConfig{
			Host: config.GetAndValidateEnv("SERVER_HOST"),
			Port: config.GetAndValidateEnv("USER_SERVICE_GRPC_PORT"),
		},
		JWTConfig: config.JWTConfig{
			SigningKey: config.GetAndValidateEnv("JWT_SIGNING_KEY"),
		},
		LogLevel:        config.ConvertLogLevel(config.GetAndValidateEnv("LOG_LEVEL")),
		DatabaseTimeout: 15 * time.Second,
	}, nil
}
