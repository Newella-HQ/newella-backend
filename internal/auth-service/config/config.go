package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/Newella-HQ/newella-backend/internal/config"
)

type AuthServiceConfig struct {
	PostgresConfig config.PostgresConfig
	ServerConfig   config.ServerConfig
	OAuthConfig    config.OAuthConfig
	JWTConfig      config.JWTConfig
	LogLevel       config.LogLevel
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
			Host:     GetAndValidateEnv("POSTGRES_HOST"),
			Port:     GetAndValidateEnv("POSTGRES_PORT"),
			Username: GetAndValidateEnv("POSTGRES_USERNAME"),
			Password: GetAndValidateEnv("POSTGRES_PASSWORD"),
			Name:     GetAndValidateEnv("POSTGRES_NAME"),
			SSLMode:  GetAndValidateEnv("POSTGRES_SSLMODE"),
		},
		ServerConfig: config.ServerConfig{
			Host: GetAndValidateEnv("SERVER_HOST"),
			Port: GetAndValidateEnv("AUTH_SERVER_PORT"),
		},
		OAuthConfig: config.OAuthConfig{
			ClientID:     GetAndValidateEnv("GOOGLE_CLIENT_ID"),
			ClientSecret: GetAndValidateEnv("GOOGLE_CLIENT_SECRET"),
		},
		JWTConfig: config.JWTConfig{
			SigningKey: GetAndValidateEnv("JWT_SIGNING_KEY"),
		},
		LogLevel: config.ConvertLogLevel(GetAndValidateEnv("LOG_LEVEL")),
	}, nil
}

func GetAndValidateEnv(key string) string {
	s := os.Getenv(key)
	if s == "" {
		panic(fmt.Sprintf("empty %s parameter", key))
	}

	return s
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
