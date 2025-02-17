package service

import (
	"context"
	"errors"
	"time"

	"github.com/HJyup/translatify-chat/internal/models"
)

type Service struct {
	store models.ChatStore
}

func NewService(store models.ChatStore) *Service {
	return &Service{store: store}
}

func (s *Service) CreateChat(userA, userB, sourceLang, targetLang string) (string, error) {
	if userA == "" || userB == "" || sourceLang == "" || targetLang == "" {
		return "", errors.New("usernameA, userNameB, sourceLanguage, and targetLanguage are required")
	}

	conv := &models.Chat{
		UsernameA:  userA,
		UsernameB:  userB,
		CreatedAt:  time.Now(),
		SourceLang: sourceLang,
		TargetLang: targetLang,
	}

	chatID, err := s.store.CreateConversion(context.Background(), conv)
	if err != nil {
		return "", err
	}

	return chatID, nil
}

func (s *Service) SendMessage(chatID, senderUsername, receiverUsername, content string) (string, error) {
	if chatID == "" || senderUsername == "" || receiverUsername == "" || content == "" {
		return "", errors.New("fromUsername, toUsername, and content are required")
	}

	now := time.Now()

	msg := &models.ChatMessage{
		ChatID:            chatID,
		SenderUsername:    senderUsername,
		ReceiverUsername:  receiverUsername,
		Content:           content,
		TranslatedContent: "",
		Timestamp:         now,
	}

	messageID, err := s.store.AddMessage(context.Background(), msg)
	if err != nil {
		return "", err
	}

	return messageID, nil
}

func (s *Service) GetMessage(messageID string) (*models.ChatMessage, error) {
	if messageID == "" {
		return nil, errors.New("message id is required")
	}
	return s.store.GetMessage(context.Background(), messageID)
}

func (s *Service) ListMessages(chatID string, since *time.Time, limit int, pageToken string) ([]*models.ChatMessage, string, error) {
	if chatID == "" {
		return nil, "", errors.New("chatID is required")
	}
	return s.store.ListMessages(context.Background(), chatID, since, limit, pageToken)
}

func (s *Service) StreamMessages(ctx context.Context, chatID string) (<-chan *models.ChatMessage, error) {
	if chatID == "" {
		return nil, errors.New("chatID is required")
	}

	out := make(chan *models.ChatMessage)
	startTime := time.Now()

	go func() {
		defer close(out)
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				messages, _, err := s.store.ListMessages(context.Background(), chatID, &startTime, 100, "")
				if err != nil {
					continue
				}

				for _, msg := range messages {
					select {
					case <-ctx.Done():
						return
					case out <- msg:
					}
				}

				startTime = time.Now()
			}
		}
	}()

	return out, nil
}

func (s *Service) GetChat(chatID string) (*models.Chat, error) {
	if chatID == "" {
		return nil, errors.New("chatID is required")
	}
	return s.store.GetChat(context.Background(), chatID)
}

func (s *Service) ListChats(userName string) ([]*models.Chat, error) {
	if userName == "" {
		return nil, errors.New("userName is required")
	}
	return s.store.ListChats(context.Background(), userName)
}
