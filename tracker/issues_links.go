package tracker

import (
	"context"
	"fmt"
)

// GetLinks returns the links for an issue.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/issues/get-links
func (s *IssuesService) GetLinks(ctx context.Context, issueKey string) ([]*IssueLink, *Response, error) {
	u := fmt.Sprintf("v3/issues/%v/links", issueKey)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var links []*IssueLink
	resp, err := s.client.Do(ctx, req, &links)
	if err != nil {
		return nil, resp, err
	}

	return links, resp, nil
}

// CreateLink creates a link between two issues.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/issues/link-issue
func (s *IssuesService) CreateLink(ctx context.Context, issueKey string, link *LinkRequest) (*IssueLink, *Response, error) {
	u := fmt.Sprintf("v3/issues/%v/links", issueKey)

	req, err := s.client.NewRequest("POST", u, link)
	if err != nil {
		return nil, nil, err
	}

	l := new(IssueLink)
	resp, err := s.client.Do(ctx, req, l)
	if err != nil {
		return nil, resp, err
	}

	return l, resp, nil
}

// DeleteLink deletes a link from an issue.
// Returns (*Response, error) since 204 responses have no body.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/issues/delete-link-issue
func (s *IssuesService) DeleteLink(ctx context.Context, issueKey, linkID string) (*Response, error) {
	u := fmt.Sprintf("v3/issues/%v/links/%v", issueKey, linkID)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
