package config

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/Newella-HQ/newella-backend/internal/config"
)

type AuthServiceConfig struct {
	PostgresConfig  config.PostgresConfig
	ServerConfig    config.ServerConfig
	OAuthConfig     config.OAuthConfig
	JWTConfig       config.JWTConfig
	LogLevel        config.LogLevel
	DatabaseTimeout time.Duration
}

func InitAuthServiceConfig() (cfg *AuthServiceConfig, err error) {
	defer func() {
		if recErr := recover(); recErr != nil {
			err = fmt.Errorf("can't init auth config: %s", recErr)
		}
	}()

	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("can't get env variables: %w", err)
	}

	return &AuthServiceConfig{
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
			Port: config.GetAndValidateEnv("AUTH_SERVER_PORT"),
		},
		OAuthConfig: config.OAuthConfig{
			ClientID:     config.GetAndValidateEnv("GOOGLE_CLIENT_ID"),
			ClientSecret: config.GetAndValidateEnv("GOOGLE_CLIENT_SECRET"),
		},
		JWTConfig: config.JWTConfig{
			SigningKey: config.GetAndValidateEnv("JWT_SIGNING_KEY"),
		},
		LogLevel:        config.ConvertLogLevel(config.GetAndValidateEnv("LOG_LEVEL")),
		DatabaseTimeout: 15 * time.Second,
	}, nil
}

const (
	UserInfoEmailScope   = "https://www.googleapis.com/auth/userinfo.email"
	UserInfoProfileScope = "https://www.googleapis.com/auth/userinfo.profile"
	OpenIDScope          = "openid"
)

func (cfg *AuthServiceConfig) NewOAuth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     cfg.OAuthConfig.ClientID,
		ClientSecret: cfg.OAuthConfig.ClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  getRedirectURL(cfg.ServerConfig),
		Scopes:       []string{UserInfoProfileScope, UserInfoEmailScope, OpenIDScope},
	}
}

func getRedirectURL(cfg config.ServerConfig) string {
	return fmt.Sprintf("http://%s:%s/redirect", cfg.Host, cfg.Port)
}
