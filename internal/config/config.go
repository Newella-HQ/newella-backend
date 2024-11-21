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

type ServerConfig struct {
	Port string
}

type OAuthConfig struct {
	ClientID     string
	ClientSecret string
}

type JWTConfig struct {
	SigningKey string
}

func ConvertPostgresConfigToConnectionString(p PostgresConfig) string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		p.Username, p.Password, p.Host, p.Port, p.Name, p.SSLMode)
}
