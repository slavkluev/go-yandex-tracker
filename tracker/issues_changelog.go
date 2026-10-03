package tracker

import (
	"context"
	"fmt"
	"iter"
)

// GetChangelog returns the changelog for an issue.
// Use ChangelogOptions to paginate via cursor (ID of last entry) and perPage.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/issues/get-changelog
func (s *IssuesService) GetChangelog(ctx context.Context, issueKey string, opts *ChangelogOptions) ([]*Changelog, *Response, error) {
	u := fmt.Sprintf("v3/issues/%v/changelog", issueKey)

	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var changelog []*Changelog
	resp, err := s.client.Do(ctx, req, &changelog)
	if err != nil {
		return nil, resp, err
	}

	return changelog, resp, nil
}

// GetChangelogIter yields the changelog entries of every page, starting with
// the page GetChangelog returns for opts, which it copies when called. Each
// next page starts after the ID of the previous page's last entry. On an error
// it yields the error once and stops.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/issues/get-changelog
func (s *IssuesService) GetChangelogIter(ctx context.Context, issueKey string, opts *ChangelogOptions) iter.Seq2[*Changelog, error] {
	id := func(c *Changelog) string {
		if c == nil || c.ID == nil {
			return ""
		}
		return string(*c.ID)
	}
	return cursorIter(ctx, opts, id, func(ctx context.Context, o *ChangelogOptions) ([]*Changelog, *Response, error) {
		return s.GetChangelog(ctx, issueKey, o)
	})
}
