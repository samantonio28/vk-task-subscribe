package usecase

import (
	"context"

	subpub "github.com/samantonio28/vk-task-subscribe/internal/delivery"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *SubPubService) Publish(ctx context.Context, req *subpub.PublishRequest) (*emptypb.Empty, error) {
	key, data := req.GetKey(), req.GetData()

	if key == "" {
		return nil, status.Error(codes.InvalidArgument, "key cannot be empty")
	}

	if err := s.subPub.Publish(key, data); err != nil {
		s.Logger.WithFields(&logrus.Fields{
			"key":   key,
			"error": err,
		}).Error("Failed to publish message")
		return nil, status.Errorf(codes.Internal, "failed to publish message: %v", err)
	}

	s.Logger.WithFields(&logrus.Fields{
		"key":  key,
		"data": data,
	}).Info("Published message")
	return &emptypb.Empty{}, nil
}
