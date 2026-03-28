package tracker

import (
	"context"
	"fmt"
)

// ListApplications returns the list of external applications registered in Yandex Tracker.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/api-ref/external-links/
func (s *ExternalLinksService) ListApplications(ctx context.Context) ([]*ExternalApplication, *Response, error) {
	req, err := s.client.NewRequest("GET", "v2/applications", nil)
	if err != nil {
		return nil, nil, err
	}

	var apps []*ExternalApplication
	resp, err := s.client.Do(ctx, req, &apps)
	if err != nil {
		return nil, resp, err
	}

	return apps, resp, nil
}

// ListLinks returns the external links for an issue.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/api-ref/external-links/
func (s *ExternalLinksService) ListLinks(ctx context.Context, issueKey string) ([]*ExternalLink, *Response, error) {
	u := fmt.Sprintf("v2/issues/%v/remotelinks", issueKey)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var links []*ExternalLink
	resp, err := s.client.Do(ctx, req, &links)
	if err != nil {
		return nil, resp, err
	}

	return links, resp, nil
}

// CreateLink creates an external link on an issue.
// The opts parameter allows specifying the backlink query parameter to create a reciprocal link.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/api-ref/external-links/
func (s *ExternalLinksService) CreateLink(ctx context.Context, issueKey string, link *ExternalLinkCreateRequest, opts *ExternalLinkCreateOptions) (*ExternalLink, *Response, error) {
	u := fmt.Sprintf("v2/issues/%v/remotelinks", issueKey)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("POST", u, link)
	if err != nil {
		return nil, nil, err
	}

	el := new(ExternalLink)
	resp, err := s.client.Do(ctx, req, el)
	if err != nil {
		return nil, resp, err
	}

	return el, resp, nil
}

// DeleteLink deletes an external link from an issue.
// The API returns 204 No Content on success.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/api-ref/external-links/
func (s *ExternalLinksService) DeleteLink(ctx context.Context, issueKey string, linkID int) (*Response, error) {
	u := fmt.Sprintf("v2/issues/%v/remotelinks/%v", issueKey, linkID)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
