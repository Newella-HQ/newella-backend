package server

import (
	"context"
	"net/http"
)

type AuthServiceServer struct {
	s *http.Server
}

func NewAuthServiceServer(port string, handler http.Handler) *AuthServiceServer {
	return &AuthServiceServer{
		s: &http.Server{
			Handler: handler,
			Addr:    ":" + port,
		},
	}
}

func (s *AuthServiceServer) Start() error {
	return s.s.ListenAndServe()
}

func (s *AuthServiceServer) Shutdown(ctx context.Context) error {
	return s.s.Shutdown(ctx)
}
