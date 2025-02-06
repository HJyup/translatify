package models

import (
	"context"
	pb "github.com/HJyup/translatify-common/api"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type ChatService interface {
	CreateConversation(userAID, userBID, sourceLang, targetLang string) (string, error)
	SendMessage(conversationID, senderID, receiverID, content string) (string, error)
	GetMessage(messageID string) (*pb.ChatMessage, error)
	ListMessages(conversationID string, since *timestamp.Timestamp, limit int, pageToken string) ([]*pb.ChatMessage, string, error)
	StreamMessages(ctx context.Context, conversationID string) (<-chan *pb.ChatMessage, error)
}

type ChatStore interface {
	CreateConversion(ctx context.Context, conv *pb.Conversation) (string, error)
	AddMessage(ctx context.Context, msg *pb.ChatMessage) error
	GetMessage(ctx context.Context, id string) (*pb.ChatMessage, error)
	ListMessages(ctx context.Context, conversationID string, since *timestamp.Timestamp, limit int, pageToken string) ([]*pb.ChatMessage, string, error)
}
