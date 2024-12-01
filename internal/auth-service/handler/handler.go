package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Newella-HQ/newella-backend/internal/logger"
	"github.com/Newella-HQ/newella-backend/internal/model"
)

type AuthService interface {
	GenerateAuthURL() string
	VerifyState(string) bool
	GetTokens(ctx context.Context, code, state string, jwks json.RawMessage) (*model.TokenPair, error)
	RefreshTokens(ctx context.Context, tokenPair model.TokenPair, jwks json.RawMessage) (*model.TokenPair, error)
	ParseAccessToken(signed string) (*model.NewellaJWTToken, error)
	RemoveTokens(ctx context.Context, userID string) error
}

type Handler struct {
	logger      logger.Logger
	authService AuthService
}

func NewHandler(logger logger.Logger, authService AuthService) *Handler {
	return &Handler{
		authService: authService,
		logger:      logger,
	}
}

func (h *Handler) InitRoutes() http.Handler {
	r := gin.New()

	r.GET("/auth", h.GetOAuthURL)
	r.GET("/redirect", h.RedirectHandler)
	r.POST("/refresh", h.UserIdentify, h.RefreshTokens)
	r.DELETE("/logout", h.UserIdentify, h.Logout)

	return r
}

func (h *Handler) UserIdentify(ctx *gin.Context) {
	header := ctx.GetHeader("Authorization")
	if header == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, model.Response{
			Message: "empty auth header",
			Payload: nil,
		})
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, model.Response{
			Message: "invalid auth header",
			Payload: nil,
		})
		return
	}
	if len(headerParts[1]) == 0 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, model.Response{
			Message: "empty auth token",
			Payload: nil,
		})
		return
	}

	jwtToken, err := h.authService.ParseAccessToken(headerParts[1])
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, model.Response{
			Message: "can't parse access token: " + err.Error(),
			Payload: nil,
		})
		return
	}

	ctx.Set("userID", jwtToken.UserID)
	ctx.Set("role", jwtToken.Role)
	ctx.Next()
}

func (h *Handler) GetUserDataFromContext(ctx *gin.Context) (uID string, r string, err error) {
	userID, ok := ctx.Get("userID")
	if !ok {
		return "", "", fmt.Errorf("userID is empty in context")
	}
	userIDStr, ok := userID.(string)
	if !ok {
		return "", "", fmt.Errorf("userID isn't string")
	}

	role, ok := ctx.Get("role")
	if !ok {
		return "", "", fmt.Errorf("role is empty in context")
	}
	roleStr, ok := role.(string)
	if !ok {
		return "", "", fmt.Errorf("role isn't string")
	}

	return userIDStr, roleStr, nil
}
