package config

import (
	"fmt"

	"github.com/joho/godotenv"

	"github.com/Newella-HQ/newella-backend/internal/config"
)

type StaticServerConfig struct {
	ServerConfig config.ServerConfig
	LogLevel     config.LogLevel
}

func InitStaticServerConfig() (cfg *StaticServerConfig, err error) {
	defer func() {
		if recErr := recover(); recErr != nil {
			err = fmt.Errorf("can't init static server config: %s", recErr)
		}
	}()

	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("can't get env variables: %w", err)
	}

	return &StaticServerConfig{
		ServerConfig: config.ServerConfig{
			Host: config.GetAndValidateEnv("SERVER_HOST"),
			Port: config.GetAndValidateEnv("STATIC_SERVER_PORT"),
		},
		LogLevel: config.ConvertLogLevel(config.GetAndValidateEnv("LOG_LEVEL")),
	}, nil
}
