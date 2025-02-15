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

func (s *Service) CreateChat(userA, userB, sourceLang, targetLang string) (string, error) {
	if userA == "" || userB == "" || sourceLang == "" || targetLang == "" {
		return "", errors.New("userA, userB, sourceLang, and targetLang are required")
	}
	conv := &pb.Chat{
		ChatId:         uuid.New().String(),
		UsernameA:      userA,
		UsernameB:      userB,
		CreatedAt:      time.Now().Unix(),
		SourceLanguage: sourceLang,
		TargetLanguage: targetLang,
	}

	chatID, err := s.store.CreateConversion(context.Background(), conv)
	if err != nil {
		return "", err
	}

	return chatID, nil
}

func (s *Service) SendMessage(chatID, senderUsername, receiverName, content string) (string, error) {
	if chatID == "" || senderUsername == "" || receiverName == "" || content == "" {
		return "", errors.New("chatID, senderID, receiverID, and content are required")
	}

	messageID := uuid.New().String()
	now := time.Now().Unix()

	msg := &pb.ChatMessage{
		MessageId:         messageID,
		ChatId:            chatID,
		SenderUsername:    senderUsername,
		ReceiverUsername:  receiverName,
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

func (s *Service) ListMessages(chatID string, since *timestamp.Timestamp, limit int, pageToken string) ([]*pb.ChatMessage, string, error) {
	if chatID == "" {
		return nil, "", errors.New("chatID is required")
	}
	return s.store.ListMessages(context.Background(), chatID, since, limit, pageToken)
}

func (s *Service) StreamMessages(ctx context.Context, chatID string) (<-chan *pb.ChatMessage, error) {
	if chatID == "" {
		return nil, errors.New("chatID is required")
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
				messages, _, err := s.store.ListMessages(context.Background(), chatID, since, 100, "")
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

func (s *Service) GetChat(chatID string) (*pb.Chat, error) {
	if chatID == "" {
		return nil, errors.New("chatID is required")
	}
	return s.store.GetChat(context.Background(), chatID)
}

func (s *Service) ListChats(userID string) ([]*pb.Chat, error) {
	if userID == "" {
		return nil, errors.New("userID is required")
	}
	return s.store.ListChats(context.Background(), userID)
}
