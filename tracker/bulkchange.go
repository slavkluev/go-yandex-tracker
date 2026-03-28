package tracker

import (
	"context"
	"fmt"
)

// Move initiates an asynchronous bulk move of issues to a different queue.
// The returned BulkChange contains the operation ID and initial status.
// Use GetStatus to poll for completion.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/concepts/bulkchange/bulk-move-issues
func (s *BulkChangeService) Move(ctx context.Context, move *BulkMoveRequest) (*BulkChange, *Response, error) {
	req, err := s.client.NewRequest("POST", "v3/bulkchange/_move", move)
	if err != nil {
		return nil, nil, err
	}

	bc := new(BulkChange)
	resp, err := s.client.Do(ctx, req, bc)
	if err != nil {
		return nil, resp, err
	}

	return bc, resp, nil
}

// Update initiates an asynchronous bulk update of issue field values.
// The returned BulkChange contains the operation ID and initial status.
// Use GetStatus to poll for completion.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/concepts/bulkchange/bulk-update-issues
func (s *BulkChangeService) Update(ctx context.Context, update *BulkUpdateRequest) (*BulkChange, *Response, error) {
	req, err := s.client.NewRequest("POST", "v3/bulkchange/_update", update)
	if err != nil {
		return nil, nil, err
	}

	bc := new(BulkChange)
	resp, err := s.client.Do(ctx, req, bc)
	if err != nil {
		return nil, resp, err
	}

	return bc, resp, nil
}

// Transition initiates an asynchronous bulk status transition of issues.
// The returned BulkChange contains the operation ID and initial status.
// Use GetStatus to poll for completion.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/concepts/bulkchange/bulk-transition
func (s *BulkChangeService) Transition(ctx context.Context, transition *BulkTransitionRequest) (*BulkChange, *Response, error) {
	req, err := s.client.NewRequest("POST", "v3/bulkchange/_transition", transition)
	if err != nil {
		return nil, nil, err
	}

	bc := new(BulkChange)
	resp, err := s.client.Do(ctx, req, bc)
	if err != nil {
		return nil, resp, err
	}

	return bc, resp, nil
}

// GetStatus retrieves the current status of a bulk change operation.
// The bulkChangeID is a hex string (e.g., "593cd211ef7e8a33********").
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/concepts/bulkchange/bulk-move-info
func (s *BulkChangeService) GetStatus(ctx context.Context, bulkChangeID string) (*BulkChange, *Response, error) {
	u := fmt.Sprintf("v3/bulkchange/%v", bulkChangeID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	bc := new(BulkChange)
	resp, err := s.client.Do(ctx, req, bc)
	if err != nil {
		return nil, resp, err
	}

	return bc, resp, nil
}
