package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Newella-HQ/newella-backend/internal/logger"
)

type Handler struct {
	logger logger.Logger
}

func NewHandler(logger logger.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func (h *Handler) InitRoutes() http.Handler {
	r := gin.New()
	r.Use(func(context *gin.Context) {
		h.logger.Debugf("got request for: %s", context.Request.URL.Path)
	})
	r.StaticFS("/assets", gin.Dir("./assets", false))
	return r
}
