package tracker

import (
	"fmt"
	"net/http"
	"testing"
)

func TestEntitiesService_ListAttachments(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/entities/project/{id}/attachments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{
			"self": "https://api.tracker.yandex.net/v3/attachments/file-1",
			"id": "file-1",
			"name": "report.pdf",
			"mimetype": "application/pdf",
			"size": 12345,
			"createdBy": {"id": "user1"}
		}]`)
	})

	attachments, _, err := client.Entities.ListAttachments(ctx, EntityTypeProject, "1")
	if err != nil {
		t.Fatalf("ListAttachments returned error: %v", err)
	}

	if len(attachments) != 1 {
		t.Fatalf("ListAttachments returned %d attachments, want 1", len(attachments))
	}

	a := attachments[0]
	if got := *a.ID; got != "file-1" {
		t.Errorf("ID = %q, want %q", got, "file-1")
	}
	if got := *a.Name; got != "report.pdf" {
		t.Errorf("Name = %q, want %q", got, "report.pdf")
	}
	if got := *a.Mimetype; got != "application/pdf" {
		t.Errorf("Mimetype = %q, want %q", got, "application/pdf")
	}
	if got := *a.Size; got != 12345 {
		t.Errorf("Size = %d, want %d", got, 12345)
	}
	if got := *a.CreatedBy.ID; got != "user1" {
		t.Errorf("CreatedBy.ID = %q, want %q", got, "user1")
	}
}

func TestEntitiesService_AttachFile(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/entities/project/{id}/attachments/{fileId}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		// Verify NO body is sent (temp-file-ID pattern, not multipart)
		if r.ContentLength > 0 {
			t.Errorf("ContentLength = %d, want 0 or -1 (no body)", r.ContentLength)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/entities/project/1",
			"id": "1",
			"entityType": "project",
			"fields": {"summary": "Test Project"}
		}`)
	})

	entity, _, err := client.Entities.AttachFile(ctx, EntityTypeProject, "1", "temp-file-123")
	if err != nil {
		t.Fatalf("AttachFile returned error: %v", err)
	}

	if got := *entity.ID; got != "1" {
		t.Errorf("ID = %q, want %q", got, "1")
	}
	if got := *entity.EntityType; got != "project" {
		t.Errorf("EntityType = %q, want %q", got, "project")
	}
	if entity.Fields == nil {
		t.Fatal("Fields is nil, want non-nil")
	}
	if got := *entity.Fields.Summary; got != "Test Project" {
		t.Errorf("Fields.Summary = %q, want %q", got, "Test Project")
	}
}

func TestEntitiesService_DeleteAttachment(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("DELETE /v3/entities/project/{id}/attachments/{fileId}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.Entities.DeleteAttachment(ctx, EntityTypeProject, "1", "file-1")
	if err != nil {
		t.Fatalf("DeleteAttachment returned error: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusNoContent)
	}
}
