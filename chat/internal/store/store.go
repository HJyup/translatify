package store

import (
	pb "github.com/HJyup/translatify-common/api"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type Store struct {
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) AddMessage(msg *pb.ChatMessage) error {
	return nil
}

func (s *Store) GetMessage(id string) (*pb.ChatMessage, error) {
	return nil, nil
}

func (s *Store) ListMessages(userID, correspondentID string, since *timestamp.Timestamp) ([]*pb.ChatMessage, error) {
	return nil, nil
}
