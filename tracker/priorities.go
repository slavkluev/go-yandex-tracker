package tracker

import "context"

// List returns all priorities.
// Pass nil for opts to use defaults.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/issues/get-priorities
func (s *PrioritiesService) List(ctx context.Context, opts *PriorityListOptions) ([]*Priority, *Response, error) {
	u := "v3/priorities"
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var priorities []*Priority
	resp, err := s.client.Do(ctx, req, &priorities)
	if err != nil {
		return nil, resp, err
	}

	return priorities, resp, nil
}
