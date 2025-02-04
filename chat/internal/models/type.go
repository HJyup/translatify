package models

import "github.com/golang/protobuf/ptypes/timestamp"

type ChatMessage struct {
	ID                string
	FromID            string
	ToID              string
	Content           string
	TranslatedContent string
	Timestamp         *timestamp.Timestamp
	Translated        bool
}

type ChatService interface {
	SendMessage(fromID, toID, content, sourceLang, targetLang string) (string, error)
	GetMessage(id string) (*ChatMessage, error)
	ListMessages(userID, correspondentID string, since *timestamp.Timestamp) ([]*ChatMessage, error)
}

type ChatStore interface {
	AddMessage(msg *ChatMessage) error
	GetMessage(id string) (*ChatMessage, error)
	ListMessages(userID, correspondentID string, since *timestamp.Timestamp) ([]*ChatMessage, error)
}
