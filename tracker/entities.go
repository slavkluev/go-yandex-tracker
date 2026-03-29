package tracker

import (
	"context"
	"fmt"
)

// Create creates a new entity of the specified type.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/entities/create-entity
func (s *EntitiesService) Create(ctx context.Context, entityType EntityType, req *EntityCreateRequest) (*Entity, *Response, error) {
	u := fmt.Sprintf("v3/entities/%v", entityType)

	httpReq, err := s.client.NewRequest("POST", u, req)
	if err != nil {
		return nil, nil, err
	}

	entity := new(Entity)
	resp, err := s.client.Do(ctx, httpReq, entity)
	if err != nil {
		return nil, resp, err
	}

	return entity, resp, nil
}

// Get retrieves an entity by its ID.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/entities/get-entity
func (s *EntitiesService) Get(ctx context.Context, entityType EntityType, id string, opts *EntityGetOptions) (*Entity, *Response, error) {
	u := fmt.Sprintf("v3/entities/%v/%v", entityType, id)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	entity := new(Entity)
	resp, err := s.client.Do(ctx, req, entity)
	if err != nil {
		return nil, resp, err
	}

	return entity, resp, nil
}

// Update modifies an existing entity.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/entities/update-entity
func (s *EntitiesService) Update(ctx context.Context, entityType EntityType, id string, req *EntityUpdateRequest) (*Entity, *Response, error) {
	u := fmt.Sprintf("v3/entities/%v/%v", entityType, id)

	httpReq, err := s.client.NewRequest("PATCH", u, req)
	if err != nil {
		return nil, nil, err
	}

	entity := new(Entity)
	resp, err := s.client.Do(ctx, httpReq, entity)
	if err != nil {
		return nil, resp, err
	}

	return entity, resp, nil
}

// Delete removes an entity.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/entities/delete-entity
func (s *EntitiesService) Delete(ctx context.Context, entityType EntityType, id string, opts *EntityDeleteOptions) (*Response, error) {
	u := fmt.Sprintf("v3/entities/%v/%v", entityType, id)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Search searches for entities of the specified type.
// The search uses POST with a request body for filters and query parameters for pagination.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/entities/search-entities
func (s *EntitiesService) Search(ctx context.Context, entityType EntityType, req *EntitySearchRequest, opts *EntitySearchOptions) (*EntitySearchResponse, *Response, error) {
	u := fmt.Sprintf("v3/entities/%v/_search", entityType)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	httpReq, err := s.client.NewRequest("POST", u, req)
	if err != nil {
		return nil, nil, err
	}

	result := new(EntitySearchResponse)
	resp, err := s.client.Do(ctx, httpReq, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// BulkChange performs a bulk change operation on entities of the specified type.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/entities/bulkchange-entities
func (s *EntitiesService) BulkChange(ctx context.Context, entityType EntityType, req *EntityBulkChangeRequest) (*BulkChange, *Response, error) {
	u := fmt.Sprintf("v3/entities/%v/bulkchange/_update", entityType)

	httpReq, err := s.client.NewRequest("POST", u, req)
	if err != nil {
		return nil, nil, err
	}

	bc := new(BulkChange)
	resp, err := s.client.Do(ctx, httpReq, bc)
	if err != nil {
		return nil, resp, err
	}

	return bc, resp, nil
}

// GetEvents retrieves the event history for an entity.
// Pagination is cursor-based using the From and Direction options, not page-based.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/entities/get-events-relative
func (s *EntitiesService) GetEvents(ctx context.Context, entityType EntityType, id string, opts *EntityEventsOptions) (*EntityEventsResponse, *Response, error) {
	u := fmt.Sprintf("v3/entities/%v/%v/events/_relative", entityType, id)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(EntityEventsResponse)
	resp, err := s.client.Do(ctx, req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
