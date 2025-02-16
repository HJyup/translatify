package user

import (
	"context"
	pb "github.com/HJyup/translatify-common/api"
)

type Gateway interface {
	CreateUser(ctx context.Context, payload *pb.CreateUserRequest) (*pb.CreateUserResponse, error)
	GetUser(ctx context.Context, payload *pb.GetUserRequest) (*pb.GetUserResponse, error)
	UpdateUser(ctx context.Context, payload *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error)
	ListUsers(context.Context, *pb.ListUsersRequest) (*pb.ListUsersResponse, error)
	DeleteUser(context.Context, *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error)
}
