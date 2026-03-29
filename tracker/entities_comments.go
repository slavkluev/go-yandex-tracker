package tracker

import (
	"context"
	"fmt"
)

// ListComments returns the comments for an entity.
// Use CommentListOptions to paginate via cursor (ID of last entry) and perPage.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/entities/comments/get-all-comments
func (s *EntitiesService) ListComments(ctx context.Context, entityType EntityType, id string, opts *CommentListOptions) ([]*Comment, *Response, error) {
	u := fmt.Sprintf("v3/entities/%v/%v/comments", entityType, id)

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

// CreateComment adds a new comment to an entity.
// Use EntityCommentCreateOptions to control notification behavior.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/entities/comments/add-comment
func (s *EntitiesService) CreateComment(ctx context.Context, entityType EntityType, id string, comment *CommentRequest, opts *EntityCommentCreateOptions) (*Comment, *Response, error) {
	u := fmt.Sprintf("v3/entities/%v/%v/comments", entityType, id)

	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

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

// GetComment retrieves a single comment on an entity by its numeric ID.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/entities/comments/get-comment
func (s *EntitiesService) GetComment(ctx context.Context, entityType EntityType, id string, commentID int) (*Comment, *Response, error) {
	u := fmt.Sprintf("v3/entities/%v/%v/comments/%d", entityType, id, commentID)

	req, err := s.client.NewRequest("GET", u, nil)
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

// EditComment updates an existing comment on an entity.
// The commentID parameter is numeric (int) because Comment.ID is *int.
// Use EntityCommentCreateOptions to control notification behavior.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/entities/comments/patch-comment
func (s *EntitiesService) EditComment(ctx context.Context, entityType EntityType, id string, commentID int, comment *CommentRequest, opts *EntityCommentCreateOptions) (*Comment, *Response, error) {
	u := fmt.Sprintf("v3/entities/%v/%v/comments/%d", entityType, id, commentID)

	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

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

// DeleteComment deletes a comment from an entity.
// The commentID parameter is numeric (int) because Comment.ID is *int.
// Returns (*Response, error) since 204 responses have no body.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/entities/comments/delete-comment
func (s *EntitiesService) DeleteComment(ctx context.Context, entityType EntityType, id string, commentID int) (*Response, error) {
	u := fmt.Sprintf("v3/entities/%v/%v/comments/%d", entityType, id, commentID)

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
