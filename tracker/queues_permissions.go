package tracker

import (
	"context"
	"fmt"
)

// GetPermissions returns the access permissions for a queue.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/concepts/queues/get-permissions.html
func (s *QueuesService) GetPermissions(ctx context.Context, queueKey string) (*QueuePermissions, *Response, error) {
	u := fmt.Sprintf("v2/queues/%v/permissions", queueKey)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	perms := new(QueuePermissions)
	resp, err := s.client.Do(ctx, req, perms)
	if err != nil {
		return nil, resp, err
	}

	return perms, resp, nil
}

// UpdatePermissions updates the access permissions for a queue.
// The request body specifies which users, groups, and roles to add or remove
// for each permission type (create, write, read, grant).
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/concepts/queues/change-permissions.html
func (s *QueuesService) UpdatePermissions(ctx context.Context, queueKey string, perms *QueuePermissionsUpdateRequest) (*QueuePermissions, *Response, error) {
	u := fmt.Sprintf("v2/queues/%v/permissions", queueKey)

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

// CheckPermission checks whether the current user has a specific permission
// on the queue. The permissionCode parameter is one of: "create", "write",
// "read", "grant". Returns the permissions structure if the user has the
// permission; returns an error (403) if they do not.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/concepts/queues/check-permissions.html
func (s *QueuesService) CheckPermission(ctx context.Context, queueKey, permissionCode string) (*QueuePermissions, *Response, error) {
	u := fmt.Sprintf("v2/queues/%v/checkPermissions/%v", queueKey, permissionCode)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	perms := new(QueuePermissions)
	resp, err := s.client.Do(ctx, req, perms)
	if err != nil {
		return nil, resp, err
	}

	return perms, resp, nil
}
