package storage

import (
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/Newella-HQ/newella-backend/internal/logger"
)

type UserStorage struct {
	logger    logger.Logger
	dbConn    *pgx.Conn
	dbTimeout time.Duration
}

func NewUserStorage(logger logger.Logger, dbConn *pgx.Conn, dbTimeout time.Duration) *UserStorage {
	return &UserStorage{
		logger:    logger,
		dbConn:    dbConn,
		dbTimeout: dbTimeout,
	}
}
