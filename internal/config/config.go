package config

import "fmt"

type PostgresConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Name     string
	SSLMode  string
}

func ConvertPostgresConfigToConnectionString(p PostgresConfig) string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		p.Username, p.Password, p.Host, p.Port, p.Name, p.SSLMode)
}

type ServerConfig struct {
	Host string
	Port string
}

type OAuthConfig struct {
	ClientID     string
	ClientSecret string
}

type JWTConfig struct {
	SigningKey string
}

type LogLevel string

const (
	Debug LogLevel = "debug"
	Info  LogLevel = "info"
	Warn  LogLevel = "warn"
	Error LogLevel = "error"
)

func ConvertLogLevel(lvl string) LogLevel {
	switch lvl {
	case "debug":
		return Debug
	case "info":
		return Info
	case "warn":
		return Warn
	case "error":
		return Error
	}

	return Debug
}
