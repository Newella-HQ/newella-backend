package handler

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/metadata"

	"github.com/Newella-HQ/newella-backend/internal/model"
	"github.com/Newella-HQ/newella-backend/internal/token"
)

const AuthHeader = "authorization"

func (h *UserServiceHandler) VerifyAndGetUserToken(ctx context.Context) (*model.NewellaJWTToken, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("can't get metadata from ctx")
	}

	authHeader, ok := md[AuthHeader]
	if !ok || len(authHeader) == 0 {
		return nil, fmt.Errorf("can't get authorization header from metadata")
	}

	headerParts := strings.Split(authHeader[0], " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" || len(headerParts[1]) == 0 {
		return nil, fmt.Errorf("invalid auth header")
	}

	parsedToken, err := token.ParseAccessToken(headerParts[0], h.singingKey)
	if err != nil {
		return nil, fmt.Errorf("can't parse access token: %w", err)
	}

	return parsedToken, nil
}
