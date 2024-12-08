package handler

import (
	"github.com/Newella-HQ/protos/gen/go/user"

	"github.com/Newella-HQ/newella-backend/internal/logger"
)

type UserStorage interface {
}

type UserServiceHandler struct {
	logger      logger.Logger
	userStorage UserStorage

	user.UnimplementedUserServiceServer
}

func NewUserServiceHandler(logger logger.Logger, userStorage UserStorage) *UserServiceHandler {
	return &UserServiceHandler{
		logger:      logger,
		userStorage: userStorage,
	}
}
