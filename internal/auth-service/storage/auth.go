package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/Masterminds/squirrel"

	"github.com/Newella-HQ/newella-backend/internal/logger"
	"github.com/Newella-HQ/newella-backend/internal/model"
)

type AuthStorage struct {
	logger    logger.Logger
	dbConn    *pgx.Conn
	dbTimeout time.Duration
}

const (
	defaultGreeting = "Hello, I'm "
)

var psqlBuilder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func NewAuthStorage(logger logger.Logger, dbConn *pgx.Conn, dbTimeout time.Duration) *AuthStorage {
	return &AuthStorage{
		dbConn:    dbConn,
		logger:    logger,
		dbTimeout: dbTimeout,
	}
}

func (s *AuthStorage) SaveUser(ctx context.Context, oauthJWT model.OAuthJWTToken, pair model.TokenPair) (string, error) {
	username := strings.Split(oauthJWT.Email, "@")[0]
	description := defaultGreeting + username

	insertQuery := `INSERT INTO users (id, username, real_name, description, email, picture)
					VALUES ($1, $2, $3, $4, $5, $6)
					ON CONFLICT (id)
					DO UPDATE
					SET id=excluded.id,
						username=excluded.username,
						real_name=excluded.real_name,
						description=excluded.description,
						email=excluded.email,
						picture=excluded.picture
					RETURNING role`

	insertTokens := `INSERT INTO oauth_tokens (access_token, refresh_token, user_id)
					 VALUES ($1, $2, $3)
					 ON CONFLICT (user_id)
					 DO UPDATE
					 SET access_token=excluded.access_token,
						 refresh_token=excluded.refresh_token`

	dbCtx, cancel := context.WithTimeout(ctx, s.dbTimeout)
	defer cancel()

	tx, err := s.dbConn.Begin(dbCtx)
	if err != nil {
		return "", fmt.Errorf("can't begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var role string
	if err := tx.QueryRow(dbCtx, insertQuery,
		oauthJWT.Subject, username, oauthJWT.Name,
		description, oauthJWT.Email, oauthJWT.Picture).Scan(&role); err != nil {
		return "", fmt.Errorf("can't insert user to db: %w", err)
	}

	if _, err := tx.Exec(dbCtx, insertTokens, pair.AccessToken, pair.RefreshToken, oauthJWT.Subject); err != nil {
		return "", fmt.Errorf("can't insert tokens to db: %w", err)
	}

	if err := tx.Commit(dbCtx); err != nil {
		return "", fmt.Errorf("can't commit tx: %w", err)
	}

	return role, nil
}

func (s *AuthStorage) GetTokensPair(ctx context.Context, refreshToken, userID string) (*model.TokenPair, error) {
	query := `SELECT access_token, refresh_token FROM oauth_tokens WHERE user_id=$1 AND refresh_token=$2`

	dbCtx, cancel := context.WithTimeout(ctx, s.dbTimeout)
	defer cancel()

	var pair model.TokenPair
	if err := s.dbConn.QueryRow(dbCtx, query, userID, refreshToken).Scan(&pair.AccessToken, &pair.RefreshToken); err != nil {
		return nil, fmt.Errorf("can't query tokens: %w", err)
	}
	return &pair, nil
}

func (s *AuthStorage) UpdateTokens(ctx context.Context, pair model.TokenPair, userID string) error {
	insertTokens := `INSERT INTO oauth_tokens (access_token, refresh_token, user_id)
					 VALUES ($1, $2, $3)
					 ON CONFLICT (user_id)
					 DO UPDATE
					 SET access_token=excluded.access_token,
						 refresh_token=excluded.refresh_token`

	dbCtx, cancel := context.WithTimeout(ctx, s.dbTimeout)
	defer cancel()

	if _, err := s.dbConn.Exec(dbCtx, insertTokens, pair.AccessToken, pair.RefreshToken, userID); err != nil {
		return fmt.Errorf("can't update tokens db: %w", err)
	}

	return nil
}

func (s *AuthStorage) RemoveTokens(ctx context.Context, userID string) error {
	deleteTokens := `DELETE FROM oauth_tokens WHERE user_id=$1`

	dbCtx, cancel := context.WithTimeout(ctx, s.dbTimeout)
	defer cancel()

	if _, err := s.dbConn.Exec(dbCtx, deleteTokens, userID); err != nil {
		return fmt.Errorf("can't delete tokens db: %w", err)
	}

	return nil
}
