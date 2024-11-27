package service

import (
	"time"

	"github.com/akyoto/cache"
	"golang.org/x/oauth2"

	"github.com/Newella-HQ/newella-backend/internal/logger"
)

const (
	_cacheTTL   = 5 * time.Minute
	_oauthState = "newella-auth-service"
)

type AuthStorage interface {
}

type AuthService struct {
	logger   logger.Logger
	storage  AuthStorage
	cache    *cache.Cache
	oauthCfg *oauth2.Config
}

func NewAuthService(logger logger.Logger, storage AuthStorage, cache *cache.Cache, oauthCfg *oauth2.Config) *AuthService {
	return &AuthService{
		storage:  storage,
		logger:   logger,
		cache:    cache,
		oauthCfg: oauthCfg,
	}
}

func (s *AuthService) GenerateAuthURL() string {
	code := oauth2.GenerateVerifier()
	s.cache.Set(code, true, _cacheTTL)
	s.logger.Infof("code = %s", code)
	opt := oauth2.S256ChallengeOption(code)

	return s.oauthCfg.AuthCodeURL(_oauthState, opt)
}
