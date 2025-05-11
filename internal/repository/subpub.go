package subpub

import (
	"context"
	"fmt"
	"slices"
	"sync"
)

// MessageHandler is a callback function that processes messages delivered to subscribers.
type MessageHandler func(msg interface{})

type Subscription interface {
	// Unsubscribe will remove interest in the current subject subscription is for.
	Unsubscribe()
}

type SubPub interface {
	// Subscribe creates an asynchronous queue subscriber on the given subject.
	Subscribe(subject string, cb MessageHandler) (Subscription, error)

	// Publish publishes the msg argument to the given subject.
	Publish(subject string, msg interface{}) error

	// Close will shutdown sub-pub system.
	// May be blocked by data delivery until the context is canceled.
	Close(ctx context.Context) error
}

type subPub struct {
	mu          sync.RWMutex
	subjects    map[string][]*subscription
	closed      bool
	publishChan chan publishRequest
	wg          sync.WaitGroup
}

type subscription struct {
	subject string
	handler MessageHandler
	removed bool
	sp      *subPub
}

type publishRequest struct {
	subject string
	msg     any
}

func NewSubPub() SubPub {
	sp := &subPub{
		subjects:    make(map[string][]*subscription),
		publishChan: make(chan publishRequest, 100),
	}

	sp.wg.Add(1)
	go sp.publisher()

	return sp
}

func (sp *subPub) Subscribe(subject string, cb MessageHandler) (Subscription, error) {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	if sp.closed {
		return nil, fmt.Errorf("subpub is closed")
	}

	sub := &subscription{
		subject: subject,
		handler: cb,
		sp:      sp,
	}

	sp.subjects[subject] = append(sp.subjects[subject], sub)
	return sub, nil
}

func (sp *subPub) Publish(subject string, msg any) error {
	sp.mu.RLock()
	defer sp.mu.RUnlock()

	if sp.closed {
		return fmt.Errorf("subpub is closed")
	}

	sp.publishChan <- publishRequest{subject, msg}
	return nil
}

func (sp *subPub) Close(ctx context.Context) error {
	sp.mu.Lock()
	if sp.closed {
		sp.mu.Unlock()
		return nil
	}
	sp.closed = true
	close(sp.publishChan)
	sp.mu.Unlock()

	done := make(chan struct{})
	go func() {
		sp.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (sp *subPub) publisher() {
	defer sp.wg.Done()

	for req := range sp.publishChan {
		sp.mu.RLock()
		subs, exists := sp.subjects[req.subject]
		if !exists {
			sp.mu.RUnlock()
			continue
		}

		// копия нужна, чтобы обработать всех подписчиков,
		// если они отпишутся во время обработки
		subsCopy := make([]*subscription, len(subs))
		copy(subsCopy, subs)
		sp.mu.RUnlock()

		var wg sync.WaitGroup
		for _, sub := range subsCopy {
			sub := sub
			sp.mu.RLock()
			removed := sub.removed
			sp.mu.RUnlock()

			if removed {
				continue
			}

			wg.Add(1)
			go func() {
				defer wg.Done()
				sub.handler(req.msg)
			}()
		}
		wg.Wait()
	}
}

func (s *subscription) Unsubscribe() {
	s.sp.mu.Lock()
	defer s.sp.mu.Unlock()

	subs, exists := s.sp.subjects[s.subject]
	if !exists {
		return
	}
	for i, sub := range subs {
		if sub == s {
			sub.removed = true
			s.sp.subjects[s.subject] = slices.Delete(subs, i, i+1)
			break
		}
	}

	if len(s.sp.subjects[s.subject]) == 0 {
		delete(s.sp.subjects, s.subject)
	}
}
