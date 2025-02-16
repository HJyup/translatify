package user

import (
	"context"
	"log"

	pb "github.com/HJyup/translatify-common/api"
	"github.com/HJyup/translatify-common/discovery"
)

type GrpcGateway struct {
	registry discovery.Registry
}

func NewGateway(registry discovery.Registry) *GrpcGateway {
	return &GrpcGateway{registry: registry}
}

func (g *GrpcGateway) CreateUser(ctx context.Context, payload *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	conn, err := discovery.ServiceConnection(ctx, "user", g.registry)
	if err != nil {
		log.Fatal("Failed to connect to user service")
	}
	userClient := pb.NewUserServiceClient(conn)
	return userClient.CreateUser(ctx, payload)
}

func (g *GrpcGateway) GetUser(ctx context.Context, payload *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	conn, err := discovery.ServiceConnection(ctx, "user", g.registry)
	if err != nil {
		log.Fatal("Failed to connect to user service")
	}
	userClient := pb.NewUserServiceClient(conn)
	return userClient.GetUser(ctx, payload)
}

func (g *GrpcGateway) UpdateUser(ctx context.Context, payload *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	conn, err := discovery.ServiceConnection(ctx, "user", g.registry)
	if err != nil {
		log.Fatal("Failed to connect to user service")
	}
	userClient := pb.NewUserServiceClient(conn)
	return userClient.UpdateUser(ctx, payload)
}

func (g *GrpcGateway) ListUsers(ctx context.Context, payload *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	conn, err := discovery.ServiceConnection(ctx, "user", g.registry)
	if err != nil {
		log.Fatal("Failed to connect to user service")
	}
	userClient := pb.NewUserServiceClient(conn)
	return userClient.ListUsers(ctx, payload)
}

func (g *GrpcGateway) DeleteUser(ctx context.Context, payload *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	conn, err := discovery.ServiceConnection(ctx, "user", g.registry)
	if err != nil {
		log.Fatal("Failed to connect to user service")
	}
	userClient := pb.NewUserServiceClient(conn)
	return userClient.DeleteUser(ctx, payload)
}
