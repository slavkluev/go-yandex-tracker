package tracker

import (
	"context"
	"fmt"
)

// CreateChecklistItem adds a new checklist item to an entity.
// The API returns the full updated Entity with the new checklist item included.
//
// Note: Checklists are supported for projects and portfolios only, not goals.
// No client-side validation is performed on entityType; the API returns
// appropriate errors for unsupported entity types.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/entities/checklists/add-checklist
func (s *EntitiesService) CreateChecklistItem(ctx context.Context, entityType EntityType, id string, item *ChecklistItemRequest) (*Entity, *Response, error) {
	u := fmt.Sprintf("v3/entities/%v/%v/checklistItems", entityType, id)

	req, err := s.client.NewRequest("POST", u, item)
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

// EditChecklistItem updates an existing checklist item on an entity.
// The API returns the full updated Entity reflecting the checklist change.
//
// Note: Checklists are supported for projects and portfolios only, not goals.
// No client-side validation is performed on entityType; the API returns
// appropriate errors for unsupported entity types.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/entities/checklists/patch-checklist-item
func (s *EntitiesService) EditChecklistItem(ctx context.Context, entityType EntityType, id, checklistItemID string, item *ChecklistItemRequest) (*Entity, *Response, error) {
	u := fmt.Sprintf("v3/entities/%v/%v/checklistItems/%v", entityType, id, checklistItemID)

	req, err := s.client.NewRequest("PATCH", u, item)
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

// MoveChecklistItem repositions a checklist item on an entity.
// Use ChecklistMoveRequest to specify the target position (before which item).
// The API returns the full updated Entity reflecting the new order.
//
// This operation is specific to entity checklists and is not available
// on issue checklists.
//
// Note: Checklists are supported for projects and portfolios only, not goals.
// No client-side validation is performed on entityType; the API returns
// appropriate errors for unsupported entity types.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/entities/checklists/move-checklist-item
func (s *EntitiesService) MoveChecklistItem(ctx context.Context, entityType EntityType, id, checklistItemID string, move *ChecklistMoveRequest) (*Entity, *Response, error) {
	u := fmt.Sprintf("v3/entities/%v/%v/checklistItems/%v/_move", entityType, id, checklistItemID)

	req, err := s.client.NewRequest("POST", u, move)
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

// DeleteChecklistItem removes a single checklist item from an entity.
// Unlike most delete operations, the API returns HTTP 200 with the full
// updated Entity in the response body (not HTTP 204 No Content).
//
// Note: Checklists are supported for projects and portfolios only, not goals.
// No client-side validation is performed on entityType; the API returns
// appropriate errors for unsupported entity types.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/entities/checklists/delete-checklist-item
func (s *EntitiesService) DeleteChecklistItem(ctx context.Context, entityType EntityType, id, checklistItemID string) (*Entity, *Response, error) {
	u := fmt.Sprintf("v3/entities/%v/%v/checklistItems/%v", entityType, id, checklistItemID)

	req, err := s.client.NewRequest("DELETE", u, nil)
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

// DeleteChecklist removes the entire checklist from an entity.
// Unlike most delete operations, the API returns HTTP 200 with the full
// updated Entity in the response body (not HTTP 204 No Content).
//
// This operation is specific to entity checklists and is not available
// on issue checklists.
//
// Note: Checklists are supported for projects and portfolios only, not goals.
// No client-side validation is performed on entityType; the API returns
// appropriate errors for unsupported entity types.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/entities/checklists/delete-checklist
func (s *EntitiesService) DeleteChecklist(ctx context.Context, entityType EntityType, id string) (*Entity, *Response, error) {
	u := fmt.Sprintf("v3/entities/%v/%v/checklistItems", entityType, id)

	req, err := s.client.NewRequest("DELETE", u, nil)
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
