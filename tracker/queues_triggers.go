package tracker

import (
	"context"
	"fmt"
)

// ListTriggers returns all triggers for a queue.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/api-ref/queues/get-triggers.html
func (s *QueuesService) ListTriggers(ctx context.Context, queueKey string) ([]*Trigger, *Response, error) {
	u := fmt.Sprintf("v3/queues/%v/triggers", queueKey)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var triggers []*Trigger
	resp, err := s.client.Do(ctx, req, &triggers)
	if err != nil {
		return nil, resp, err
	}

	return triggers, resp, nil
}

// GetTrigger returns a single trigger by its ID within a queue.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/api-ref/queues/get-trigger.html
func (s *QueuesService) GetTrigger(ctx context.Context, queueKey string, triggerID int) (*Trigger, *Response, error) {
	u := fmt.Sprintf("v3/queues/%v/triggers/%v", queueKey, triggerID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	t := new(Trigger)
	resp, err := s.client.Do(ctx, req, t)
	if err != nil {
		return nil, resp, err
	}

	return t, resp, nil
}

// CreateTrigger creates a new trigger in a queue.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/api-ref/queues/create-trigger.html
func (s *QueuesService) CreateTrigger(ctx context.Context, queueKey string, trigger *TriggerCreateRequest) (*Trigger, *Response, error) {
	u := fmt.Sprintf("v3/queues/%v/triggers", queueKey)

	req, err := s.client.NewRequest("POST", u, trigger)
	if err != nil {
		return nil, nil, err
	}

	t := new(Trigger)
	resp, err := s.client.Do(ctx, req, t)
	if err != nil {
		return nil, resp, err
	}

	return t, resp, nil
}

// UpdateTrigger updates an existing trigger in a queue.
// The opts parameter must include Version for optimistic locking.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/api-ref/queues/change-trigger.html
func (s *QueuesService) UpdateTrigger(ctx context.Context, queueKey string, triggerID int, opts *TriggerUpdateOptions, trigger *TriggerUpdateRequest) (*Trigger, *Response, error) {
	u := fmt.Sprintf("v3/queues/%v/triggers/%v", queueKey, triggerID)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("PATCH", u, trigger)
	if err != nil {
		return nil, nil, err
	}

	t := new(Trigger)
	resp, err := s.client.Do(ctx, req, t)
	if err != nil {
		return nil, resp, err
	}

	return t, resp, nil
}
