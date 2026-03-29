package tracker

import (
	"context"
	"fmt"
)

// GetAccess retrieves the access settings for an entity.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/entities/get-access
func (s *EntitiesService) GetAccess(ctx context.Context, entityType EntityType, id string) (*EntityAccess, *Response, error) {
	u := fmt.Sprintf("v3/entities/%v/%v/extendedPermissions", entityType, id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	access := new(EntityAccess)
	resp, err := s.client.Do(ctx, req, access)
	if err != nil {
		return nil, resp, err
	}

	return access, resp, nil
}

// UpdateAccess modifies the access settings for an entity.
// The request uses grant/revoke semantics with UPPERCASE permission keys (READ, WRITE, GRANT).
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/entities/patch-access
func (s *EntitiesService) UpdateAccess(ctx context.Context, entityType EntityType, id string, body *EntityAccessUpdateRequest) (*EntityAccess, *Response, error) {
	u := fmt.Sprintf("v3/entities/%v/%v/extendedPermissions", entityType, id)

	req, err := s.client.NewRequest("PATCH", u, body)
	if err != nil {
		return nil, nil, err
	}

	access := new(EntityAccess)
	resp, err := s.client.Do(ctx, req, access)
	if err != nil {
		return nil, resp, err
	}

	return access, resp, nil
}
