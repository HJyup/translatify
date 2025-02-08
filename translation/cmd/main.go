package main

import (
	"context"
	"fmt"
	"github.com/HJyup/translatify-common/broker"
	"github.com/HJyup/translatify-common/discovery"
	"github.com/HJyup/translatify-common/discovery/consul"
	common "github.com/HJyup/translatify-common/utils"
	"github.com/HJyup/translatify-translation/internal/consumer"
	"github.com/HJyup/translatify-translation/internal/handler"
	"github.com/HJyup/translatify-translation/internal/service"
	"github.com/HJyup/translatify-translation/internal/translator"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

var (
	serviceName = common.EnvString("SERVICE_NAME")
	grpcAddr    = common.EnvString("GRPC_ADDR")
	consulAddr  = common.EnvString("CONSUL_ADDR")

	openaiAPIKey = common.EnvString("OPENAI_API_KEY")

	amqpUser = common.EnvString("AMQP_USER")
	amqpPass = common.EnvString("AMQP_PASS")
	amqpHost = common.EnvString("AMQP_HOST")
	amqpPort = common.EnvString("AMQP_PORT")
)

func main() {
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err = registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
		panic(err)
	}

	ch, closeConn := broker.Connect(amqpUser, amqpPass, amqpHost, amqpPort)
	defer func() {
		_ = closeConn()
		_ = ch.Close()
	}()

	go func() {
		for {
			if err = registry.HealthCheck(instanceID, serviceName); err != nil {
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

	trans := translator.NewModel(openaiAPIKey)
	srv := service.NewTranslationService(trans)

	cons := consumer.NewConsumer(srv)
	go cons.Listen(ch)

	handler.NewGrpcHandler(grpcServer, srv, ch)

	fmt.Println("Starting gRPC server", grpcAddr)

	if err = grpcServer.Serve(conn); err != nil {
		log.Fatal(err)
	}

}
