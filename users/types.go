package main

import "context"

type UsersService interface {
	CreateUser(ctx context.Context) error
}

type UsersStore interface {
	Create(ctx context.Context) error
}
