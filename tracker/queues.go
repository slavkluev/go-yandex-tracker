package tracker

import (
	"context"
	"fmt"
)

// Create creates a new queue.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/queues/create-queue
func (s *QueuesService) Create(ctx context.Context, queue *QueueCreateRequest) (*Queue, *Response, error) {
	req, err := s.client.NewRequest("POST", "v2/queues/", queue)
	if err != nil {
		return nil, nil, err
	}

	q := new(Queue)
	resp, err := s.client.Do(ctx, req, q)
	if err != nil {
		return nil, resp, err
	}

	return q, resp, nil
}

// Get fetches a queue by its key or ID.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/queues/get-queue
func (s *QueuesService) Get(ctx context.Context, queueKey string, opts *QueueGetOptions) (*Queue, *Response, error) {
	u := fmt.Sprintf("v2/queues/%v", queueKey)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	queue := new(Queue)
	resp, err := s.client.Do(ctx, req, queue)
	if err != nil {
		return nil, resp, err
	}

	return queue, resp, nil
}

// List returns all queues available to the user.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/queues/get-queues
func (s *QueuesService) List(ctx context.Context, opts *QueueListOptions) ([]*Queue, *Response, error) {
	u := "v2/queues"
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var queues []*Queue
	resp, err := s.client.Do(ctx, req, &queues)
	if err != nil {
		return nil, resp, err
	}

	return queues, resp, nil
}

// Delete deletes a queue.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/queues/delete-queue
func (s *QueuesService) Delete(ctx context.Context, queueKey string) (*Response, error) {
	u := fmt.Sprintf("v2/queues/%v", queueKey)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Restore restores a previously deleted queue.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/queues/restore-queue
func (s *QueuesService) Restore(ctx context.Context, queueKey string) (*Queue, *Response, error) {
	u := fmt.Sprintf("v2/queues/%v/_restore", queueKey)

	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, nil, err
	}

	queue := new(Queue)
	resp, err := s.client.Do(ctx, req, queue)
	if err != nil {
		return nil, resp, err
	}

	return queue, resp, nil
}
