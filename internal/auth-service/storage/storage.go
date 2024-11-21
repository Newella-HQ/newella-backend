package storage

import "github.com/jackc/pgx/v5"

type Auth interface {
}

type AuthStorage struct {
	dbConn *pgx.Conn
}

func NewAuthStorage(dbConn *pgx.Conn) *AuthStorage {
	return &AuthStorage{dbConn: dbConn}
}
