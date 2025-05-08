package usecase

import (
	subpub "github.com/samantonio28/vk-task-subscribe/internal/delivery"
)

func (s *SubPubService) Subscribe(*subpub.SubscribeRequest, subpub.PubSub_SubscribeServer) error {
	return nil
}
