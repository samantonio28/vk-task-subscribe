package usecase

import (
	"context"
	"io"

	subpub "github.com/samantonio28/vk-task-subscribe/internal/delivery"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *SubPubService) Subscribe(req *subpub.SubscribeRequest, stream subpub.PubSub_SubscribeServer) error {
	key := req.GetKey()
	if key == "" {
		return status.Error(codes.InvalidArgument, "key cannot be empty")
	}

	ctx := stream.Context()

	s.mu.Lock()
	if s.subs == nil {
		s.subs = make(map[string]map[subpub.PubSub_SubscribeServer]struct{})
	}
	if _, exists := s.subs[key]; !exists {
		s.subs[key] = make(map[subpub.PubSub_SubscribeServer]struct{})
	}
	s.subs[key][stream] = struct{}{}
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.subs[key], stream)
		if len(s.subs[key]) == 0 {
			delete(s.subs, key)
		}
		s.mu.Unlock()
	}()

	msgChan := make(chan string, 100)

	sub, err := s.subPub.Subscribe(key, func(msg any) {
		if data, ok := msg.(string); ok {
			select {
			case msgChan <- data:
			default:
				s.Logger.WithFields(&logrus.Fields{
					"key":  key,
					"data": data,
				}).Warn("Message channel full, dropping message")
			}
		}
	})

	if err != nil {
		s.Logger.WithFields(&logrus.Fields{
			"key":   key,
			"error": err,
		}).Error("Failed to subscribe to key")
		return status.Errorf(codes.Internal, "failed to subscribe: %v", err)
	}
	defer sub.Unsubscribe()

	for {
		select {
		case <-ctx.Done():
			switch ctx.Err() {
			case context.Canceled:
				s.Logger.WithFields(&logrus.Fields{
					"key": key,
				}).Debug("Subscription canceled by client")
				return status.Error(codes.Canceled, "subscription canceled by client")
			case context.DeadlineExceeded:
				s.Logger.WithFields(&logrus.Fields{
					"key": key,
				}).Debug("Subscription deadline exceeded")
				return status.Error(codes.DeadlineExceeded, "subscription deadline exceeded")
			default:
				return status.FromContextError(ctx.Err()).Err()
			}

		case data, ok := <-msgChan:
			if !ok {
				s.Logger.WithFields(&logrus.Fields{
					"key": key,
				}).Debug("Message channel closed for key")
				return status.Error(codes.Aborted, "message channel closed")
			}

			if err := stream.Send(&subpub.Event{Data: data}); err != nil {
				if err == io.EOF {
					s.Logger.WithFields(&logrus.Fields{
						"key": key,
					}).Debug("Client disconnected from key")
					return nil
				}
				s.Logger.WithFields(&logrus.Fields{
					"key":   key,
					"error": err,
				}).Error("Failed to send message to subscriber")
				return status.Errorf(codes.Internal, "failed to send message: %v", err)
			}
		}
	}
}
