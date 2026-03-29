package tracker

import (
	"context"
	"fmt"
)

// ListWorklogs returns the worklog entries for an issue.
// The API does not support pagination for worklog listings.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/issues/issue-worklog
func (s *IssuesService) ListWorklogs(ctx context.Context, issueKey string) ([]*Worklog, *Response, error) {
	u := fmt.Sprintf("v2/issues/%v/worklog", issueKey)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var worklogs []*Worklog
	resp, err := s.client.Do(ctx, req, &worklogs)
	if err != nil {
		return nil, resp, err
	}

	return worklogs, resp, nil
}

// CreateWorklog adds a worklog entry to an issue.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/issues/new-worklog
func (s *IssuesService) CreateWorklog(ctx context.Context, issueKey string, worklog *WorklogRequest) (*Worklog, *Response, error) {
	u := fmt.Sprintf("v2/issues/%v/worklog", issueKey)

	req, err := s.client.NewRequest("POST", u, worklog)
	if err != nil {
		return nil, nil, err
	}

	wl := new(Worklog)
	resp, err := s.client.Do(ctx, req, wl)
	if err != nil {
		return nil, resp, err
	}

	return wl, resp, nil
}

// EditWorklog updates a worklog entry on an issue.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/issues/patch-worklog
func (s *IssuesService) EditWorklog(ctx context.Context, issueKey, worklogID string, worklog *WorklogRequest) (*Worklog, *Response, error) {
	u := fmt.Sprintf("v2/issues/%v/worklog/%v", issueKey, worklogID)

	req, err := s.client.NewRequest("PATCH", u, worklog)
	if err != nil {
		return nil, nil, err
	}

	wl := new(Worklog)
	resp, err := s.client.Do(ctx, req, wl)
	if err != nil {
		return nil, resp, err
	}

	return wl, resp, nil
}

// DeleteWorklog deletes a worklog entry from an issue.
// The API returns HTTP 204 No Content on success.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/issues/delete-worklog
func (s *IssuesService) DeleteWorklog(ctx context.Context, issueKey, worklogID string) (*Response, error) {
	u := fmt.Sprintf("v2/issues/%v/worklog/%v", issueKey, worklogID)

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
