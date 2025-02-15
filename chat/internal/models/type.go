package models

import (
	"context"
	pb "github.com/HJyup/translatify-common/api"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type ChatService interface {
	CreateChat(userNameA, userNameB, sourceLang, targetLang string) (string, error)
	SendMessage(chatID, senderUserName, receiverUserName, content string) (string, error)
	GetMessage(messageID string) (*pb.ChatMessage, error)
	ListMessages(chatID string, since *timestamp.Timestamp, limit int, pageToken string) ([]*pb.ChatMessage, string, error)
	StreamMessages(ctx context.Context, chatID string) (<-chan *pb.ChatMessage, error)
	GetChat(chatID string) (*pb.Chat, error)
	ListChats(userName string) ([]*pb.Chat, error)
}

type ChatStore interface {
	CreateConversion(ctx context.Context, conv *pb.Chat) (string, error)
	AddMessage(ctx context.Context, msg *pb.ChatMessage) error
	GetMessage(ctx context.Context, id string) (*pb.ChatMessage, error)
	ListMessages(ctx context.Context, chatID string, since *timestamp.Timestamp, limit int, pageToken string) ([]*pb.ChatMessage, string, error)
	GetChat(ctx context.Context, id string) (*pb.Chat, error)
	ListChats(ctx context.Context, userName string) ([]*pb.Chat, error)
}
