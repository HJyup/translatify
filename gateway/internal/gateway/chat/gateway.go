package chat

import (
	"context"
	pb "github.com/HJyup/translatify-common/api"
)

type Gateway interface {
	CreateChat(ctx context.Context, payload *pb.CreateChatRequest) (*pb.CreateChatResponse, error)
	SendMessage(ctx context.Context, payload *pb.SendMessageRequest) (*pb.SendMessageResponse, error)
	GetMessage(context.Context, *pb.GetMessageRequest) (*pb.GetMessageResponse, error)
	ListMessages(context.Context, *pb.ListMessagesRequest) (*pb.ListMessagesResponse, error)
	StreamMessages(context.Context, *pb.StreamMessagesRequest) (pb.ChatService_StreamMessagesClient, error)
	GetChat(context.Context, *pb.GetChatRequest) (*pb.GetChatResponse, error)
	ListChats(context.Context, *pb.ListChatsRequest) (*pb.ListChatsResponse, error)
}
