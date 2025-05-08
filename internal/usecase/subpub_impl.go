package usecase

import (
	"fmt"

	subpub "github.com/samantonio28/vk-task-subscribe/internal/delivery"
	repository "github.com/samantonio28/vk-task-subscribe/internal/repository"
	"github.com/samantonio28/vk-task-subscribe/logger"
)

type SubPubService struct {
	subpub.UnimplementedPubSubServer
	subPub repository.SubPub
	Logger *logger.LogrusLogger
}

func NewSubPubService(logger *logger.LogrusLogger, subPub repository.SubPub) (*SubPubService, error) {
	if logger == nil || subPub == nil {
		return nil, fmt.Errorf("logger is nil or subPub is nil")
	}

	return &SubPubService{
		subPub: subPub,
		Logger: logger,
	}, nil
}
