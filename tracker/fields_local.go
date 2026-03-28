package tracker

import (
	"context"
	"fmt"
)

// ListLocal returns all local fields for the specified queue.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/concepts/queues/get-local-fields.html
func (s *FieldsService) ListLocal(ctx context.Context, queueKey string) ([]*Field, *Response, error) {
	u := fmt.Sprintf("v3/queues/%v/localFields", queueKey)

	req, err := s.client.NewRequest("GET", u, nil)
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

// GetLocal fetches a single local field for the specified queue.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/concepts/queues/get-local-field.html
func (s *FieldsService) GetLocal(ctx context.Context, queueKey, fieldKey string) (*Field, *Response, error) {
	u := fmt.Sprintf("v3/queues/%v/localFields/%v", queueKey, fieldKey)

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

// CreateLocal creates a new local field for the specified queue.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/concepts/queues/create-local-field.html
func (s *FieldsService) CreateLocal(ctx context.Context, queueKey string, field *FieldCreateRequest) (*Field, *Response, error) {
	u := fmt.Sprintf("v3/queues/%v/localFields", queueKey)

	req, err := s.client.NewRequest("POST", u, field)
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

// EditLocal updates an existing local field for the specified queue.
// Unlike global Edit, local field edits do not require a version parameter.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/concepts/queues/edit-local-field.html
func (s *FieldsService) EditLocal(ctx context.Context, queueKey, fieldKey string, field *FieldEditRequest) (*Field, *Response, error) {
	u := fmt.Sprintf("v3/queues/%v/localFields/%v", queueKey, fieldKey)

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
