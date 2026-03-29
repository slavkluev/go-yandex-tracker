package tracker

import (
	"context"
	"fmt"
)

// List returns all saved filters for the current user.
func (s *FiltersService) List(ctx context.Context) ([]*Filter, *Response, error) {
	req, err := s.client.NewRequest("GET", "v3/filters", nil)
	if err != nil {
		return nil, nil, err
	}

	var filters []*Filter
	resp, err := s.client.Do(ctx, req, &filters)
	if err != nil {
		return nil, resp, err
	}

	return filters, resp, nil
}

// Get fetches a saved filter by its numeric ID.
func (s *FiltersService) Get(ctx context.Context, filterID int) (*Filter, *Response, error) {
	u := fmt.Sprintf("v3/filters/%v", filterID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	f := new(Filter)
	resp, err := s.client.Do(ctx, req, f)
	if err != nil {
		return nil, resp, err
	}

	return f, resp, nil
}

// Create creates a new saved filter.
func (s *FiltersService) Create(ctx context.Context, filter *FilterRequest) (*Filter, *Response, error) {
	req, err := s.client.NewRequest("POST", "v3/filters", filter)
	if err != nil {
		return nil, nil, err
	}

	f := new(Filter)
	resp, err := s.client.Do(ctx, req, f)
	if err != nil {
		return nil, resp, err
	}

	return f, resp, nil
}

// Update updates an existing saved filter.
func (s *FiltersService) Update(ctx context.Context, filterID int, filter *FilterRequest) (*Filter, *Response, error) {
	u := fmt.Sprintf("v3/filters/%v", filterID)

	req, err := s.client.NewRequest("PATCH", u, filter)
	if err != nil {
		return nil, nil, err
	}

	f := new(Filter)
	resp, err := s.client.Do(ctx, req, f)
	if err != nil {
		return nil, resp, err
	}

	return f, resp, nil
}

// Delete deletes a saved filter.
// The API returns 204 No Content on success.
func (s *FiltersService) Delete(ctx context.Context, filterID int) (*Response, error) {
	u := fmt.Sprintf("v3/filters/%v", filterID)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
