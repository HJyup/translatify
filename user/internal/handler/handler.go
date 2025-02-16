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
	user, err := h.service.CreateUser(
		req.GetUsername(),
		req.GetEmail(),
		req.GetFullName(),
		req.GetPassword(),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}
	return &pb.CreateUserResponse{
		User: user,
	}, nil
}

func (h *GrpcHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := h.service.GetUser(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}
	return &pb.GetUserResponse{
		User: user,
	}, nil
}

func (h *GrpcHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	user, err := h.service.UpdateUser(
		req.GetUsername(),
		req.GetEmail(),
		req.GetFullName(),
		req.GetPassword(),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}
	return &pb.UpdateUserResponse{
		User: user,
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
	users, nextPageToken, err := h.service.ListUsers(int(req.GetLimit()), req.GetPageToken())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list users: %v", err)
	}
	return &pb.ListUsersResponse{
		Users:         users,
		NextPageToken: nextPageToken,
	}, nil
}
