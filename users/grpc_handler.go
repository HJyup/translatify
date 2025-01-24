package main

import (
	"context"
	pb "github.com/HJyup/translatify-common/api"
	"google.golang.org/grpc"
	"log"
)

type GrpcHandler struct {
	pb.UnimplementedUserServiceServer
}

func NewGRPCServer(grpcServer *grpc.Server) {
	handler := &GrpcHandler{}
	pb.RegisterUserServiceServer(grpcServer, handler)
}

func (h *GrpcHandler) CreateUser(_ context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	log.Println("CreateUser", req)
	userResponse := &pb.CreateUserResponse{
		UserId:     "1",
		Message:    "User created",
		StatusCode: 200,
	}

	return userResponse, nil
}
func (h *GrpcHandler) GetUser(_ context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	log.Println("GetUser", req)

	user := &pb.User{
		UserId: "1",
		Name:   "John Doe",
		Email:  "john@gmail.com",
	}

	return user, nil
}
