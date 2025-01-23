package main

import "context"

type Service struct {
	store UsersStore
}

func NewService(store UsersStore) *Service {
	return &Service{store}
}

func (s *Service) CreateUser(ctx context.Context) error {
	return nil
}
