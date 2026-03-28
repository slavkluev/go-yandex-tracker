package tracker

import (
	"context"
	"fmt"
)

// List returns all components.
// The API does not support pagination for this endpoint; all components are returned in a single response.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/concepts/components/get-components.html
func (s *ComponentsService) List(ctx context.Context) ([]*Component, *Response, error) {
	req, err := s.client.NewRequest("GET", "v2/components", nil)
	if err != nil {
		return nil, nil, err
	}

	var components []*Component
	resp, err := s.client.Do(ctx, req, &components)
	if err != nil {
		return nil, resp, err
	}

	return components, resp, nil
}

// Create creates a new component.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/concepts/components/create-component.html
func (s *ComponentsService) Create(ctx context.Context, component *ComponentRequest) (*Component, *Response, error) {
	req, err := s.client.NewRequest("POST", "v2/components", component)
	if err != nil {
		return nil, nil, err
	}

	c := new(Component)
	resp, err := s.client.Do(ctx, req, c)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, nil
}

// Get fetches a component by its ID.
// The componentID parameter is numeric (int) because Component.ID is *int.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/concepts/components/get-component.html
func (s *ComponentsService) Get(ctx context.Context, componentID int) (*Component, *Response, error) {
	u := fmt.Sprintf("v2/components/%d", componentID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	c := new(Component)
	resp, err := s.client.Do(ctx, req, c)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, nil
}

// Edit updates an existing component.
// The componentID parameter is numeric (int) because Component.ID is *int.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/concepts/components/edit-component.html
func (s *ComponentsService) Edit(ctx context.Context, componentID int, component *ComponentRequest) (*Component, *Response, error) {
	u := fmt.Sprintf("v2/components/%d", componentID)

	req, err := s.client.NewRequest("PATCH", u, component)
	if err != nil {
		return nil, nil, err
	}

	c := new(Component)
	resp, err := s.client.Do(ctx, req, c)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, nil
}

// Delete deletes a component by its ID.
// The componentID parameter is numeric (int) because Component.ID is *int.
// Returns (*Response, error) since 204 responses have no body.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/concepts/components/delete-component.html
func (s *ComponentsService) Delete(ctx context.Context, componentID int) (*Response, error) {
	u := fmt.Sprintf("v2/components/%d", componentID)

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
