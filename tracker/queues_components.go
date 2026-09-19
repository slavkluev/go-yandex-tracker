package tracker

import (
	"context"
	"fmt"
)

// ListComponents returns the components of a queue. By default each
// component is returned in full; opts.Fields limits the response to the
// listed fields, and Self and ID are always included.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api/queues/get-queue-components
func (s *QueuesService) ListComponents(ctx context.Context, queueKey string, opts *QueueComponentsListOptions) ([]*Component, *Response, error) {
	u := fmt.Sprintf("v3/queues/%v/components", queueKey)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var components []*Component
	resp, err := s.client.Do(ctx, req, &components)
	if err != nil {
		return nil, resp, err
	}

	return components, resp, nil
}
