package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"github.com/Newella-HQ/newella-backend/internal/auth-service/config"
	"github.com/Newella-HQ/newella-backend/internal/auth-service/handler"
	"github.com/Newella-HQ/newella-backend/internal/auth-service/server"
	"github.com/Newella-HQ/newella-backend/internal/auth-service/service"
	"github.com/Newella-HQ/newella-backend/internal/auth-service/storage"
	sharedcfg "github.com/Newella-HQ/newella-backend/internal/config"
)

func main() {
	l, _ := zap.NewProduction()
	defer func(l *zap.Logger) {
		if err := l.Sync(); err != nil {
			log.Fatalln("can't flush logs")
		}
	}(l)
	logger := l.Sugar()

	authConfig, err := config.InitAuthServiceConfig()
	if err != nil {
		logger.Fatalln("can't init config")
	}

	serverContext, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	dbConn, err := pgx.Connect(serverContext, sharedcfg.ConvertPostgresConfigToConnectionString(authConfig.PostgresConfig))
	if err != nil {
		logger.Fatalln("can't connect to db")
	}
	defer func() {
		if err := dbConn.Close(serverContext); err != nil {
			logger.Errorf("can't close db connection: %s", err)
		}
	}()

	authStorage := storage.NewAuthStorage(dbConn)
	authService := service.NewAuthService(authStorage)
	authHandler := handler.NewHandler(authService)

	srv := server.NewAuthServiceServer(authConfig.ServerConfig.Port, authHandler.InitRoutes())

	go func() {
		if err := srv.Start(); !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("can't start server: %s", err)
		}
	}()
	defer func() {
		if err := srv.Shutdown(serverContext); err != nil {
			logger.Errorf("can't shutdown server: %s", err)
		}
	}()

	logger.Infof("server started on port: %s", authConfig.ServerConfig.Port)

	<-serverContext.Done()
}
