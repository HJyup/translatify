package service

import (
	"context"
	"errors"

	pb "github.com/HJyup/translatify-common/api"
	models "github.com/HJyup/translatify-user/internal/model"
	"golang.org/x/crypto/bcrypt"
)

const MinPasswordLength = 8

type Service struct {
	store models.UserStore
}

func NewService(store models.UserStore) *Service {
	return &Service{store: store}
}

func (s *Service) CreateUser(username, email, fullName, password string) (*pb.User, error) {
	ctx := context.Background()

	if username == "" || email == "" || fullName == "" || password == "" {
		return nil, errors.New("username, email, fullName, and password are required")
	}

	if len(password) < MinPasswordLength {
		return nil, errors.New("password must be at least 8 characters long")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user, err := s.store.CreateUser(ctx, username, email, fullName, string(hashedPassword))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GetUser(userId string) (*pb.User, error) {
	ctx := context.Background()

	if userId == "" {
		return nil, errors.New("userId is required")
	}

	user, err := s.store.GetUser(ctx, userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) UpdateUser(username, email, fullName, password string) (*pb.User, error) {
	ctx := context.Background()

	if username == "" || email == "" || fullName == "" {
		return nil, errors.New("username, email, and fullName are required")
	}

	var hashedPassword string
	if password != "" {
		if len(password) < MinPasswordLength {
			return nil, errors.New("password must be at least 8 characters long")
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("failed to hash password")
		}
		hashedPassword = string(hash)
	}

	user, err := s.store.UpdateUser(ctx, username, email, fullName, hashedPassword)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) DeleteUser(userId string) (bool, error) {
	ctx := context.Background()

	if userId == "" {
		return false, errors.New("userId is required")
	}

	success, err := s.store.DeleteUser(ctx, userId)
	if err != nil {
		return false, err
	}

	return success, nil
}

func (s *Service) ListUsers(limit int, paginationToken string) ([]*pb.User, string, error) {
	ctx := context.Background()

	if limit <= 0 {
		limit = 10
	}

	users, token, err := s.store.ListUsers(ctx, limit, paginationToken)
	if err != nil {
		return nil, "", err
	}

	return users, token, nil
}
