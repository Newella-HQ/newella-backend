package server

import (
	"context"
	"net/http"
	"time"
)

type ServiceServer struct {
	s *http.Server
}

func NewServiceServer(port string, handler http.Handler) *ServiceServer {
	return &ServiceServer{
		s: &http.Server{
			Handler:           handler,
			Addr:              ":" + port,
			ReadHeaderTimeout: 10 * time.Second,
		},
	}
}

func (s *ServiceServer) Start() error {
	return s.s.ListenAndServe()
}

func (s *ServiceServer) Shutdown(ctx context.Context) error {
	return s.s.Shutdown(ctx)
}
