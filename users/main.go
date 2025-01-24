package main

import (
	"context"
	common "github.com/HJyup/translatify-common"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	grpcAddr = common.EnvString("GRPC_ADDR", "localhost:50051")
)

func main() {
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer listener.Close()

	store := NewStore()
	svc := NewService(store)
	NewGRPCServer(grpcServer)

	svc.CreateUser(context.Background())

	log.Printf("gRPC server is running on %s", grpcAddr)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err.Error())
	}
}
