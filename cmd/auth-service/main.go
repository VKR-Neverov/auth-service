package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fidesy-pay/auth-service/internal/app"
	"github.com/fidesy-pay/auth-service/internal/config"
	authservice "github.com/fidesy-pay/auth-service/internal/pkg/auth-service"
	clientsconsumer "github.com/fidesy-pay/auth-service/internal/pkg/consumers/client-consumer"
	"github.com/fidesy-pay/auth-service/internal/pkg/storage"
	"github.com/fidesy/sdk/common/grpc"
	"github.com/fidesy/sdk/common/kafka"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	defer cancel()

	err := config.Init()
	if err != nil {
		log.Fatalf("config.Init: %v", err)
	}

	redisClient, err := storage.NewRedis(ctx)
	if err != nil {
		log.Fatalf("storage.NewRedis: %v", err)
	}

	store := storage.New(redisClient)

	authService := authservice.New(store)

	impl := app.New(authService)

	server, err := grpc.NewServer(
		grpc.WithPort(os.Getenv("GRPC_PORT")),
		grpc.WithMetricsPort(os.Getenv("METRICS_PORT")),
		grpc.WithDomainNameService(ctx, "domain-name-service:10000"),
		grpc.WithGraylog("graylog:5555"),
		grpc.WithTracer("http://jaeger:14268/api/traces"),
	)
	if err != nil {
		log.Fatalf("grpc.NewServer: %v", err)
	}

	err = kafka.RegisterConsumer(
		ctx,
		clientsconsumer.New(store),
		config.Get(config.KafkaBrokers).([]string),
		"clients-json",
	)
	if err != nil {
		log.Fatalf("kafka.RegisterConsumer: %v", err)
	}

	err = server.Run(ctx, impl)
	if err != nil {
		log.Fatalf("app.Run: %v", err)
	}
}
