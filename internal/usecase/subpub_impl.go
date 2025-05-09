package usecase

import (
	"sync"

	subpub "github.com/samantonio28/vk-task-subscribe/internal/delivery"
	repository "github.com/samantonio28/vk-task-subscribe/internal/repository"
	"github.com/samantonio28/vk-task-subscribe/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SubPubService struct {
	subpub.UnimplementedPubSubServer
	subPub repository.SubPub
	Logger *logger.LogrusLogger
	mu     sync.Mutex
	subs   map[string]map[subpub.PubSub_SubscribeServer]struct{}
}

func NewSubPubService(logger *logger.LogrusLogger, subPub repository.SubPub) (*SubPubService, error) {
    if logger == nil || subPub == nil {
        return nil, status.Error(codes.InvalidArgument, "logger and subPub must not be nil")
    }

    return &SubPubService{
        subPub: subPub,
        Logger: logger,
        subs:   make(map[string]map[subpub.PubSub_SubscribeServer]struct{}), // Инициализация основной map
    }, nil
}
