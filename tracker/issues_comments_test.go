package tracker

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestIssuesService_ListComments(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/issues/{key}/comments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{
			"id": 123,
			"text": "Test comment",
			"createdBy": {"id": "user1", "display": "User One"}
		}]`)
	})

	ctx := context.Background()
	comments, _, err := client.Issues.ListComments(ctx, "QUEUE-1", nil)
	if err != nil {
		t.Fatalf("ListComments returned error: %v", err)
	}

	if len(comments) != 1 {
		t.Fatalf("ListComments returned %d comments, want 1", len(comments))
	}

	comment := comments[0]
	if got := *comment.ID; got != 123 {
		t.Errorf("ID = %d, want %d", got, 123)
	}
	if got := *comment.Text; got != "Test comment" {
		t.Errorf("Text = %q, want %q", got, "Test comment")
	}
	if got := *comment.CreatedBy.ID; got != "user1" {
		t.Errorf("CreatedBy.ID = %q, want %q", got, "user1")
	}
}

func TestIssuesService_ListComments_WithOptions(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/issues/{key}/comments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		q := r.URL.Query()
		if got := q.Get("perPage"); got != "10" {
			t.Errorf("perPage = %q, want %q", got, "10")
		}
		if got := q.Get("id"); got != "abc" {
			t.Errorf("id = %q, want %q", got, "abc")
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[]`)
	})

	ctx := context.Background()
	opts := &CommentListOptions{ID: "abc", PerPage: 10}
	_, _, err := client.Issues.ListComments(ctx, "QUEUE-1", opts)
	if err != nil {
		t.Fatalf("ListComments returned error: %v", err)
	}
}

func TestIssuesService_CreateComment(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/issues/{key}/comments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"text":"New comment"}`)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"id": 456,
			"text": "New comment",
			"createdBy": {"id": "user1", "display": "User One"}
		}`)
	})

	ctx := context.Background()
	comment, resp, err := client.Issues.CreateComment(ctx, "QUEUE-1", &CommentRequest{
		Text: Ptr("New comment"),
	})
	if err != nil {
		t.Fatalf("CreateComment returned error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusCreated)
	}
	if got := *comment.ID; got != 456 {
		t.Errorf("ID = %d, want %d", got, 456)
	}
	if got := *comment.Text; got != "New comment" {
		t.Errorf("Text = %q, want %q", got, "New comment")
	}
}

func TestIssuesService_EditComment(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("PATCH /v3/issues/{key}/comments/{id}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testBody(t, r, `{"text":"Updated comment"}`)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"id": 123,
			"text": "Updated comment",
			"version": 2
		}`)
	})

	ctx := context.Background()
	comment, _, err := client.Issues.EditComment(ctx, "QUEUE-1", 123, &CommentRequest{
		Text: Ptr("Updated comment"),
	})
	if err != nil {
		t.Fatalf("EditComment returned error: %v", err)
	}
	if got := *comment.ID; got != 123 {
		t.Errorf("ID = %d, want %d", got, 123)
	}
	if got := *comment.Text; got != "Updated comment" {
		t.Errorf("Text = %q, want %q", got, "Updated comment")
	}
	if got := *comment.Version; got != 2 {
		t.Errorf("Version = %d, want %d", got, 2)
	}
}

func TestIssuesService_DeleteComment(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("DELETE /v3/issues/{key}/comments/{id}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	resp, err := client.Issues.DeleteComment(ctx, "QUEUE-1", 123)
	if err != nil {
		t.Fatalf("DeleteComment returned error: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusNoContent)
	}
}
