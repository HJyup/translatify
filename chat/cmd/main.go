package main

import (
	"context"
	"github.com/HJyup/translatify-chat/internal/handler"
	"github.com/HJyup/translatify-chat/internal/service"
	"github.com/HJyup/translatify-chat/internal/store"
	"github.com/HJyup/translatify-common/discovery"
	"github.com/HJyup/translatify-common/discovery/consul"
	common "github.com/HJyup/translatify-common/utils"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

var (
	serviceName = "chat"
	grpcAddr    = common.EnvString("GRPC_ADDR", "localhost:5050")
	consulAddr  = common.EnvString("CONSUL_ADDR", "localhost:8500")
)

func main() {
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Fatal("Failed to health check", err)
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.DeRegister(ctx, instanceID, serviceName)

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
