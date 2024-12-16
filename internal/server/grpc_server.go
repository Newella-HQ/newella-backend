package server

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Newella-HQ/newella-backend/internal/logger"
)

type GRPCServer struct {
	srv    *grpc.Server
	logger logger.Logger
}

func NewGRPCServer(logger logger.Logger, withLoggingInterceptor, withRecoveryInterceptor bool) *GRPCServer {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived, logging.PayloadSent,
		),
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p any) (err error) {
			logger.Errorf("recovered from panic: %s", p)

			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	var interceptors []grpc.UnaryServerInterceptor

	if withLoggingInterceptor {
		interceptors = append(interceptors,
			logging.UnaryServerInterceptor(InterceptorLogger(logger), loggingOpts...),
		)
	}

	if withRecoveryInterceptor {
		interceptors = append(interceptors,
			recovery.UnaryServerInterceptor(recoveryOpts...),
		)
	}

	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptors...))

	return &GRPCServer{
		srv:    srv,
		logger: logger,
	}
}

func InterceptorLogger(l logger.Logger) logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, _ logging.Level, msg string, fields ...any) {
		l.With(fields...).Debugf(msg)
	})
}

func filterNilsAndEmptyStrings(l []any) []any {
	var filtered []any
	for _, v := range l {
		if v != nil {
			if str := fmt.Sprintf("%s", v); str == "" {
				filtered = append(filtered, " ")
				continue
			}
			filtered = append(filtered, v)
		}
	}
	return filtered
}

func (g *GRPCServer) Run(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("can't start tcp listener: %w", err)
	}

	if err := g.srv.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return fmt.Errorf("server was stopped: %w", err)
	}

	return nil
}

func (g *GRPCServer) GracefulShutdown() {
	g.logger.Infof("gracefully stopping the grpc server")

	g.srv.GracefulStop()
}

func (g *GRPCServer) Register(sd *grpc.ServiceDesc, ss any) {
	g.logger.Infof("register new service: %T", ss)

	g.srv.RegisterService(sd, ss)
}
