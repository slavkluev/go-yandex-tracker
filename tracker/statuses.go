package tracker

import "context"

// List returns all statuses.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/admin/get-statuses
func (s *StatusesService) List(ctx context.Context) ([]*Status, *Response, error) {
	req, err := s.client.NewRequest("GET", "v3/statuses", nil)
	if err != nil {
		return nil, nil, err
	}

	var statuses []*Status
	resp, err := s.client.Do(ctx, req, &statuses)
	if err != nil {
		return nil, resp, err
	}

	return statuses, resp, nil
}
