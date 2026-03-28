package tracker

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestImportService_ImportIssue(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v2/issues/_import", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"queue":"QUEUE","summary":"Imported issue","createdAt":"2017-08-29T12:34:41.740+0000","createdBy":"user1"}`)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v2/issues/QUEUE-1",
			"id": "1",
			"key": "QUEUE-1",
			"summary": "Imported issue"
		}`)
	})

	ctx := context.Background()
	issue, _, err := client.Import.ImportIssue(ctx, &ImportIssueRequest{
		Queue:     Ptr("QUEUE"),
		Summary:   Ptr("Imported issue"),
		CreatedAt: Ptr("2017-08-29T12:34:41.740+0000"),
		CreatedBy: Ptr("user1"),
	})
	if err != nil {
		t.Fatalf("ImportIssue returned error: %v", err)
	}

	if got := *issue.Key; got != "QUEUE-1" {
		t.Errorf("Key = %q, want %q", got, "QUEUE-1")
	}
	if got := *issue.Summary; got != "Imported issue" {
		t.Errorf("Summary = %q, want %q", got, "Imported issue")
	}
}

func TestImportService_ImportComment(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v2/issues/{key}/comments/_import", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"text":"Imported comment","createdAt":"2017-08-29T12:34:41.740+0000","createdBy":"user1"}`)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v2/issues/QUEUE-1/comments/1",
			"id": 1,
			"text": "Imported comment",
			"createdBy": {"self": "https://api.tracker.yandex.net/v2/users/1", "id": "1", "display": "user1"},
			"createdAt": "2017-08-29T12:34:41.740+0000"
		}`)
	})

	ctx := context.Background()
	comment, _, err := client.Import.ImportComment(ctx, "QUEUE-1", &ImportCommentRequest{
		Text:      Ptr("Imported comment"),
		CreatedAt: Ptr("2017-08-29T12:34:41.740+0000"),
		CreatedBy: Ptr("user1"),
	})
	if err != nil {
		t.Fatalf("ImportComment returned error: %v", err)
	}

	if got := *comment.Text; got != "Imported comment" {
		t.Errorf("Text = %q, want %q", got, "Imported comment")
	}
	if got := *comment.ID; got != 1 {
		t.Errorf("ID = %d, want %d", got, 1)
	}
}

func TestImportService_ImportLink(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v2/issues/{key}/links/_import", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"relationship":"relates","issue":"QUEUE-2","createdAt":"2017-08-29T12:34:41.740+0000","createdBy":"user1"}`)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v2/issues/QUEUE-1/links/1",
			"id": "1",
			"type": {"self": "https://api.tracker.yandex.net/v2/linktypes/relates", "id": "relates", "inward": "relates", "outward": "relates"},
			"direction": "outward",
			"object": {"self": "https://api.tracker.yandex.net/v2/issues/QUEUE-2", "id": "2", "key": "QUEUE-2", "summary": "Linked issue"}
		}`)
	})

	ctx := context.Background()
	link, _, err := client.Import.ImportLink(ctx, "QUEUE-1", &ImportLinkRequest{
		Relationship: Ptr("relates"),
		Issue:        Ptr("QUEUE-2"),
		CreatedAt:    Ptr("2017-08-29T12:34:41.740+0000"),
		CreatedBy:    Ptr("user1"),
	})
	if err != nil {
		t.Fatalf("ImportLink returned error: %v", err)
	}

	if got := *link.ID; got != "1" {
		t.Errorf("ID = %q, want %q", got, "1")
	}
	if got := *link.Direction; got != "outward" {
		t.Errorf("Direction = %q, want %q", got, "outward")
	}
	if got := *link.Object.Key; got != "QUEUE-2" {
		t.Errorf("Object.Key = %q, want %q", got, "QUEUE-2")
	}
}

func TestImportService_ImportFile(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v2/issues/{key}/attachments/_import", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		// Verify query parameters
		if got := r.URL.Query().Get("filename"); got != "test.txt" {
			t.Errorf("query param filename = %q, want %q", got, "test.txt")
		}
		if got := r.URL.Query().Get("createdAt"); got != "2017-08-29T12:34:41.740+0000" {
			t.Errorf("query param createdAt = %q, want %q", got, "2017-08-29T12:34:41.740+0000")
		}
		if got := r.URL.Query().Get("createdBy"); got != "user1" {
			t.Errorf("query param createdBy = %q, want %q", got, "user1")
		}

		// Verify multipart content type
		ct := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ct, "multipart/form-data") {
			t.Errorf("Content-Type = %q, want prefix %q", ct, "multipart/form-data")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v2/issues/QUEUE-1/attachments/1",
			"id": "1",
			"name": "test.txt",
			"mimetype": "text/plain",
			"size": 12
		}`)
	})

	ctx := context.Background()
	attachment, _, err := client.Import.ImportFile(ctx, "QUEUE-1", &ImportFileOptions{
		Filename:  "test.txt",
		CreatedAt: "2017-08-29T12:34:41.740+0000",
		CreatedBy: "user1",
	}, "test.txt", strings.NewReader("file content"))
	if err != nil {
		t.Fatalf("ImportFile returned error: %v", err)
	}

	if got := *attachment.Name; got != "test.txt" {
		t.Errorf("Name = %q, want %q", got, "test.txt")
	}
	if got := *attachment.ID; got != "1" {
		t.Errorf("ID = %q, want %q", got, "1")
	}
}

func TestImportService_ImportCommentFile(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v2/issues/{key}/comments/{commentId}/attachments/_import", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		// Verify query parameters
		if got := r.URL.Query().Get("filename"); got != "comment-file.txt" {
			t.Errorf("query param filename = %q, want %q", got, "comment-file.txt")
		}
		if got := r.URL.Query().Get("createdAt"); got != "2017-08-29T12:34:41.740+0000" {
			t.Errorf("query param createdAt = %q, want %q", got, "2017-08-29T12:34:41.740+0000")
		}
		if got := r.URL.Query().Get("createdBy"); got != "user1" {
			t.Errorf("query param createdBy = %q, want %q", got, "user1")
		}

		// Verify multipart content type
		ct := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ct, "multipart/form-data") {
			t.Errorf("Content-Type = %q, want prefix %q", ct, "multipart/form-data")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v2/issues/QUEUE-1/comments/123/attachments/2",
			"id": "2",
			"name": "comment-file.txt",
			"mimetype": "text/plain",
			"size": 12
		}`)
	})

	ctx := context.Background()
	attachment, _, err := client.Import.ImportCommentFile(ctx, "QUEUE-1", 123, &ImportFileOptions{
		Filename:  "comment-file.txt",
		CreatedAt: "2017-08-29T12:34:41.740+0000",
		CreatedBy: "user1",
	}, "comment-file.txt", strings.NewReader("file content"))
	if err != nil {
		t.Fatalf("ImportCommentFile returned error: %v", err)
	}

	if got := *attachment.Name; got != "comment-file.txt" {
		t.Errorf("Name = %q, want %q", got, "comment-file.txt")
	}
	if got := *attachment.ID; got != "2" {
		t.Errorf("ID = %q, want %q", got, "2")
	}
}
