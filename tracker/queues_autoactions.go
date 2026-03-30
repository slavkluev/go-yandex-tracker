package tracker

import (
	"context"
	"fmt"
)

// ListAutoActions returns all auto-actions for a queue.
func (s *QueuesService) ListAutoActions(ctx context.Context, queueKey string) ([]*AutoAction, *Response, error) {
	u := fmt.Sprintf("v3/queues/%v/autoactions", queueKey)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var autoActions []*AutoAction
	resp, err := s.client.Do(ctx, req, &autoActions)
	if err != nil {
		return nil, resp, err
	}

	return autoActions, resp, nil
}

// GetAutoAction returns a single auto-action by its ID within a queue.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/queues/get-autoaction
func (s *QueuesService) GetAutoAction(ctx context.Context, queueKey string, autoActionID string) (*AutoAction, *Response, error) {
	u := fmt.Sprintf("v3/queues/%v/autoactions/%v", queueKey, autoActionID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	a := new(AutoAction)
	resp, err := s.client.Do(ctx, req, a)
	if err != nil {
		return nil, resp, err
	}

	return a, resp, nil
}

// CreateAutoAction creates a new auto-action in a queue.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/queues/create-autoaction
func (s *QueuesService) CreateAutoAction(ctx context.Context, queueKey string, autoAction *AutoActionCreateRequest) (*AutoAction, *Response, error) {
	u := fmt.Sprintf("v3/queues/%v/autoactions", queueKey)

	req, err := s.client.NewRequest("POST", u, autoAction)
	if err != nil {
		return nil, nil, err
	}

	a := new(AutoAction)
	resp, err := s.client.Do(ctx, req, a)
	if err != nil {
		return nil, resp, err
	}

	return a, resp, nil
}

// UpdateAutoAction updates an existing auto-action in a queue.
func (s *QueuesService) UpdateAutoAction(ctx context.Context, queueKey string, autoActionID string, autoAction *AutoActionUpdateRequest) (*AutoAction, *Response, error) {
	u := fmt.Sprintf("v3/queues/%v/autoactions/%v", queueKey, autoActionID)

	req, err := s.client.NewRequest("PATCH", u, autoAction)
	if err != nil {
		return nil, nil, err
	}

	a := new(AutoAction)
	resp, err := s.client.Do(ctx, req, a)
	if err != nil {
		return nil, resp, err
	}

	return a, resp, nil
}
