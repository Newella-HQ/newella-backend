package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Newella-HQ/newella-backend/internal/model"
)

func (h *Handler) GetOAuthURL(ctx *gin.Context) {
	url := h.authService.GenerateAuthURL()

	ctx.JSON(http.StatusOK, model.Response{
		Message: "ok",
		Payload: gin.H{
			"url": url,
		},
	})
}

func (h *Handler) RedirectHandler(ctx *gin.Context) {
	h.logger.Infoln(ctx.Request.URL)

	ctx.String(http.StatusOK, "success")
}

func (h *Handler) RefreshTokens(ctx *gin.Context) {

}

func (h *Handler) Logout(ctx *gin.Context) {

}
