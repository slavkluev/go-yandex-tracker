package tracker

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// IssueGetOptions specifies the optional parameters to the Get method.
type IssueGetOptions struct {
	Expand string `url:"expand,omitempty"`
}

// IssueEditOptions specifies the optional parameters to the Edit method.
type IssueEditOptions struct {
	Version int `url:"version,omitempty"`
}

// Create creates a new issue.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/issues/create-issue
func (s *IssuesService) Create(ctx context.Context, issue *IssueRequest) (*Issue, *Response, error) {
	req, err := s.client.NewRequest("POST", "v2/issues/", issue)
	if err != nil {
		return nil, nil, err
	}

	i := new(Issue)
	resp, err := s.client.Do(ctx, req, i)
	if err != nil {
		return nil, resp, err
	}

	return i, resp, nil
}

// Get fetches an issue by its key or ID.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/issues/get-issue
func (s *IssuesService) Get(ctx context.Context, issueKey string, opts *IssueGetOptions) (*Issue, *Response, error) {
	u := fmt.Sprintf("v2/issues/%v", issueKey)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	issue := new(Issue)
	resp, err := s.client.Do(ctx, req, issue)
	if err != nil {
		return nil, resp, err
	}

	return issue, resp, nil
}

// Edit updates an existing issue.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/issues/patch-issue
func (s *IssuesService) Edit(ctx context.Context, issueKey string, issue *IssueRequest, opts *IssueEditOptions) (*Issue, *Response, error) {
	u := fmt.Sprintf("v2/issues/%v", issueKey)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("PATCH", u, issue)
	if err != nil {
		return nil, nil, err
	}

	i := new(Issue)
	resp, err := s.client.Do(ctx, req, i)
	if err != nil {
		return nil, resp, err
	}

	return i, resp, nil
}

// Search finds issues matching the specified criteria.
// Results are returned using page-based pagination controlled by IssueSearchOptions.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/issues/search-issues
func (s *IssuesService) Search(ctx context.Context, search *IssueSearchRequest, opts *IssueSearchOptions) ([]*Issue, *Response, error) {
	u := "v2/issues/_search"
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("POST", u, search)
	if err != nil {
		return nil, nil, err
	}

	var issues []*Issue
	resp, err := s.client.Do(ctx, req, &issues)
	if err != nil {
		return nil, resp, err
	}

	return issues, resp, nil
}

// ScrollSearch starts a scroll-based search for issues.
// Use ScrollNext to retrieve subsequent pages.
// The Response will contain ScrollID and ScrollToken for continuation.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/issues/search-issues
func (s *IssuesService) ScrollSearch(ctx context.Context, search *IssueSearchRequest, opts *ScrollSearchOptions) ([]*Issue, *Response, error) {
	u := "v2/issues/_search"
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("POST", u, search)
	if err != nil {
		return nil, nil, err
	}

	var issues []*Issue
	resp, err := s.client.Do(ctx, req, &issues)
	if err != nil {
		return nil, resp, err
	}

	return issues, resp, nil
}

// ScrollNext retrieves the next page of scroll search results.
// Pass the scrollID and scrollToken from the previous Response.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/issues/search-issues
func (s *IssuesService) ScrollNext(ctx context.Context, search *IssueSearchRequest, scrollID, scrollToken string) ([]*Issue, *Response, error) {
	u := fmt.Sprintf("v2/issues/_search?scrollId=%s&scrollToken=%s", scrollID, scrollToken)

	req, err := s.client.NewRequest("POST", u, search)
	if err != nil {
		return nil, nil, err
	}

	var issues []*Issue
	resp, err := s.client.Do(ctx, req, &issues)
	if err != nil {
		return nil, resp, err
	}

	return issues, resp, nil
}

// Count returns the number of issues matching the search criteria.
// The Yandex Tracker API returns a plain integer in the response body
// (not JSON), so this method uses BareDo and parses the integer manually.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/issues/count-issues
func (s *IssuesService) Count(ctx context.Context, search *IssueSearchRequest) (int, *Response, error) {
	req, err := s.client.NewRequest("POST", "v2/issues/_count", search)
	if err != nil {
		return 0, nil, err
	}

	resp, err := s.client.BareDo(ctx, req)
	if err != nil {
		return 0, resp, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, resp, err
	}

	count, err := strconv.Atoi(strings.TrimSpace(string(body)))
	if err != nil {
		return 0, resp, err
	}

	return count, resp, nil
}

// Move moves an issue to another queue.
// The target queue is specified in opts.Queue (required).
// An optional issue parameter can override fields during the move.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/issues/move-issue
func (s *IssuesService) Move(ctx context.Context, issueKey string, opts *IssueMoveOptions, issue *IssueRequest) (*Issue, *Response, error) {
	u := fmt.Sprintf("v2/issues/%v/_move", issueKey)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("POST", u, issue)
	if err != nil {
		return nil, nil, err
	}

	i := new(Issue)
	resp, err := s.client.Do(ctx, req, i)
	if err != nil {
		return nil, resp, err
	}

	return i, resp, nil
}
