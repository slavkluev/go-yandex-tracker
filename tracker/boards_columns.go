package tracker

import (
	"context"
	"fmt"
)

// ListColumns returns all columns for a board.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/boards/get-columns
func (s *BoardsService) ListColumns(ctx context.Context, boardID string) ([]*Column, *Response, error) {
	u := fmt.Sprintf("v3/boards/%s/columns", boardID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var columns []*Column
	resp, err := s.client.Do(ctx, req, &columns)
	if err != nil {
		return nil, resp, err
	}

	return columns, resp, nil
}

// CreateColumn creates a new column on a board.
// The version parameter is the current board version for optimistic locking,
// sent as the If-Match header value.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/boards/post-column
func (s *BoardsService) CreateColumn(ctx context.Context, boardID string, version string, column *ColumnCreateRequest) (*Column, *Response, error) {
	u := fmt.Sprintf("v3/boards/%s/columns/", boardID)

	req, err := s.client.NewRequest("POST", u, column)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("If-Match", version)

	c := new(Column)
	resp, err := s.client.Do(ctx, req, c)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, nil
}

// EditColumn updates an existing column on a board.
// The version parameter is the current board version for optimistic locking,
// sent as the If-Match header value.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/boards/patch-column
func (s *BoardsService) EditColumn(ctx context.Context, boardID string, columnID string, version string, column *ColumnEditRequest) (*Column, *Response, error) {
	u := fmt.Sprintf("v3/boards/%s/columns/%s", boardID, columnID)

	req, err := s.client.NewRequest("PATCH", u, column)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("If-Match", version)

	c := new(Column)
	resp, err := s.client.Do(ctx, req, c)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, nil
}

// DeleteColumn deletes a column from a board.
// The version parameter is the current board version for optimistic locking,
// sent as the If-Match header value.
// Returns (*Response, error) since 204 responses have no body.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/boards/delete-column
func (s *BoardsService) DeleteColumn(ctx context.Context, boardID string, columnID string, version string) (*Response, error) {
	u := fmt.Sprintf("v3/boards/%s/columns/%s", boardID, columnID)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("If-Match", version)

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
