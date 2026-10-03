package tracker

import (
	"context"
	"iter"
)

// List returns one page of statuses. Pass nil for opts to get the first page
// at the server's default page size; ListIter returns every page.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/admin/get-statuses
func (s *StatusesService) List(ctx context.Context, opts *ListOptions) ([]*Status, *Response, error) {
	u, err := addOptions("v3/statuses", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
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

// ListIter yields the statuses of every page, starting with the page List
// returns for opts, which it copies when called. On an error it yields the
// error once and stops.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/admin/get-statuses
func (s *StatusesService) ListIter(ctx context.Context, opts *ListOptions) iter.Seq2[*Status, error] {
	return pageIter(ctx, opts, s.List)
}
