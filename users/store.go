package main

import "context"

type Store struct {
	// add postgres connection
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) Create(ctx context.Context) error {
	return nil
}
