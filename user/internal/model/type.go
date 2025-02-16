package models

import (
	"context"
	pb "github.com/HJyup/translatify-common/api"
)

type UserService interface {
	CreateUser(username, email, fullName, password string) (*pb.User, error)
	GetUser(userId string) (*pb.User, error)
	UpdateUser(username, email, fullName, password string) (*pb.User, error)
	DeleteUser(userId string) (bool, error)
	ListUsers(limit int, paginationToken string) ([]*pb.User, string, error)
}

type UserStore interface {
	CreateUser(ctx context.Context, username, email, fullName, password string) (*pb.User, error)
	GetUser(ctx context.Context, userId string) (*pb.User, error)
	UpdateUser(ctx context.Context, username, email, fullName, password string) (*pb.User, error)
	DeleteUser(ctx context.Context, userId string) (bool, error)
	ListUsers(ctx context.Context, limit int, paginationToken string) ([]*pb.User, string, error)
}
