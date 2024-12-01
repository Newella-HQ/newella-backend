package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Newella-HQ/newella-backend/internal/model"
)

func (h *Handler) GetOAuthURL(ctx *gin.Context) {
	ctx.Header("Cache-Control", "no-cache, no-store")
	urlString := h.authService.GenerateAuthURL()

	ctx.JSON(http.StatusOK, model.Response{
		Message: "ok",
		Payload: gin.H{
			"url": urlString,
		},
	})
}

func (h *Handler) RedirectHandler(ctx *gin.Context) {
	ctx.Header("Cache-Control", "no-cache, no-store")
	state := ctx.Query("state")
	code := ctx.Query("code")

	jwks, err := h.getOAuthKeys(ctx)
	if err != nil {
		h.logger.Errorf("can't get jwks: %s", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, model.Response{
			Message: "can't get jwks",
			Payload: nil,
		})
		return
	}

	if exists := h.authService.VerifyState(state); !exists {
		h.logger.Warnf("got unexpected state: %s", state)
		ctx.AbortWithStatusJSON(http.StatusForbidden, model.Response{
			Message: "bad state value",
			Payload: nil,
		})
		return
	}
	h.logger.Debugf("got expected state: %s", state)

	tokenPair, err := h.authService.GetTokens(ctx, code, state, jwks)
	if err != nil {
		h.logger.Errorf("can't get tokens: %s", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, model.Response{
			Message: "can't get tokens",
			Payload: nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Message: "ok",
		Payload: tokenPair,
	})
}

func (h *Handler) getOAuthKeys(ctx context.Context) (json.RawMessage, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.googleapis.com/oauth2/v3/certs", nil)
	if err != nil {
		return nil, fmt.Errorf("can't creaete request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't do request: %w", err)
	}
	defer res.Body.Close()

	jwks, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("can' read body: %w", err)
	}

	return jwks, nil
}

func (h *Handler) RefreshTokens(ctx *gin.Context) {
	var tokenPair model.TokenPair
	if err := ctx.ShouldBindJSON(&tokenPair); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, model.Response{
			Message: "can't parse json body",
			Payload: nil,
		})
	}

	jwks, err := h.getOAuthKeys(ctx)
	if err != nil {
		h.logger.Errorf("can't get jwks: %s", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, model.Response{
			Message: "can't get jwks",
			Payload: nil,
		})
		return
	}

	newTokenPair, err := h.authService.RefreshTokens(ctx, tokenPair, jwks)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.logger.Debugf("tokens not found: %s", err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, model.Response{
				Message: "you're logged out",
				Payload: nil,
			})
			return
		}
		h.logger.Errorf("can't refresh tokens: %s", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, model.Response{
			Message: "can't refresh tokens",
			Payload: nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Message: "ok",
		Payload: newTokenPair,
	})
}

func (h *Handler) Logout(ctx *gin.Context) {
	userID, _, err := h.GetUserDataFromContext(ctx)
	if err != nil {
		h.logger.Errorf("can't get user data ctx: %w", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, model.Response{
			Message: "can't get user data ctx",
			Payload: nil,
		})
		return
	}

	if err := h.authService.RemoveTokens(ctx, userID); err != nil {
		h.logger.Errorf("can't remove tokens: %w", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, model.Response{
			Message: "can't remove tokens",
			Payload: nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Message: "ok",
		Payload: nil,
	})
}
