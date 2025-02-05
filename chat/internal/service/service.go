package service

import (
	"errors"
	"github.com/HJyup/translatify-chat/internal/models"
	pb "github.com/HJyup/translatify-common/api"
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

func (s *Service) GetMessage(id string) (*pb.ChatMessage, error) {
	return nil, nil
}

func (s *Service) ListMessages(userID, correspondentID string, since *timestamp.Timestamp) ([]*pb.ChatMessage, error) {
	if userID == "" || correspondentID == "" {
		return nil, errors.New("userID and correspondentID are required")
	}

	messages, err := s.store.ListMessages(userID, correspondentID, since)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
