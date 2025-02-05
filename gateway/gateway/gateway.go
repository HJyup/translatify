package gateway

import (
	"context"
	pb "github.com/HJyup/translatify-common/api"
)

type ChatGateway interface {
	GetMessage(context.Context, *pb.GetMessageRequest) (*pb.GetMessageResponse, error)
	ListMessages(context.Context, *pb.ListMessagesRequest) (*pb.ListMessagesResponse, error)
}
