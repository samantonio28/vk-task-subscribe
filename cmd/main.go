package main

import (
	"fmt"
	"net"

	subpubpb "github.com/samantonio28/vk-task-subscribe/internal/delivery"
	impl "github.com/samantonio28/vk-task-subscribe/internal/usecase"
	logger "github.com/samantonio28/vk-task-subscribe/logger"
	"google.golang.org/grpc"
)

func main() {
	logger, err := logger.NewLogrusLogger("./logs/access.log")
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		return
	}

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(fmt.Errorf("not able to listen port: %v", err))
		return
	}

	server := grpc.NewServer()

	subPubService := &impl.SubPubService{
		Logger: logger,
	}

	subpubpb.RegisterPubSubServer(server, subPubService)

	fmt.Println("starting server at :8080")
	server.Serve(listener)
}
