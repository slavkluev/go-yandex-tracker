package tracker

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestIssuesService_ListAttachments(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v2/issues/{key}/attachments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{
			"self": "https://api.tracker.yandex.net/v2/issues/QUEUE-1/attachments/1",
			"id": "1",
			"name": "test.txt",
			"content": "https://tracker.yandex.net/issues/QUEUE-1/attachments/1/test.txt",
			"mimetype": "text/plain",
			"size": 12345,
			"createdBy": {"self": "https://api.tracker.yandex.net/v2/users/1", "id": "1", "display": "Test User"},
			"createdAt": "2024-01-15T10:00:00.000+0000"
		}]`)
	})

	ctx := context.Background()
	attachments, _, err := client.Issues.ListAttachments(ctx, "QUEUE-1")
	if err != nil {
		t.Fatalf("ListAttachments returned error: %v", err)
	}

	if len(attachments) != 1 {
		t.Fatalf("ListAttachments returned %d attachments, want 1", len(attachments))
	}

	a := attachments[0]
	if got := *a.Name; got != "test.txt" {
		t.Errorf("Name = %q, want %q", got, "test.txt")
	}
	if got := *a.Size; got != 12345 {
		t.Errorf("Size = %d, want %d", got, 12345)
	}
	if got := *a.Mimetype; got != "text/plain" {
		t.Errorf("Mimetype = %q, want %q", got, "text/plain")
	}
	if got := *a.ID; got != "1" {
		t.Errorf("ID = %q, want %q", got, "1")
	}
}

func TestIssuesService_UploadAttachment(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v2/issues/{key}/attachments/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		ct := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ct, "multipart/form-data") {
			t.Errorf("Content-Type = %q, want prefix %q", ct, "multipart/form-data")
		}

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			t.Fatalf("ParseMultipartForm returned error: %v", err)
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			t.Fatalf("FormFile(\"file\") returned error: %v", err)
		}
		defer file.Close()

		if header.Filename != "test.txt" {
			t.Errorf("filename = %q, want %q", header.Filename, "test.txt")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v2/issues/QUEUE-1/attachments/1",
			"id": "1",
			"name": "test.txt",
			"mimetype": "text/plain",
			"size": 11
		}`)
	})

	ctx := context.Background()
	attachment, resp, err := client.Issues.UploadAttachment(ctx, "QUEUE-1", "test.txt", strings.NewReader("hello world"))
	if err != nil {
		t.Fatalf("UploadAttachment returned error: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusCreated)
	}

	if got := *attachment.Name; got != "test.txt" {
		t.Errorf("Name = %q, want %q", got, "test.txt")
	}
	if got := *attachment.Size; got != 11 {
		t.Errorf("Size = %d, want %d", got, 11)
	}
}

func TestIssuesService_UploadAttachment_VerifyContent(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v2/issues/{key}/attachments/", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			t.Fatalf("ParseMultipartForm returned error: %v", err)
		}

		file, _, err := r.FormFile("file")
		if err != nil {
			t.Fatalf("FormFile(\"file\") returned error: %v", err)
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			t.Fatalf("ReadAll returned error: %v", err)
		}
		if got := string(content); got != "hello world" {
			t.Errorf("file content = %q, want %q", got, "hello world")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"id": "1", "name": "test.txt"}`)
	})

	ctx := context.Background()
	_, _, err := client.Issues.UploadAttachment(ctx, "QUEUE-1", "test.txt", strings.NewReader("hello world"))
	if err != nil {
		t.Fatalf("UploadAttachment returned error: %v", err)
	}
}

func TestIssuesService_DeleteAttachment(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("DELETE /v2/issues/{key}/attachments/{id}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Issues.DeleteAttachment(ctx, "QUEUE-1", "123")
	if err != nil {
		t.Fatalf("DeleteAttachment returned error: %v", err)
	}
}

func TestIssuesService_UploadTempFile(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v2/attachments/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		ct := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ct, "multipart/form-data") {
			t.Errorf("Content-Type = %q, want prefix %q", ct, "multipart/form-data")
		}

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			t.Fatalf("ParseMultipartForm returned error: %v", err)
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			t.Fatalf("FormFile(\"file\") returned error: %v", err)
		}
		defer file.Close()

		if header.Filename != "temp.txt" {
			t.Errorf("filename = %q, want %q", header.Filename, "temp.txt")
		}

		content, err := io.ReadAll(file)
		if err != nil {
			t.Fatalf("ReadAll returned error: %v", err)
		}
		if got := string(content); got != "temp content" {
			t.Errorf("file content = %q, want %q", got, "temp content")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v2/attachments/2",
			"id": "2",
			"name": "temp.txt",
			"mimetype": "text/plain",
			"size": 12
		}`)
	})

	ctx := context.Background()
	attachment, resp, err := client.Issues.UploadTempFile(ctx, "temp.txt", strings.NewReader("temp content"))
	if err != nil {
		t.Fatalf("UploadTempFile returned error: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusCreated)
	}

	if got := *attachment.Name; got != "temp.txt" {
		t.Errorf("Name = %q, want %q", got, "temp.txt")
	}
	if got := *attachment.ID; got != "2" {
		t.Errorf("ID = %q, want %q", got, "2")
	}
}

func TestIssuesService_DownloadAttachment(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v2/issues/{key}/attachments/{id}/{filename}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", "11")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello world"))
	})

	ctx := context.Background()
	rc, _, err := client.Issues.DownloadAttachment(ctx, "QUEUE-1", "1", "test.txt")
	if err != nil {
		t.Fatalf("DownloadAttachment returned error: %v", err)
	}
	defer rc.Close()

	content, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("ReadAll returned error: %v", err)
	}
	if got := string(content); got != "hello world" {
		t.Errorf("content = %q, want %q", got, "hello world")
	}
}

func TestIssuesService_DownloadAttachment_Headers(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v2/issues/{key}/attachments/{id}/{filename}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Length", "11")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello world"))
	})

	ctx := context.Background()
	rc, resp, err := client.Issues.DownloadAttachment(ctx, "QUEUE-1", "1", "report.pdf")
	if err != nil {
		t.Fatalf("DownloadAttachment returned error: %v", err)
	}
	defer rc.Close()

	if got := resp.Header.Get("Content-Type"); got != "application/pdf" {
		t.Errorf("Content-Type = %q, want %q", got, "application/pdf")
	}
	if got := resp.Header.Get("Content-Length"); got != "11" {
		t.Errorf("Content-Length = %q, want %q", got, "11")
	}
}

func TestIssuesService_DownloadAttachment_MustClose(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v2/issues/{key}/attachments/{id}/{filename}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("data"))
	})

	ctx := context.Background()
	rc, _, err := client.Issues.DownloadAttachment(ctx, "QUEUE-1", "1", "file.bin")
	if err != nil {
		t.Fatalf("DownloadAttachment returned error: %v", err)
	}

	if err := rc.Close(); err != nil {
		t.Errorf("Close() returned error: %v", err)
	}
}
