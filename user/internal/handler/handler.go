package handler

import (
	"context"

	pb "github.com/HJyup/translatify-common/api"
	models "github.com/HJyup/translatify-user/internal/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcHandler struct {
	pb.UnimplementedUserServiceServer
	service models.UserService
}

func NewGrpcHandler(grpcServer *grpc.Server, service models.UserService) {
	handler := &GrpcHandler{
		service: service,
	}
	pb.RegisterUserServiceServer(grpcServer, handler)
}

func (h *GrpcHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	token, err := h.service.CreateUser(
		req.GetUsername(),
		req.GetEmail(),
		req.GetPassword(),
		req.GetLanguage(),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}
	return &pb.CreateUserResponse{
		Token: token,
	}, nil
}

func (h *GrpcHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := h.service.GetUser(req.GetUsername())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}
	return &pb.GetUserResponse{
		User: &pb.User{
			Username:  user.Username,
			Email:     user.Email,
			Language:  user.Language,
			CreatedAt: user.CreatedAt.Unix(),
		},
	}, nil
}

func (h *GrpcHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	success, err := h.service.DeleteUser(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}
	return &pb.DeleteUserResponse{
		Success: success,
	}, nil
}

func (h *GrpcHandler) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	domainUsers, nextPageToken, err := h.service.ListUsers(int(req.GetLimit()), req.GetPageToken())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list users: %v", err)
	}

	var users []*pb.User

	for _, u := range domainUsers {
		users = append(users, &pb.User{
			Username:  u.Username,
			Email:     u.Email,
			Language:  u.Language,
			CreatedAt: u.CreatedAt.Unix(),
		})
	}

	return &pb.ListUsersResponse{
		Users:         users,
		NextPageToken: nextPageToken,
	}, nil
}
