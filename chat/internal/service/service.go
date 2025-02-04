package service

import (
	"github.com/HJyup/translatify-chat/internal/models"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type Service struct {
	store models.ChatStore
}

func NewService(store models.ChatStore) *Service {
	return &Service{store: store}
}

func (s *Service) SendMessage(fromID, toID, content, sourceLang, targetLang string) (string, error) {
	return "", nil
}

func (s *Service) GetMessage(id string) (*models.ChatMessage, error) {
	return nil, nil
}

func (s *Service) ListMessages(userID, correspondentID string, since *timestamp.Timestamp) ([]*models.ChatMessage, error) {
	return nil, nil
}
