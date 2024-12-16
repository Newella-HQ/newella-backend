package handler

import (
	"context"

	"github.com/Newella-HQ/protos/gen/go/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Newella-HQ/newella-backend/internal/logger"
)

type UserStorage interface {
	GetUser(ctx context.Context, id string) (*user.User, error)
	GetUsers(ctx context.Context, search string, limit, offset int) (int, []*user.User, error)
}

type UserServiceHandler struct {
	logger      logger.Logger
	userStorage UserStorage
	singingKey  string

	user.UnimplementedUserServiceServer
}

func NewUserServiceHandler(logger logger.Logger, userStorage UserStorage, singingKey string) *UserServiceHandler {
	return &UserServiceHandler{
		logger:      logger,
		userStorage: userStorage,
		singingKey:  singingKey,
	}
}

// GetUser with authority
func (h *UserServiceHandler) GetUser(ctx context.Context, request *user.GetUserRequest) (*user.GetUserResponse, error) {
	userToken, err := h.VerifyAndGetUserToken(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}
	h.logger.Debugf("got user token: %v", userToken)

	id := request.GetId()
	if id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "empty ID field")
	}

	userData, err := h.userStorage.GetUser(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can't get user from db: %s", err)
	}

	return &user.GetUserResponse{User: userData}, nil
}

func (h *UserServiceHandler) GetUsers(ctx context.Context, request *user.GetUsersRequest) (*user.GetUsersResponse, error) {
	limit, page := 50, 0
	if l := request.GetLimit(); l >= 0 {
		limit = int(l)
	}
	if p := request.GetPage(); p >= 0 {
		page = int(p)
	}
	offset := limit * page
	search := request.GetSearch()

	count, users, err := h.userStorage.GetUsers(ctx, search, limit, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can't get users from db: %s", err)
	}

	return &user.GetUsersResponse{
		Count: int64(count),
		Users: users,
	}, nil
}

func (h *UserServiceHandler) ChangeUsername(ctx context.Context, request *user.ChangeUsernameRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (h *UserServiceHandler) ChangeUserData(ctx context.Context, request *user.ChangeUserDataRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (h *UserServiceHandler) ChangePicture(server user.UserService_ChangePictureServer) error {
	//TODO implement me
	panic("implement me")
}

func (h *UserServiceHandler) ChangeRole(ctx context.Context, request *user.ChangeRoleRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (h *UserServiceHandler) GetSubscribers(ctx context.Context, request *user.GetSubsRequest) (*user.GetSubsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (h *UserServiceHandler) GetSubscriptions(ctx context.Context, request *user.GetSubsRequest) (*user.GetSubsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (h *UserServiceHandler) Subscribe(ctx context.Context, request *user.SubRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (h *UserServiceHandler) Unsubscribe(ctx context.Context, request *user.SubRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (h *UserServiceHandler) DeleteSubscriber(ctx context.Context, request *user.SubRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}
