package tracker

import (
	"context"
	"iter"
)

// List returns one page of issue types. Pass nil for opts to get the first page
// at the server's default page size; ListIter returns every page.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/admin/get-issue-types
func (s *IssueTypesService) List(ctx context.Context, opts *ListOptions) ([]*IssueType, *Response, error) {
	u, err := addOptions("v3/issuetypes", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var issueTypes []*IssueType
	resp, err := s.client.Do(ctx, req, &issueTypes)
	if err != nil {
		return nil, resp, err
	}

	return issueTypes, resp, nil
}

// ListIter yields the issue types of every page, starting with the page List
// returns for opts, which it copies when called. On an error it yields the
// error once and stops.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/admin/get-issue-types
func (s *IssueTypesService) ListIter(ctx context.Context, opts *ListOptions) iter.Seq2[*IssueType, error] {
	return pageIter(ctx, opts, s.List)
}
