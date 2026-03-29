package tracker

import (
	"context"
	"fmt"
)

// UpdatePermissions updates the access permissions for a queue.
// The request body specifies which users, groups, and roles to add or remove
// for each permission type (create, write, read, grant).
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/queues/manage-access
func (s *QueuesService) UpdatePermissions(ctx context.Context, queueKey string, perms *QueuePermissionsUpdateRequest) (*QueuePermissions, *Response, error) {
	u := fmt.Sprintf("v3/queues/%v/permissions", queueKey)

	req, err := s.client.NewRequest("PATCH", u, perms)
	if err != nil {
		return nil, nil, err
	}

	result := new(QueuePermissions)
	resp, err := s.client.Do(ctx, req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
