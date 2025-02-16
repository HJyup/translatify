package service

import (
	"context"
	"errors"
	models "github.com/HJyup/translatify-user/internal/model"

	"github.com/HJyup/translatify-common/utils"
	"golang.org/x/crypto/bcrypt"
)

const MinPasswordLength = 8

type UserService struct {
	store models.UserStore
}

func NewService(store models.UserStore) *UserService {
	return &UserService{store: store}
}

func (s *UserService) CreateUser(username, email, password string) (string, error) {
	ctx := context.Background()

	if username == "" || email == "" || password == "" {
		return "", errors.New("username, email, and password are required")
	}

	if len(password) < MinPasswordLength {
		return "", errors.New("password must be at least 8 characters long")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("failed to hash password")
	}

	user, err := s.store.CreateUser(ctx, username, email, string(hashedPassword))
	if err != nil {
		return "", err
	}

	token, err := utils.CreateToken(user.UserId, user.Username, user.Email)
	if err != nil {
		return "", errors.New("failed to generate token: " + err.Error())
	}

	return token, nil
}

func (s *UserService) GetUser(username string) (*models.User, error) {
	ctx := context.Background()

	if username == "" {
		return nil, errors.New("userID is required")
	}

	user, err := s.store.GetUser(ctx, username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) DeleteUser(userId string) (bool, error) {
	ctx := context.Background()

	if userId == "" {
		return false, errors.New("userID is required")
	}

	_, err := s.store.DeleteUser(ctx, userId)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *UserService) ListUsers(limit int, paginationToken string) ([]*models.User, string, error) {
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
