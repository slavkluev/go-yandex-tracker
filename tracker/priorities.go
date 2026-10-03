package tracker

import (
	"context"
	"iter"
)

// List returns one page of priorities. Pass nil for opts to get the first
// page at the server's default page size; ListIter returns every page.
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

// ListIter yields the priorities of every page, starting with the page List
// returns for opts, which it copies when called. On an error it yields the
// error once and stops.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/issues/get-priorities
func (s *PrioritiesService) ListIter(ctx context.Context, opts *PriorityListOptions) iter.Seq2[*Priority, error] {
	return pageIter(ctx, opts, s.List)
}
