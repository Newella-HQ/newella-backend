package auth

import (
	"context"
	ssov2 "github.com/Newella-HQ/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context,
		email string,
		password string,
		appID int,
	) (token string, err error)
	RegisterNewUser(ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
	IsAdmin(ctx context.Context,
		userID int64,
	) (bool, error)
}

type serverApi struct {
	ssov2.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov2.RegisterAuthServer(gRPC, &serverApi{auth: auth})
}

const emptyValue = 0

func (s *serverApi) Login(
	ctx context.Context,
	req *ssov2.LoginRequest,
) (*ssov2.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov2.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverApi) Register(
	ctx context.Context,
	req *ssov2.RegisterRequest,
) (*ssov2.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov2.RegisterResponse{UserId: userID}, nil
}

func (s *serverApi) IsAdmin(
	ctx context.Context,
	req *ssov2.IsAdminRequest,
) (*ssov2.IsAdminResponse, error) {
	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &ssov2.IsAdminResponse{IsAdmin: isAdmin}, nil
}

func validateLogin(req *ssov2.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "email")
	}

	if req.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "email")
	}

	return nil
}

func validateRegister(req *ssov2.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email id required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "email id required")
	}

	return nil
}

func validateIsAdmin(req *ssov2.IsAdminRequest) error {
	if req.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}

	return nil
}
