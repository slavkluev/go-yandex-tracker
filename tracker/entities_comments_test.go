package tracker

import (
	"fmt"
	"net/http"
	"testing"
)

func TestEntitiesService_ListComments(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/entities/project/{id}/comments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{
			"id": 123,
			"text": "Entity comment",
			"createdBy": {"id": "user1", "display": "User One"}
		}]`)
	})

	comments, _, err := client.Entities.ListComments(ctx, EntityTypeProject, "1", nil)
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
	if got := *comment.Text; got != "Entity comment" {
		t.Errorf("Text = %q, want %q", got, "Entity comment")
	}
	if got := *comment.CreatedBy.ID; got != "user1" {
		t.Errorf("CreatedBy.ID = %q, want %q", got, "user1")
	}
}

func TestEntitiesService_ListComments_WithOptions(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/entities/project/{id}/comments", func(w http.ResponseWriter, r *http.Request) {
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

	opts := &CommentListOptions{ID: "abc", PerPage: 10}
	_, _, err := client.Entities.ListComments(ctx, EntityTypeProject, "1", opts)
	if err != nil {
		t.Fatalf("ListComments returned error: %v", err)
	}
}

func TestEntitiesService_CreateComment(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/entities/project/{id}/comments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"text":"New entity comment"}`)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"id": 456,
			"text": "New entity comment",
			"createdBy": {"id": "user1", "display": "User One"}
		}`)
	})

	comment, resp, err := client.Entities.CreateComment(ctx, EntityTypeProject, "1", &CommentRequest{
		Text: Ptr("New entity comment"),
	}, nil)
	if err != nil {
		t.Fatalf("CreateComment returned error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusCreated)
	}
	if got := *comment.ID; got != 456 {
		t.Errorf("ID = %d, want %d", got, 456)
	}
	if got := *comment.Text; got != "New entity comment" {
		t.Errorf("Text = %q, want %q", got, "New entity comment")
	}
}

func TestEntitiesService_CreateComment_WithOptions(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/entities/project/{id}/comments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		q := r.URL.Query()
		if got := q.Get("isAddToFollowers"); got != "false" {
			t.Errorf("isAddToFollowers = %q, want %q", got, "false")
		}
		if got := q.Get("notify"); got != "true" {
			t.Errorf("notify = %q, want %q", got, "true")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"id": 789, "text": "comment with opts"}`)
	})

	opts := &EntityCommentCreateOptions{
		IsAddToFollowers: Ptr(false),
		Notify:           Ptr(true),
	}
	comment, _, err := client.Entities.CreateComment(ctx, EntityTypeProject, "1", &CommentRequest{
		Text: Ptr("comment with opts"),
	}, opts)
	if err != nil {
		t.Fatalf("CreateComment returned error: %v", err)
	}
	if got := *comment.ID; got != 789 {
		t.Errorf("ID = %d, want %d", got, 789)
	}
}

func TestEntitiesService_GetComment(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/entities/project/{id}/comments/{commentId}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"id": 123,
			"text": "Single comment",
			"createdBy": {"id": "user1", "display": "User One"}
		}`)
	})

	comment, _, err := client.Entities.GetComment(ctx, EntityTypeProject, "1", 123)
	if err != nil {
		t.Fatalf("GetComment returned error: %v", err)
	}
	if got := *comment.ID; got != 123 {
		t.Errorf("ID = %d, want %d", got, 123)
	}
	if got := *comment.Text; got != "Single comment" {
		t.Errorf("Text = %q, want %q", got, "Single comment")
	}
}

func TestEntitiesService_EditComment(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("PATCH /v3/entities/project/{id}/comments/{commentId}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testBody(t, r, `{"text":"Updated entity comment"}`)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"id": 123,
			"text": "Updated entity comment",
			"version": 2
		}`)
	})

	comment, _, err := client.Entities.EditComment(ctx, EntityTypeProject, "1", 123, &CommentRequest{
		Text: Ptr("Updated entity comment"),
	}, nil)
	if err != nil {
		t.Fatalf("EditComment returned error: %v", err)
	}
	if got := *comment.ID; got != 123 {
		t.Errorf("ID = %d, want %d", got, 123)
	}
	if got := *comment.Text; got != "Updated entity comment" {
		t.Errorf("Text = %q, want %q", got, "Updated entity comment")
	}
	if got := *comment.Version; got != 2 {
		t.Errorf("Version = %d, want %d", got, 2)
	}
}

func TestEntitiesService_DeleteComment(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("DELETE /v3/entities/project/{id}/comments/{commentId}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.Entities.DeleteComment(ctx, EntityTypeProject, "1", 123)
	if err != nil {
		t.Fatalf("DeleteComment returned error: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusNoContent)
	}
}
