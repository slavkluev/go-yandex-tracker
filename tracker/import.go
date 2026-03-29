package tracker

import (
	"context"
	"fmt"
	"io"
)

// ImportIssue imports an issue with preserved creation timestamps.
// The createdAt and createdBy fields in the request are used as the issue's actual creation metadata,
// enabling historical data migration.
//
// Requires admin permissions in the target queue.
// createdAt must not exceed current time; updatedAt must fall between createdAt and current time.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/import/import-ticket
func (s *ImportService) ImportIssue(ctx context.Context, issue *ImportIssueRequest) (*Issue, *Response, error) {
	req, err := s.client.NewRequest("POST", "v3/issues/_import", issue)
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

// ImportComment imports a comment to an issue with preserved creation timestamps.
// Comments must be imported in chronological order of createdAt.
//
// Requires admin permissions in the target queue.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/import/import-comments
func (s *ImportService) ImportComment(ctx context.Context, issueKey string, comment *ImportCommentRequest) (*Comment, *Response, error) {
	u := fmt.Sprintf("v3/issues/%v/comments/_import", issueKey)

	req, err := s.client.NewRequest("POST", u, comment)
	if err != nil {
		return nil, nil, err
	}

	c := new(Comment)
	resp, err := s.client.Do(ctx, req, c)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, nil
}

// ImportLink imports a link to an issue with preserved creation timestamps.
//
// Requires admin permissions in the target queue.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/import/import-links
func (s *ImportService) ImportLink(ctx context.Context, issueKey string, link *ImportLinkRequest) (*IssueLink, *Response, error) {
	u := fmt.Sprintf("v3/issues/%v/links/_import", issueKey)

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

// ImportFile imports a file attachment to an issue with preserved creation metadata.
// The file metadata (filename, createdAt, createdBy) is passed as query parameters.
// The file content is sent as multipart/form-data in the request body.
//
// Requires admin permissions in the target queue.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/import/import-attachments
func (s *ImportService) ImportFile(ctx context.Context, issueKey string, opts *ImportFileOptions, filename string, file io.Reader) (*Attachment, *Response, error) {
	u := fmt.Sprintf("v3/issues/%v/attachments/_import", issueKey)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewUploadRequest("POST", u, filename, file)
	if err != nil {
		return nil, nil, err
	}

	attachment := new(Attachment)
	resp, err := s.client.Do(ctx, req, attachment)
	if err != nil {
		return nil, resp, err
	}

	return attachment, resp, nil
}

// ImportCommentFile imports a file attachment to a comment with preserved creation metadata.
// The file metadata (filename, createdAt, createdBy) is passed as query parameters.
// The file content is sent as multipart/form-data in the request body.
//
// Requires admin permissions in the target queue.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/import/import-attachments
func (s *ImportService) ImportCommentFile(ctx context.Context, issueKey string, commentID int, opts *ImportFileOptions, filename string, file io.Reader) (*Attachment, *Response, error) {
	u := fmt.Sprintf("v3/issues/%v/comments/%v/attachments/_import", issueKey, commentID)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewUploadRequest("POST", u, filename, file)
	if err != nil {
		return nil, nil, err
	}

	attachment := new(Attachment)
	resp, err := s.client.Do(ctx, req, attachment)
	if err != nil {
		return nil, resp, err
	}

	return attachment, resp, nil
}
