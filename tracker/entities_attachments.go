package tracker

import (
	"context"
	"fmt"
)

// ListAttachments returns the list of attachments on an entity.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/ru/api-ref/entities/attachments/get-all-attachments
func (s *EntitiesService) ListAttachments(ctx context.Context, entityType EntityType, id string) ([]*Attachment, *Response, error) {
	u := fmt.Sprintf("v3/entities/%v/%v/attachments", entityType, id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var attachments []*Attachment
	resp, err := s.client.Do(ctx, req, &attachments)
	if err != nil {
		return nil, resp, err
	}

	return attachments, resp, nil
}

// AttachFile attaches a previously uploaded temporary file to an entity.
// The file must be uploaded first using IssuesService.UploadTempFile,
// then attached using the returned file ID.
//
// Unlike issue attachments which use multipart upload, entity attachments
// use a two-step pattern: upload a temp file, then attach it by ID.
// The fileID is passed in the URL path and no request body is sent.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/ru/api-ref/entities/attachments/add-attachment
func (s *EntitiesService) AttachFile(ctx context.Context, entityType EntityType, id, fileID string) (*Entity, *Response, error) {
	u := fmt.Sprintf("v3/entities/%v/%v/attachments/%v", entityType, id, fileID)

	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, nil, err
	}

	entity := new(Entity)
	resp, err := s.client.Do(ctx, req, entity)
	if err != nil {
		return nil, resp, err
	}

	return entity, resp, nil
}

// DeleteAttachment removes an attachment from an entity.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/ru/api-ref/entities/attachments/delete-attachment
func (s *EntitiesService) DeleteAttachment(ctx context.Context, entityType EntityType, id, fileID string) (*Response, error) {
	u := fmt.Sprintf("v3/entities/%v/%v/attachments/%v", entityType, id, fileID)

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
