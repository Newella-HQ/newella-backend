package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/akyoto/cache"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"

	"github.com/Newella-HQ/newella-backend/internal/logger"
	"github.com/Newella-HQ/newella-backend/internal/model"
)

const (
	_cacheTTL = 5 * time.Minute
)

type AuthStorage interface {
	SaveUser(ctx context.Context, oauthJWT model.OAuthJWTToken, pair model.TokenPair) (string, error)
	GetTokensPair(ctx context.Context, refreshToken, userID string) (*model.TokenPair, error)
	UpdateTokens(ctx context.Context, pair model.TokenPair, userID string) error
	RemoveTokens(ctx context.Context, userID string) error
}

type AuthService struct {
	logger     logger.Logger
	storage    AuthStorage
	cache      *cache.Cache
	oauthCfg   *oauth2.Config
	singingKey string
}

func NewAuthService(logger logger.Logger, storage AuthStorage, cache *cache.Cache, oauthCfg *oauth2.Config, sgnKey string) *AuthService {
	return &AuthService{
		storage:    storage,
		logger:     logger,
		cache:      cache,
		oauthCfg:   oauthCfg,
		singingKey: sgnKey,
	}
}

func (s *AuthService) GenerateAuthURL() string {
	code := oauth2.GenerateVerifier()
	s.cache.Set(code, true, _cacheTTL)

	options := []oauth2.AuthCodeOption{
		oauth2.S256ChallengeOption(code),
		oauth2.SetAuthURLParam("access_type", "offline"),
		oauth2.SetAuthURLParam("prompt", "consent"),
	}

	return s.oauthCfg.AuthCodeURL(code, options...)
}

func (s *AuthService) VerifyState(state string) bool {
	if _, exists := s.cache.Get(state); exists {
		s.cache.Delete(state)
		return true
	}
	return false
}

func (s *AuthService) getOAuthJWTToken(token *oauth2.Token, jwks json.RawMessage) (model.OAuthJWTToken, error) {
	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		return model.OAuthJWTToken{}, fmt.Errorf("can't cast id_token to string")
	}

	k, err := keyfunc.NewJWKSetJSON(jwks)
	if err != nil {
		return model.OAuthJWTToken{}, fmt.Errorf("can't create a keyfunc: %w", err)
	}

	parsedToken, err := jwt.ParseWithClaims(idToken, &model.OAuthJWTToken{}, k.Keyfunc)
	if err != nil {
		return model.OAuthJWTToken{}, fmt.Errorf("can't parse jwt: %w", err)
	}

	claims, ok := parsedToken.Claims.(*model.OAuthJWTToken)
	if !ok {
		return model.OAuthJWTToken{}, fmt.Errorf("can't cast claims to %T", &model.OAuthJWTToken{})
	}

	validatedJWT, err := claims.Validate()
	if err != nil {
		return model.OAuthJWTToken{}, fmt.Errorf("can't validate jwt: %w", err)
	}

	return validatedJWT, nil
}

func (s *AuthService) getSignedNewellaJWTToken(oauthToken model.OAuthJWTToken, role string) (string, error) {
	unsignedJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, &model.NewellaJWTToken{
		UserID:         oauthToken.Subject,
		Role:           role,
		Email:          oauthToken.Email,
		EmailVerified:  oauthToken.EmailVerified,
		Audience:       oauthToken.Audience,
		ExpirationTime: oauthToken.ExpirationTime,
		IssuedAt:       oauthToken.IssuedAt,
		Issuer:         oauthToken.Issuer,
	})

	signedAccessJWT, err := unsignedJWT.SignedString([]byte(s.singingKey))
	if err != nil {
		return "", fmt.Errorf("can't sing JWT: %w", err)
	}

	return signedAccessJWT, nil
}

func (s *AuthService) GetTokens(ctx context.Context, code, state string, jwks json.RawMessage) (*model.TokenPair, error) {
	token, err := s.oauthCfg.Exchange(ctx, code, oauth2.VerifierOption(state))
	if err != nil {
		return nil, fmt.Errorf("can't exchange: %w", err)
	}

	validatedJWT, err := s.getOAuthJWTToken(token, jwks)
	if err != nil {
		return nil, fmt.Errorf("can't get OAuthJWTToken: %w", err)
	}

	role, err := s.storage.SaveUser(ctx, validatedJWT, model.TokenPair{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	})
	if err != nil {
		return nil, fmt.Errorf("can't save user to db: %w", err)
	}

	signedAccessJWT, err := s.getSignedNewellaJWTToken(validatedJWT, role)
	if err != nil {
		return nil, fmt.Errorf("can't singed JWT: %w", err)
	}

	s.logger.Debugf("user authenticated: %s", validatedJWT.Email)

	return &model.TokenPair{
		AccessToken:  signedAccessJWT,
		RefreshToken: token.RefreshToken,
	}, nil
}

func (s *AuthService) RefreshTokens(ctx context.Context, tokenPair model.TokenPair, jwks json.RawMessage) (*model.TokenPair, error) {
	accessToken, err := s.ParseAccessToken(tokenPair.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("can't parse access token: %w", err)
	}

	dbTokens, err := s.storage.GetTokensPair(ctx, tokenPair.RefreshToken, accessToken.UserID)
	if err != nil {
		return nil, fmt.Errorf("can't get tokens from db: %w", err)
	}

	newToken, err := s.oauthCfg.TokenSource(ctx, &oauth2.Token{
		AccessToken:  dbTokens.AccessToken,
		TokenType:    "Bearer",
		RefreshToken: dbTokens.RefreshToken,
		Expiry:       time.Unix(accessToken.ExpirationTime, 0).UTC(),
		ExpiresIn:    0,
	}).Token()
	if err != nil {
		return nil, fmt.Errorf("can't refresh tokens: %w", err)
	}

	s.logger.Infof("access: %s, expiry: %s", newToken.AccessToken, newToken.Expiry)

	validatedNewJWT := model.OAuthJWTToken{
		Audience:       accessToken.Audience,
		ExpirationTime: newToken.Expiry.Unix(),
		IssuedAt:       accessToken.IssuedAt,
		Issuer:         accessToken.Issuer,
		Subject:        accessToken.UserID,
		Email:          accessToken.Email,
		EmailVerified:  accessToken.EmailVerified,
	}

	signedAccessJWT, err := s.getSignedNewellaJWTToken(validatedNewJWT, accessToken.Role)
	if err != nil {
		return nil, fmt.Errorf("can't singed JWT: %w", err)
	}

	if err := s.storage.UpdateTokens(ctx, model.TokenPair{
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
	}, validatedNewJWT.Subject); err != nil {
		return nil, fmt.Errorf("can't update tokens: %w", err)
	}

	return &model.TokenPair{
		AccessToken:  signedAccessJWT,
		RefreshToken: newToken.RefreshToken,
	}, nil
}

func (s *AuthService) ParseAccessToken(signed string) (*model.NewellaJWTToken, error) {
	token, err := jwt.ParseWithClaims(signed, &model.NewellaJWTToken{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("ivalid signing method")
		}
		return []byte(s.singingKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("can't parse jwt token: %w", err)
	}

	accessToken, ok := token.Claims.(*model.NewellaJWTToken)
	if !ok {
		return nil, fmt.Errorf("can't cast to newella jwt token: %w", err)
	}

	return accessToken, nil
}

func (s *AuthService) RemoveTokens(ctx context.Context, userID string) error {
	return s.storage.RemoveTokens(ctx, userID)
}
