package main

import (
	"fmt"
	"net"

	"github.com/samantonio28/vk-task-subscribe/internal/config"
	subpubpb "github.com/samantonio28/vk-task-subscribe/internal/delivery"
	subpub "github.com/samantonio28/vk-task-subscribe/internal/repository"
	"github.com/samantonio28/vk-task-subscribe/internal/usecase"
	"github.com/samantonio28/vk-task-subscribe/logger"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}

	logger, err := logger.NewLogrusLogger(cfg.Logging.FilePath)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		return
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
	if err != nil {
		fmt.Printf("Failed to listen port: %v\n", err)
		return
	}

	server := grpc.NewServer(
		grpc.ConnectionTimeout(cfg.GRPC.Timeout),
	)

	subpub := subpub.NewSubPub()

	subPubService, err := usecase.NewSubPubService(logger, subpub)
	if err != nil {
		fmt.Printf("Failed to initialize subPubService: %v\n", err)
		return
	}

	subpubpb.RegisterPubSubServer(server, subPubService)

	fmt.Printf("Starting server at :%d\n", cfg.GRPC.Port)
	if err := server.Serve(listener); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
