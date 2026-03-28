package tracker

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestIssuesService_GetChangelog(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v2/issues/{key}/changelog", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{
			"id": "1",
			"type": "IssueUpdated",
			"updatedAt": "2017-06-11T05:16:01.339+0000",
			"updatedBy": {"id": "user1"},
			"fields": [{"field": "status", "from": "open", "to": "closed"}]
		}]`)
	})

	ctx := context.Background()
	changelog, _, err := client.Issues.GetChangelog(ctx, "QUEUE-1", nil)
	if err != nil {
		t.Fatalf("GetChangelog returned error: %v", err)
	}

	if len(changelog) != 1 {
		t.Fatalf("GetChangelog returned %d entries, want 1", len(changelog))
	}

	entry := changelog[0]
	if got := *entry.ID; got != "1" {
		t.Errorf("ID = %q, want %q", got, "1")
	}
	if got := *entry.Type; got != "IssueUpdated" {
		t.Errorf("Type = %q, want %q", got, "IssueUpdated")
	}
	if len(entry.Fields) != 1 {
		t.Fatalf("Fields length = %d, want 1", len(entry.Fields))
	}
	if got := *entry.Fields[0].Field; got != "status" {
		t.Errorf("Fields[0].Field = %q, want %q", got, "status")
	}
}

func TestIssuesService_GetChangelog_WithOptions(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v2/issues/{key}/changelog", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		gotID := r.URL.Query().Get("id")
		if gotID != "last123" {
			t.Errorf("query param id = %q, want %q", gotID, "last123")
		}
		gotPerPage := r.URL.Query().Get("perPage")
		if gotPerPage != "50" {
			t.Errorf("query param perPage = %q, want %q", gotPerPage, "50")
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[]`)
	})

	ctx := context.Background()
	opts := &ChangelogOptions{ID: "last123", PerPage: 50}
	_, _, err := client.Issues.GetChangelog(ctx, "QUEUE-1", opts)
	if err != nil {
		t.Fatalf("GetChangelog returned error: %v", err)
	}
}

func TestIssuesService_GetChangelog_Empty(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v2/issues/{key}/changelog", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[]`)
	})

	ctx := context.Background()
	changelog, _, err := client.Issues.GetChangelog(ctx, "QUEUE-1", nil)
	if err != nil {
		t.Fatalf("GetChangelog returned error: %v", err)
	}
	if len(changelog) != 0 {
		t.Errorf("GetChangelog returned %d entries, want 0", len(changelog))
	}
}
