package tracker

import (
	"context"
	"fmt"
)

// ListMacros returns all macros for a queue.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/api-ref/queues/get-macros.html
func (s *QueuesService) ListMacros(ctx context.Context, queueKey string) ([]*Macro, *Response, error) {
	u := fmt.Sprintf("v3/queues/%v/macros", queueKey)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var macros []*Macro
	resp, err := s.client.Do(ctx, req, &macros)
	if err != nil {
		return nil, resp, err
	}

	return macros, resp, nil
}

// GetMacro returns a single macro by its numeric ID within a queue.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/api-ref/queues/get-macro.html
func (s *QueuesService) GetMacro(ctx context.Context, queueKey string, macroID int) (*Macro, *Response, error) {
	u := fmt.Sprintf("v3/queues/%v/macros/%v", queueKey, macroID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	m := new(Macro)
	resp, err := s.client.Do(ctx, req, m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// CreateMacro creates a new macro in a queue.
// The API returns 201 Created on success.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/api-ref/queues/create-macro.html
func (s *QueuesService) CreateMacro(ctx context.Context, queueKey string, macro *MacroCreateRequest) (*Macro, *Response, error) {
	u := fmt.Sprintf("v3/queues/%v/macros", queueKey)

	req, err := s.client.NewRequest("POST", u, macro)
	if err != nil {
		return nil, nil, err
	}

	m := new(Macro)
	resp, err := s.client.Do(ctx, req, m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// EditMacro updates an existing macro in a queue.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/api-ref/queues/change-macro.html
func (s *QueuesService) EditMacro(ctx context.Context, queueKey string, macroID int, macro *MacroEditRequest) (*Macro, *Response, error) {
	u := fmt.Sprintf("v3/queues/%v/macros/%v", queueKey, macroID)

	req, err := s.client.NewRequest("PATCH", u, macro)
	if err != nil {
		return nil, nil, err
	}

	m := new(Macro)
	resp, err := s.client.Do(ctx, req, m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// DeleteMacro deletes a macro from a queue.
// The API returns 204 No Content on success.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/api-ref/queues/delete-macro.html
func (s *QueuesService) DeleteMacro(ctx context.Context, queueKey string, macroID int) (*Response, error) {
	u := fmt.Sprintf("v3/queues/%v/macros/%v", queueKey, macroID)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
