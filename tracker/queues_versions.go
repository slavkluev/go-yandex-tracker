package tracker

import (
	"context"
	"fmt"
)

// ListVersions returns all versions defined in a queue.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/concepts/queues/get-versions.html
func (s *QueuesService) ListVersions(ctx context.Context, queueKey string) ([]*QueueVersion, *Response, error) {
	u := fmt.Sprintf("v2/queues/%v/versions", queueKey)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var versions []*QueueVersion
	resp, err := s.client.Do(ctx, req, &versions)
	if err != nil {
		return nil, resp, err
	}

	return versions, resp, nil
}

// ListTags returns all tags used in a queue.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/concepts/queues/get-tags.html
func (s *QueuesService) ListTags(ctx context.Context, queueKey string) ([]string, *Response, error) {
	u := fmt.Sprintf("v2/queues/%v/tags", queueKey)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var tags []string
	resp, err := s.client.Do(ctx, req, &tags)
	if err != nil {
		return nil, resp, err
	}

	return tags, resp, nil
}
