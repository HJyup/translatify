package gateway

import (
	"context"
	pb "github.com/HJyup/translatify-common/api"
)

type ChatGateway interface {
	CreateConversation(ctx context.Context, payload *pb.CreateConversationRequest) (*pb.CreateConversationResponse, error)
	SendMessage(ctx context.Context, payload *pb.SendMessageRequest) (*pb.SendMessageResponse, error)
	GetMessage(context.Context, *pb.GetMessageRequest) (*pb.GetMessageResponse, error)
	ListMessages(context.Context, *pb.ListMessagesRequest) (*pb.ListMessagesResponse, error)
	StreamMessages(context.Context, *pb.StreamMessagesRequest) (pb.ChatService_StreamMessagesClient, error)
}
