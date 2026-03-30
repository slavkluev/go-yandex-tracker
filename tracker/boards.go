package tracker

import (
	"context"
	"fmt"
)

// Create creates a new board.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/boards/post-board
func (s *BoardsService) Create(ctx context.Context, board *BoardCreateRequest) (*Board, *Response, error) {
	req, err := s.client.NewRequest("POST", "v3/boards/", board)
	if err != nil {
		return nil, nil, err
	}

	b := new(Board)
	resp, err := s.client.Do(ctx, req, b)
	if err != nil {
		return nil, resp, err
	}

	return b, resp, nil
}

// Get fetches a board by its ID.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/boards/get-board
func (s *BoardsService) Get(ctx context.Context, boardID string) (*Board, *Response, error) {
	u := fmt.Sprintf("v3/boards/%s", boardID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	b := new(Board)
	resp, err := s.client.Do(ctx, req, b)
	if err != nil {
		return nil, resp, err
	}

	return b, resp, nil
}

// List returns all boards.
// The API does not support pagination for this endpoint; all boards are returned in a single response.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/boards/get-boards
func (s *BoardsService) List(ctx context.Context) ([]*Board, *Response, error) {
	req, err := s.client.NewRequest("GET", "v3/boards", nil)
	if err != nil {
		return nil, nil, err
	}

	var boards []*Board
	resp, err := s.client.Do(ctx, req, &boards)
	if err != nil {
		return nil, resp, err
	}

	return boards, resp, nil
}

// Edit updates an existing board.
// The version parameter is required for optimistic locking and is sent
// as the If-Match header value.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/boards/patch-board
func (s *BoardsService) Edit(ctx context.Context, boardID string, version string, board *BoardEditRequest) (*Board, *Response, error) {
	u := fmt.Sprintf("v3/boards/%s", boardID)

	req, err := s.client.NewRequest("PATCH", u, board)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("If-Match", version)

	b := new(Board)
	resp, err := s.client.Do(ctx, req, b)
	if err != nil {
		return nil, resp, err
	}

	return b, resp, nil
}

// Delete deletes a board by its ID.
// Returns (*Response, error) since 204 responses have no body.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/boards/delete-board
func (s *BoardsService) Delete(ctx context.Context, boardID string) (*Response, error) {
	u := fmt.Sprintf("v3/boards/%s", boardID)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
