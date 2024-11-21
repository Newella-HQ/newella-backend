package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Newella-HQ/newella-backend/internal/auth-service/service"
)

type Handler struct {
	authService service.Auth
}

func NewHandler(authService service.Auth) *Handler {
	return &Handler{authService: authService}
}

func (h *Handler) InitRoutes() http.Handler {
	r := chi.NewRouter()

	return r
}
