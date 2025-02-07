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

func (s *Service) CreateConversation(userA, userB, sourceLang, targetLang string) (string, error) {
	if userA == "" || userB == "" || sourceLang == "" || targetLang == "" {
		return "", errors.New("userA, userB, sourceLang, and targetLang are required")
	}
	conv := &pb.Conversation{
		ConversationId: uuid.New().String(),
		UserAId:        userA,
		UserBId:        userB,
		CreatedAt:      time.Now().Unix(),
		SourceLanguage: sourceLang,
		TargetLanguage: targetLang,
	}

	conversationID, err := s.store.CreateConversion(context.Background(), conv)
	if err != nil {
		return "", err
	}

	return conversationID, nil
}

func (s *Service) SendMessage(conversationID, senderID, receiverID, content string) (string, error) {
	if conversationID == "" || senderID == "" || receiverID == "" || content == "" {
		return "", errors.New("conversationID, senderID, receiverID, and content are required")
	}

	messageID := uuid.New().String()
	now := time.Now().Unix()

	msg := &pb.ChatMessage{
		MessageId:         messageID,
		ConversationId:    conversationID,
		SenderId:          senderID,
		ReceiverId:        receiverID,
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

func (s *Service) GetMessage(messageID string) (*pb.ChatMessage, error) {
	if messageID == "" {
		return nil, errors.New("message id is required")
	}
	return s.store.GetMessage(context.Background(), messageID)
}

func (s *Service) ListMessages(conversationID string, since *timestamp.Timestamp, limit int, pageToken string) ([]*pb.ChatMessage, string, error) {
	if conversationID == "" {
		return nil, "", errors.New("conversationID is required")
	}
	return s.store.ListMessages(context.Background(), conversationID, since, limit, pageToken)
}

func (s *Service) StreamMessages(ctx context.Context, conversationID string) (<-chan *pb.ChatMessage, error) {
	if conversationID == "" {
		return nil, errors.New("conversationID is required")
	}

	out := make(chan *pb.ChatMessage)
	startTime := time.Now().Unix()

	go func() {
		defer close(out)
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				since := &timestamp.Timestamp{Seconds: startTime}
				messages, _, err := s.store.ListMessages(context.Background(), conversationID, since, 100, "")
				if err != nil {
					// In production, you might want to log this error.
					continue
				}

				for _, msg := range messages {
					select {
					case <-ctx.Done():
						return
					case out <- msg:
					}
				}

				startTime = time.Now().Unix()
			}
		}
	}()

	return out, nil
}

func (s *Service) GetConversation(conversationID string) (*pb.Conversation, error) {
	if conversationID == "" {
		return nil, errors.New("conversationID is required")
	}
	return s.store.GetConversation(context.Background(), conversationID)
}

func (s *Service) ListConversations(userID string) ([]*pb.Conversation, error) {
	if userID == "" {
		return nil, errors.New("userID is required")
	}
	return s.store.ListConversations(context.Background(), userID)
}
