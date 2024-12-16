package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/akyoto/cache"
	"github.com/jackc/pgx/v5"

	"github.com/Newella-HQ/newella-backend/internal/auth-service/config"
	"github.com/Newella-HQ/newella-backend/internal/auth-service/handler"
	"github.com/Newella-HQ/newella-backend/internal/auth-service/service"
	"github.com/Newella-HQ/newella-backend/internal/auth-service/storage"
	sharedcfg "github.com/Newella-HQ/newella-backend/internal/config"
	"github.com/Newella-HQ/newella-backend/internal/logger"
	"github.com/Newella-HQ/newella-backend/internal/server"
)

const (
	_cacheTTL = 5 * time.Minute
)

func main() {
	authConfig, err := config.InitAuthServiceConfig()
	if err != nil {
		log.Fatalf("can't init config: %s\n", err)
	}

	zapLogger, err := logger.NewZapLogger(authConfig.LogLevel)
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

	dbConn, err := pgx.Connect(serverContext, sharedcfg.ConvertPostgresConfigToConnectionString(authConfig.PostgresConfig))
	if err != nil {
		zapLogger.Fatalf("can't connect to db: %s", err)
	}
	defer func() {
		if err := dbConn.Close(serverContext); err != nil {
			zapLogger.Errorf("can't close db connection: %s", err)
		}
	}()

	oauthCodesCache := cache.New(_cacheTTL)
	oauthConfig := authConfig.NewOAuth2Config()

	authStorage := storage.NewAuthStorage(zapLogger, dbConn, authConfig.DatabaseTimeout)
	authService := service.NewAuthService(zapLogger, authStorage, oauthCodesCache, oauthConfig, authConfig.JWTConfig.SigningKey)
	authHandler := handler.NewHandler(zapLogger, authService)

	srv := server.NewHTTPServer(authConfig.ServerConfig.Port, authHandler.InitRoutes())

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

	zapLogger.Infof("server started on port: %s", authConfig.ServerConfig.Port)

	<-serverContext.Done()
}
