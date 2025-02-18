package models

import (
	"context"
	"time"
)

type ChatService interface {
	CreateChat(userNameA, userNameB, sourceLang, targetLang string) (string, error)
	SendMessage(ctx context.Context, chatID, senderUserName, receiverUserName, content string) (string, error)
	GetMessage(messageID string) (*ChatMessage, error)
	ListMessages(chatID string, since *time.Time, limit int, pageToken string) ([]*ChatMessage, string, error)
	StreamMessages(ctx context.Context, chatID string) (<-chan *ChatMessage, error)
	GetChat(chatID string) (*Chat, error)
	ListChats(userName string) ([]*Chat, error)
	UpdateMessageTranslation(messageID string, translatedContent string) error
}

type ChatStore interface {
	CreateConversion(ctx context.Context, conv *Chat) (string, error)
	AddMessage(ctx context.Context, msg *ChatMessage) (string, error)
	GetMessage(ctx context.Context, id string) (*ChatMessage, error)
	ListMessages(ctx context.Context, chatID string, since *time.Time, limit int, pageToken string) ([]*ChatMessage, string, error)
	GetChat(ctx context.Context, id string) (*Chat, error)
	ListChats(ctx context.Context, userName string) ([]*Chat, error)
	UpdateMessageTranslation(ctx context.Context, messageID string, translatedContent string) error
}

type ChatMessage struct {
	MessageID         string
	ChatID            string
	SenderUsername    string
	ReceiverUsername  string
	Content           string
	TranslatedContent string
	Timestamp         time.Time
}
type Chat struct {
	ChatID     string
	UsernameA  string
	UsernameB  string
	CreatedAt  time.Time
	SourceLang string
	TargetLang string
}

type ConsumerResponse struct {
	MessageId         string `json:"messageId"`
	TranslatedContent string `json:"translatedContent"`
	Success           bool   `json:"Success"`
}
