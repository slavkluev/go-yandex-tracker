package tracker

import (
	"context"
	"fmt"
)

// ListFields returns the fields that a queue configures for its issues:
// Field.Schema.Required says whether the queue requires each one, and is nil
// when Tracker omits it (as it does for array fields), and
// Field.OptionsProvider can carry per-queue values (QueueValues) with
// fallback Defaults. It is often a subset of the issue fields, and it can be
// empty. It is not the queue's local fields, which FieldsService.ListLocal
// returns.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api/queues/get-fields
func (s *QueuesService) ListFields(ctx context.Context, queueKey string) ([]*Field, *Response, error) {
	u := fmt.Sprintf("v3/queues/%v/fields", queueKey)

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
