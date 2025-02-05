package service

import (
	"context"
	"errors"
	"time"

	"github.com/HJyup/translatify-chat/internal/models"
	pb "github.com/HJyup/translatify-common/api"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
)

type Service struct {
	store models.ChatStore
}

func NewService(store models.ChatStore) *Service {
	return &Service{store: store}
}

func (s *Service) SendMessage(fromID, toID, content, sourceLang, targetLang string) (string, error) {
	if fromID == "" || toID == "" || content == "" {
		return "", errors.New("fromID, toID, and content are required")
	}

	messageID := uuid.New().String()
	now := time.Now().Unix()

	msg := &pb.ChatMessage{
		MessageId:         messageID,
		FromUserId:        fromID,
		ToUserId:          toID,
		Content:           content,
		TranslatedContent: "",
		Timestamp:         now,
		Translated:        false,
	}

	if err := s.store.AddMessage(context.Background(), msg); err != nil {
		return "", err
	}

	return messageID, nil
}

func (s *Service) GetMessage(id string) (*pb.ChatMessage, error) {
	if id == "" {
		return nil, errors.New("message id is required")
	}

	return s.store.GetMessage(context.Background(), id)
}

func (s *Service) ListMessages(userID, correspondentID string, since *timestamp.Timestamp) ([]*pb.ChatMessage, error) {
	if userID == "" || correspondentID == "" {
		return nil, errors.New("userID and correspondentID are required")
	}

	messages, err := s.store.ListMessages(context.Background(), userID, correspondentID, since)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (s *Service) StreamMessages(ctx context.Context, userID, correspondentID string, sinceTimestamp int64) (<-chan *pb.ChatMessage, error) {
	if userID == "" || correspondentID == "" {
		return nil, errors.New("userID and correspondentID are required")
	}

	out := make(chan *pb.ChatMessage)

	go func() {
		defer close(out)
		currentSince := sinceTimestamp

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				since := &timestamp.Timestamp{Seconds: currentSince}
				messages, err := s.store.ListMessages(context.Background(), userID, correspondentID, since)
				if err != nil {
					// In production, you might want to log this error.
					continue
				}

				for _, msg := range messages {
					// Only send messages that are new.
					if msg.Timestamp > currentSince {
						select {
						case <-ctx.Done():
							return
						case out <- msg:
							currentSince = msg.Timestamp
						}
					}
				}
			}
		}
	}()

	return out, nil
}
