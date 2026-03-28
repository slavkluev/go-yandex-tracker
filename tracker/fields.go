package tracker

import (
	"context"
	"fmt"
)

// List returns all global issue fields.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/concepts/issues/get-global-fields.html
func (s *FieldsService) List(ctx context.Context) ([]*Field, *Response, error) {
	req, err := s.client.NewRequest("GET", "v3/fields", nil)
	if err != nil {
		return nil, nil, err
	}

	var fields []*Field
	resp, err := s.client.Do(ctx, req, &fields)
	if err != nil {
		return nil, resp, err
	}

	return fields, resp, nil
}

// Get fetches a single global issue field by its ID.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/concepts/issues/get-global-field.html
func (s *FieldsService) Get(ctx context.Context, fieldID string) (*Field, *Response, error) {
	u := fmt.Sprintf("v3/fields/%v", fieldID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	field := new(Field)
	resp, err := s.client.Do(ctx, req, field)
	if err != nil {
		return nil, resp, err
	}

	return field, resp, nil
}

// Create creates a new global issue field.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/concepts/issues/create-global-field.html
func (s *FieldsService) Create(ctx context.Context, field *FieldCreateRequest) (*Field, *Response, error) {
	req, err := s.client.NewRequest("POST", "v3/fields", field)
	if err != nil {
		return nil, nil, err
	}

	f := new(Field)
	resp, err := s.client.Do(ctx, req, f)
	if err != nil {
		return nil, resp, err
	}

	return f, resp, nil
}

// Edit updates an existing global issue field.
// The opts parameter must include Version for optimistic locking.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/concepts/issues/edit-global-field.html
func (s *FieldsService) Edit(ctx context.Context, fieldID string, field *FieldEditRequest, opts *FieldEditOptions) (*Field, *Response, error) {
	u := fmt.Sprintf("v3/fields/%v", fieldID)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("PATCH", u, field)
	if err != nil {
		return nil, nil, err
	}

	f := new(Field)
	resp, err := s.client.Do(ctx, req, f)
	if err != nil {
		return nil, resp, err
	}

	return f, resp, nil
}
