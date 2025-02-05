package models

import (
	pb "github.com/HJyup/translatify-common/api"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type ChatService interface {
	SendMessage(fromID, toID, content, sourceLang, targetLang string) (string, error)
	GetMessage(id string) (*pb.ChatMessage, error)
	ListMessages(userID, correspondentID string, since *timestamp.Timestamp) ([]*pb.ChatMessage, error)
}

type ChatStore interface {
	AddMessage(msg *pb.ChatMessage) error
	GetMessage(id string) (*pb.ChatMessage, error)
	ListMessages(userID, correspondentID string, since *timestamp.Timestamp) ([]*pb.ChatMessage, error)
}
