package models

import (
	"context"
	pb "github.com/HJyup/translatify-common/api"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type ChatService interface {
	SendMessage(fromID, toID, content, sourceLang, targetLang string) (string, error)
	GetMessage(id string) (*pb.ChatMessage, error)
	ListMessages(userID, correspondentID string, since *timestamp.Timestamp) ([]*pb.ChatMessage, error)
	StreamMessages(ctx context.Context, userID, correspondentID string, sinceTimestamp int64) (<-chan *pb.ChatMessage, error)
}

type ChatStore interface {
	AddMessage(ctx context.Context, msg *pb.ChatMessage) error
	GetMessage(ctx context.Context, id string) (*pb.ChatMessage, error)
	ListMessages(ctx context.Context, userID, correspondentID string, since *timestamp.Timestamp) ([]*pb.ChatMessage, error)
}
