package tracker

import (
	"context"
	"fmt"
	"io"
)

// ListAttachments returns a list of attachments on an issue.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/concepts/issues/get-attachments-list.html
func (s *IssuesService) ListAttachments(ctx context.Context, issueKey string) ([]*Attachment, *Response, error) {
	u := fmt.Sprintf("v2/issues/%v/attachments", issueKey)

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

// UploadAttachment uploads a file to an issue using multipart/form-data.
// The filename is used in the multipart Content-Disposition header.
// The file content is read from the provided io.Reader.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/concepts/issues/post-attachment.html
func (s *IssuesService) UploadAttachment(ctx context.Context, issueKey, filename string, file io.Reader) (*Attachment, *Response, error) {
	u := fmt.Sprintf("v2/issues/%v/attachments/", issueKey)

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

// DeleteAttachment deletes an attachment from an issue.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/concepts/issues/delete-attachment.html
func (s *IssuesService) DeleteAttachment(ctx context.Context, issueKey, attachmentID string) (*Response, error) {
	u := fmt.Sprintf("v2/issues/%v/attachments/%v", issueKey, attachmentID)

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

// UploadTempFile uploads a temporary file not attached to any issue.
// Temporary files can later be attached to an issue using their ID.
// The filename is used in the multipart Content-Disposition header.
// The file content is read from the provided io.Reader.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/concepts/issues/temp-attachment.html
func (s *IssuesService) UploadTempFile(ctx context.Context, filename string, file io.Reader) (*Attachment, *Response, error) {
	u := "v2/attachments/"

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

// DownloadAttachment downloads an attachment as a streaming response.
// The caller is responsible for closing the returned io.ReadCloser.
// Content-Type and Content-Length are available via resp.Header.Get().
//
// The filename parameter is required in the URL path; the API returns 404
// without it. Use the Attachment.Name field from a prior ListAttachments call.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/concepts/issues/get-attachment.html
func (s *IssuesService) DownloadAttachment(ctx context.Context, issueKey, attachmentID, filename string) (io.ReadCloser, *Response, error) {
	u := fmt.Sprintf("v2/issues/%v/attachments/%v/%v", issueKey, attachmentID, filename)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.BareDo(ctx, req)
	if err != nil {
		return nil, resp, err
	}

	return resp.Body, resp, nil
}
