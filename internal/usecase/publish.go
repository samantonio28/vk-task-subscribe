package usecase

import (
	"context"

	subpub "github.com/samantonio28/vk-task-subscribe/internal/delivery"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *SubPubService) Publish(ctx context.Context, req *subpub.PublishRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
