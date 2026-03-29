package tracker

import (
	"context"
	"fmt"
)

// ListComments returns the comments for an issue.
// Use CommentListOptions to paginate via cursor (ID of last entry) and perPage.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/issues/get-comments
func (s *IssuesService) ListComments(ctx context.Context, issueKey string, opts *CommentListOptions) ([]*Comment, *Response, error) {
	u := fmt.Sprintf("v3/issues/%v/comments", issueKey)

	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var comments []*Comment
	resp, err := s.client.Do(ctx, req, &comments)
	if err != nil {
		return nil, resp, err
	}

	return comments, resp, nil
}

// CreateComment adds a new comment to an issue.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/issues/add-comment
func (s *IssuesService) CreateComment(ctx context.Context, issueKey string, comment *CommentRequest) (*Comment, *Response, error) {
	u := fmt.Sprintf("v3/issues/%v/comments", issueKey)

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

// EditComment updates an existing comment on an issue.
// The commentID parameter is numeric (int) because Comment.ID is *int.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/issues/edit-comment
func (s *IssuesService) EditComment(ctx context.Context, issueKey string, commentID int, comment *CommentRequest) (*Comment, *Response, error) {
	u := fmt.Sprintf("v3/issues/%v/comments/%d", issueKey, commentID)

	req, err := s.client.NewRequest("PATCH", u, comment)
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

// DeleteComment deletes a comment from an issue.
// The commentID parameter is numeric (int) because Comment.ID is *int.
// Returns (*Response, error) since 204 responses have no body.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/issues/delete-comment
func (s *IssuesService) DeleteComment(ctx context.Context, issueKey string, commentID int) (*Response, error) {
	u := fmt.Sprintf("v3/issues/%v/comments/%d", issueKey, commentID)

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
