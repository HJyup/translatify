package service

import (
	"context"
	"time"
)

type UserService interface {
	CreateUser(username, email, password, language string) (string, error)
	GetUser(username string) (*User, error)
	DeleteUser(userId string) (bool, error)
	ListUsers(limit int, paginationToken string) ([]*User, string, error)
}

type UserStore interface {
	CreateUser(ctx context.Context, username, email, password, language string) (*User, error)
	GetUser(ctx context.Context, username string) (*User, error)
	DeleteUser(ctx context.Context, userId string) (bool, error)
	ListUsers(ctx context.Context, limit int, paginationToken string) ([]*User, string, error)
}

type User struct {
	UserId    string
	Username  string
	Email     string
	Password  string
	Language  string
	CreatedAt time.Time
}
