package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/HJyup/translatify-chat/internal/handler"
	"github.com/HJyup/translatify-chat/internal/service"
	"github.com/HJyup/translatify-chat/internal/store"
	"github.com/HJyup/translatify-common/broker"
	"github.com/HJyup/translatify-common/discovery"
	"github.com/HJyup/translatify-common/discovery/consul"
	common "github.com/HJyup/translatify-common/utils"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"

	_ "github.com/joho/godotenv/autoload"
)

var (
	serviceName = common.EnvString("SERVICE_NAME")
	grpcAddr    = common.EnvString("GRPC_ADDR")
	consulAddr  = common.EnvString("CONSUL_ADDR")

	amqpUser    = common.EnvString("AMQP_USER")
	amqpPass    = common.EnvString("AMQP_PASS")
	amqpHost    = common.EnvString("AMQP_HOST")
	amqpPort    = common.EnvString("AMQP_PORT")
	databaseURL = common.EnvString("DATABASE_URL")
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

	go func() {
		for {
			if err = registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Fatal("Failed to health check", err)
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.DeRegister(ctx, instanceID, serviceName)

	ch, closeConn := broker.Connect(amqpUser, amqpPass, amqpHost, amqpPort)
	defer func() {
		_ = closeConn()
		_ = ch.Close()
	}()

	config, err := pgx.ParseConfig(databaseURL)
	if err != nil {
		log.Fatalf("Failed to parse database config: %v", err)
	}
	config.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	dbConn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer dbConn.Close(context.Background())

	grpcServer := grpc.NewServer()
	conn, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer conn.Close()

	str := store.NewStore(dbConn)
	srv := service.NewService(str)
	handler.NewGrpcHandler(grpcServer, srv, ch)

	log.Printf("Starting chat server on %s", grpcAddr)
	if err = grpcServer.Serve(conn); err != nil {
		log.Fatal(err.Error())
	}
}
