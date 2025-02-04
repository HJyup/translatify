package main

import (
	"github.com/HJyup/translatify-chat/internal/handler"
	"github.com/HJyup/translatify-chat/internal/service"
	"github.com/HJyup/translatify-chat/internal/store"
	common "github.com/HJyup/translatify-common"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	grpcAddr = common.EnvString("GRPC_ADDR", "localhost:5050")
)

func main() {
	grpcServer := grpc.NewServer()
	conn, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer conn.Close()

	str := store.NewStore()
	srv := service.NewService(str)
	handler.NewGrpcHandler(grpcServer)

	srv.GetMessage("25")

	log.Printf("Starting chat server on %s", grpcAddr)

	if err := grpcServer.Serve(conn); err != nil {
		log.Fatal(err.Error())
	}
}
