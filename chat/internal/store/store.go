package store

import (
	"github.com/HJyup/translatify-chat/internal/models"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type Store struct {
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) AddMessage(msg *models.ChatMessage) error {
	return nil
}

func (s *Store) GetMessage(id string) (*models.ChatMessage, error) {
	return nil, nil
}

func (s *Store) ListMessages(userID, correspondentID string, since *timestamp.Timestamp) ([]*models.ChatMessage, error) {
	return nil, nil
}
