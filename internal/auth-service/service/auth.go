package service

import (
	"github.com/Newella-HQ/newella-backend/internal/auth-service/storage"
)

type Auth interface {
	SignUp()
}

type AuthService struct {
	storage storage.Auth
}

func NewAuthService(storage storage.Auth) *AuthService {
	return &AuthService{storage: storage}
}

func (s *AuthService) SignUp() {

}
