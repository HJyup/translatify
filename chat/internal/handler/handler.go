package handler

import (
	"context"
	pb "github.com/HJyup/translatify-common/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type GrpcHandler struct {
	pb.UnimplementedChatServiceServer
}

func NewGrpcHandler(grpcServer *grpc.Server) {
	handler := &GrpcHandler{}
	pb.RegisterChatServiceServer(grpcServer, handler)
}

func (h *GrpcHandler) SendMessage(context.Context, *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}

func (h *GrpcHandler) StreamMessages(*pb.StreamMessagesRequest, grpc.ServerStreamingServer[pb.ChatMessage]) error {
	return status.Errorf(codes.Unimplemented, "method StreamMessages not implemented")
}

func (h *GrpcHandler) GetMessage(context.Context, *pb.GetMessageRequest) (*pb.GetMessageResponse, error) {
	log.Println("Are u trying to get a message?")
	msg := &pb.ChatMessage{
		MessageId:         "25",
		FromUserId:        "22",
		ToUserId:          "34",
		Content:           "Hello, World!",
		TranslatedContent: "Hallo, Welt!",
		Timestamp:         234234,
		Translated:        false,
	}
	return &pb.GetMessageResponse{Message: msg}, nil
}
func (h *GrpcHandler) ListMessages(context.Context, *pb.ListMessagesRequest) (*pb.ListMessagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListMessages not implemented")
}
