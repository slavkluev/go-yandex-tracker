package tracker

import (
	"context"
	"iter"
)

// List returns one page of resolutions. Pass nil for opts to get the first page
// at the server's default page size; ListIter returns every page.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/admin/get-resolutions
func (s *ResolutionsService) List(ctx context.Context, opts *ListOptions) ([]*Resolution, *Response, error) {
	u, err := addOptions("v3/resolutions", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var resolutions []*Resolution
	resp, err := s.client.Do(ctx, req, &resolutions)
	if err != nil {
		return nil, resp, err
	}

	return resolutions, resp, nil
}

// ListIter yields the resolutions of every page, starting with the page List
// returns for opts, which it copies when called. On an error it yields the
// error once and stops.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/admin/get-resolutions
func (s *ResolutionsService) ListIter(ctx context.Context, opts *ListOptions) iter.Seq2[*Resolution, error] {
	return pageIter(ctx, opts, s.List)
}
