package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Newella-HQ/newella-backend/internal/logger"
	"github.com/Newella-HQ/newella-backend/internal/server"
	"github.com/Newella-HQ/newella-backend/internal/static-server/config"
	"github.com/Newella-HQ/newella-backend/internal/static-server/handler"
)

func main() {
	cfg, err := config.InitStaticServerConfig()
	if err != nil {
		log.Fatalf("can't init config: %s\n", err)
	}

	zapLogger, err := logger.NewZapLogger(cfg.LogLevel)
	if err != nil {
		log.Fatalf("can't init logger: %s\n", err)
	}
	defer func(zapLogger *logger.ZapLogger) {
		err := zapLogger.Sync()
		if err != nil && (!errors.Is(err, syscall.EBADF) && !errors.Is(err, syscall.ENOTTY)) {
			zapLogger.Errorf("can't sync logger: %s", err)
		}
	}(zapLogger)

	serverContext, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := createAssetsDirectory("./assets"); err != nil {
		zapLogger.Fatalf("can't create assets directory: %s", err)
	}

	h := handler.NewHandler(zapLogger)

	srv := server.NewHTTPServer(cfg.ServerConfig.Port, h.InitRoutes())

	go func() {
		if err := srv.Start(); !errors.Is(err, http.ErrServerClosed) {
			zapLogger.Fatalf("can't start server: %s", err)
		}
	}()
	defer func() {
		if err := srv.Shutdown(serverContext); err != nil {
			zapLogger.Errorf("can't shutdown server: %s", err)
		}
	}()

	zapLogger.Infof("server started on port: %s", cfg.ServerConfig.Port)

	<-serverContext.Done()
}

func createAssetsDirectory(path string) error {
	return os.MkdirAll(path, os.ModeDir)
}
