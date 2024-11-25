package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Newella-HQ/newella-backend/internal/logger"
)

type AuthService interface {
	GenerateAuthURL() string
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
	r.POST("/refresh", h.RefreshTokens)
	r.DELETE("/logout", h.Logout)

	return r
}
