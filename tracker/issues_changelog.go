package tracker

import (
	"context"
	"fmt"
)

// GetChangelog returns the changelog for an issue.
// Use ChangelogOptions to paginate via cursor (ID of last entry) and perPage.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/issues/get-changelog
func (s *IssuesService) GetChangelog(ctx context.Context, issueKey string, opts *ChangelogOptions) ([]*Changelog, *Response, error) {
	u := fmt.Sprintf("v2/issues/%v/changelog", issueKey)

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
