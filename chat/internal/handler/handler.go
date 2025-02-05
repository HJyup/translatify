package handler

import (
	"context"
	"github.com/HJyup/translatify-common/broker"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"

	"github.com/HJyup/translatify-chat/internal/models"
	pb "github.com/HJyup/translatify-common/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcHandler struct {
	pb.UnimplementedChatServiceServer

	service models.ChatService
	channel *amqp.Channel
}

func NewGrpcHandler(grpcServer *grpc.Server, service models.ChatService, channel *amqp.Channel) {
	handler := &GrpcHandler{
		service: service,
		channel: channel,
	}

	pb.RegisterChatServiceServer(grpcServer, handler)
}

func (h *GrpcHandler) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	fromID := req.GetFromUserId()
	toID := req.GetToUserId()
	content := req.GetContent()

	if fromID == "" || toID == "" || content == "" {
		return nil, status.Error(codes.InvalidArgument, "from_user_id, to_user_id, and content must be provided")
	}

	messageID, err := h.service.SendMessage(fromID, toID, content, "", "")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to send message: %v", err)
	}

	q, err := h.channel.QueueDeclare(broker.MessageSentEvent, true, false, false, false, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to declare queue: %v", err)
	}

	err = h.channel.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         []byte("25"),
		DeliveryMode: amqp.Persistent,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to publish message: %v", err)
	}

	return &pb.SendMessageResponse{MessageId: messageID}, nil
}

func (h *GrpcHandler) StreamMessages(req *pb.StreamMessagesRequest, stream pb.ChatService_StreamMessagesServer) error {
	userID := req.GetUserId()
	correspondentUserID := req.GetCorrespondentUserId()
	sinceTimestamp := req.GetSinceTimestamp()

	msgCh, err := h.service.StreamMessages(stream.Context(), userID, correspondentUserID, sinceTimestamp)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to start message stream: %v", err)
	}

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case msg, ok := <-msgCh:
			if !ok {
				return nil
			}
			if err = stream.Send(msg); err != nil {
				return status.Errorf(codes.Internal, "failed to send message: %v", err)
			}
		}
	}
}

func (h *GrpcHandler) GetMessage(ctx context.Context, req *pb.GetMessageRequest) (*pb.GetMessageResponse, error) {
	msg, err := h.service.GetMessage(req.GetMessageId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get message: %v", err)
	}

	return &pb.GetMessageResponse{Message: msg}, nil
}

func (h *GrpcHandler) ListMessages(ctx context.Context, req *pb.ListMessagesRequest) (*pb.ListMessagesResponse, error) {
	if req.GetUserId() == "" || req.GetCorrespondentUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and correspondent_user_id must be provided")
	}

	var since *timestamppb.Timestamp
	if req.GetSinceTimestamp() > 0 {
		since = timestamppb.New(time.Unix(req.GetSinceTimestamp(), 0))
	}

	messages, err := h.service.ListMessages(req.GetUserId(), req.GetCorrespondentUserId(), since)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list messages: %v", err)
	}

	resp := &pb.ListMessagesResponse{
		Messages:      messages,
		NextPageToken: "",
	}
	return resp, nil
}
