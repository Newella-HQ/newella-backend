package storage

import (
	"github.com/jackc/pgx/v5"

	"github.com/Newella-HQ/newella-backend/internal/logger"
)

type AuthStorage struct {
	logger logger.Logger
	dbConn *pgx.Conn
}

func NewAuthStorage(logger logger.Logger, dbConn *pgx.Conn) *AuthStorage {
	return &AuthStorage{
		dbConn: dbConn,
		logger: logger,
	}
}
