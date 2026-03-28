package tracker

import (
	"context"
	"fmt"
)

// ListChecklistItems returns the checklist items for an issue.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/concepts/issues/get-checklist.html
func (s *IssuesService) ListChecklistItems(ctx context.Context, issueKey string) ([]*ChecklistItem, *Response, error) {
	u := fmt.Sprintf("v2/issues/%v/checklistItems", issueKey)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var items []*ChecklistItem
	resp, err := s.client.Do(ctx, req, &items)
	if err != nil {
		return nil, resp, err
	}

	return items, resp, nil
}

// CreateChecklistItem creates a new checklist item on an issue.
// The API returns the full updated Issue, not the created checklist item.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/concepts/issues/add-checklist-item.html
func (s *IssuesService) CreateChecklistItem(ctx context.Context, issueKey string, item *ChecklistItemRequest) (*Issue, *Response, error) {
	u := fmt.Sprintf("v2/issues/%v/checklistItems", issueKey)

	req, err := s.client.NewRequest("POST", u, item)
	if err != nil {
		return nil, nil, err
	}

	issue := new(Issue)
	resp, err := s.client.Do(ctx, req, issue)
	if err != nil {
		return nil, resp, err
	}

	return issue, resp, nil
}

// EditChecklistItem updates a checklist item on an issue.
// The API returns the full updated Issue, not the edited checklist item.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/concepts/issues/edit-checklist.html
func (s *IssuesService) EditChecklistItem(ctx context.Context, issueKey, checklistItemID string, item *ChecklistItemRequest) (*Issue, *Response, error) {
	u := fmt.Sprintf("v2/issues/%v/checklistItems/%v", issueKey, checklistItemID)

	req, err := s.client.NewRequest("PATCH", u, item)
	if err != nil {
		return nil, nil, err
	}

	issue := new(Issue)
	resp, err := s.client.Do(ctx, req, issue)
	if err != nil {
		return nil, resp, err
	}

	return issue, resp, nil
}

// DeleteChecklistItem deletes a checklist item from an issue.
// Unlike most delete methods, the API returns HTTP 200 with the full
// updated Issue, not HTTP 204 No Content.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/concepts/issues/delete-checklist-item.html
func (s *IssuesService) DeleteChecklistItem(ctx context.Context, issueKey, checklistItemID string) (*Issue, *Response, error) {
	u := fmt.Sprintf("v2/issues/%v/checklistItems/%v", issueKey, checklistItemID)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, nil, err
	}

	issue := new(Issue)
	resp, err := s.client.Do(ctx, req, issue)
	if err != nil {
		return nil, resp, err
	}

	return issue, resp, nil
}
