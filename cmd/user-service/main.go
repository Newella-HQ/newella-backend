package main

import (
	"context"
	"errors"
	"log"
	"os/signal"
	"syscall"

	"github.com/Newella-HQ/protos/gen/go/user"
	"github.com/jackc/pgx/v5"

	sharedcfg "github.com/Newella-HQ/newella-backend/internal/config"
	"github.com/Newella-HQ/newella-backend/internal/logger"
	"github.com/Newella-HQ/newella-backend/internal/server"
	"github.com/Newella-HQ/newella-backend/internal/user-service/config"
	"github.com/Newella-HQ/newella-backend/internal/user-service/handler"
	"github.com/Newella-HQ/newella-backend/internal/user-service/storage"
)

func main() {
	cfg, err := config.InitUserServiceConfig()
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

	dbConn, err := pgx.Connect(serverContext, sharedcfg.ConvertPostgresConfigToConnectionString(cfg.PostgresConfig))
	if err != nil {
		zapLogger.Fatalf("can't connect to db: %s", err)
	}
	defer func() {
		if err := dbConn.Close(serverContext); err != nil {
			zapLogger.Errorf("can't close db connection: %s", err)
		}
	}()

	postgresStorage := storage.NewUserStorage(zapLogger, dbConn, cfg.DatabaseTimeout)
	h := handler.NewUserServiceHandler(zapLogger, postgresStorage, cfg.JWTConfig.SigningKey)

	srv := server.NewGRPCServer(zapLogger, true, true)

	srv.Register(&user.UserService_ServiceDesc, h)

	go func() {
		if err := srv.Run(cfg.ServerConfig.Port); err != nil {
			zapLogger.Fatalf("can't ")
		}
	}()
	defer srv.GracefulShutdown()

	zapLogger.Infof("started grpc server at port: %s", cfg.ServerConfig.Port)

	<-serverContext.Done()
}
