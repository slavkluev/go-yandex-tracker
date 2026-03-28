package tracker

import (
	"context"
	"fmt"
)

// List returns all sprints for a board.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/get-sprints
func (s *SprintsService) List(ctx context.Context, boardID int) ([]*Sprint, *Response, error) {
	u := fmt.Sprintf("v3/boards/%d/sprints", boardID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var sprints []*Sprint
	resp, err := s.client.Do(ctx, req, &sprints)
	if err != nil {
		return nil, resp, err
	}

	return sprints, resp, nil
}

// Create creates a new sprint.
// The board ID is specified in the request body (SprintCreateRequest.Board.ID),
// not in the URL. The API endpoint is POST /v3/sprints.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/post-sprint
func (s *SprintsService) Create(ctx context.Context, sprint *SprintCreateRequest) (*Sprint, *Response, error) {
	req, err := s.client.NewRequest("POST", "v3/sprints", sprint)
	if err != nil {
		return nil, nil, err
	}

	sp := new(Sprint)
	resp, err := s.client.Do(ctx, req, sp)
	if err != nil {
		return nil, resp, err
	}

	return sp, resp, nil
}

// Get fetches a sprint by its ID.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/get-sprint
func (s *SprintsService) Get(ctx context.Context, sprintID int) (*Sprint, *Response, error) {
	u := fmt.Sprintf("v3/sprints/%d", sprintID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	sp := new(Sprint)
	resp, err := s.client.Do(ctx, req, sp)
	if err != nil {
		return nil, resp, err
	}

	return sp, resp, nil
}

// Note: Sprint editing (PATCH /v3/sprints/<id>) is NOT supported by the
// Yandex Tracker API. The official API reference documents only three sprint
// endpoints: list, get, and create. There is no PATCH or PUT endpoint for
// sprints. See SPRT-04 in REQUIREMENTS.md.
