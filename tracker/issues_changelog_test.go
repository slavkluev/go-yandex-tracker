package tracker

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestIssuesService_GetChangelog(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/issues/{key}/changelog", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{
			"id": "1",
			"type": "IssueUpdated",
			"transport": "front",
			"updatedAt": "2017-06-11T05:16:01.339+0000",
			"updatedBy": {"id": "user1"},
			"fields": [{
				"field": {
					"self": "https://api.tracker.yandex.net/v3/fields/status",
					"id": "status",
					"display": "Status"
				},
				"from": "open",
				"to": "closed"
			}],
			"comments": {
				"added": [{
					"self": "https://api.tracker.yandex.net/v3/issues/QUEUE-1/comments/10",
					"id": "10",
					"display": "Comment text"
				}]
			},
			"executedTriggers": [{
				"trigger": {
					"self": "https://api.tracker.yandex.net/v3/queues/QUEUE/triggers/29",
					"id": "29",
					"display": "Trigger-42"
				},
				"success": true,
				"message": "Success"
			}]
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
	if got := *entry.Transport; got != "front" {
		t.Errorf("Transport = %q, want %q", got, "front")
	}

	if len(entry.Fields) != 1 {
		t.Fatalf("Fields length = %d, want 1", len(entry.Fields))
	}
	if got := string(*entry.Fields[0].Field.ID); got != "status" {
		t.Errorf("Fields[0].Field.ID = %q, want %q", got, "status")
	}
	if got := *entry.Fields[0].Field.Display; got != "Status" {
		t.Errorf("Fields[0].Field.Display = %q, want %q", got, "Status")
	}

	if entry.Comments == nil {
		t.Fatal("Comments is nil")
	}
	if len(entry.Comments.Added) != 1 {
		t.Fatalf("Comments.Added length = %d, want 1", len(entry.Comments.Added))
	}
	if got := string(*entry.Comments.Added[0].ID); got != "10" {
		t.Errorf("Comments.Added[0].ID = %q, want %q", got, "10")
	}
	if got := *entry.Comments.Added[0].Display; got != "Comment text" {
		t.Errorf("Comments.Added[0].Display = %q, want %q", got, "Comment text")
	}

	if len(entry.ExecutedTriggers) != 1 {
		t.Fatalf("ExecutedTriggers length = %d, want 1", len(entry.ExecutedTriggers))
	}
	if got := string(*entry.ExecutedTriggers[0].Trigger.ID); got != "29" {
		t.Errorf("ExecutedTriggers[0].Trigger.ID = %q, want %q", got, "29")
	}
	if got := *entry.ExecutedTriggers[0].Success; got != true {
		t.Errorf("ExecutedTriggers[0].Success = %v, want true", got)
	}
	if got := *entry.ExecutedTriggers[0].Message; got != "Success" {
		t.Errorf("ExecutedTriggers[0].Message = %q, want %q", got, "Success")
	}
}

func TestIssuesService_GetChangelog_WithOptions(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/issues/{key}/changelog", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		gotID := r.URL.Query().Get("id")
		if gotID != "last123" {
			t.Errorf("query param id = %q, want %q", gotID, "last123")
		}
		gotPerPage := r.URL.Query().Get("perPage")
		if gotPerPage != "50" {
			t.Errorf("query param perPage = %q, want %q", gotPerPage, "50")
		}
		gotField := r.URL.Query().Get("field")
		if gotField != "status" {
			t.Errorf("query param field = %q, want %q", gotField, "status")
		}
		gotType := r.URL.Query().Get("type")
		if gotType != "IssueUpdated" {
			t.Errorf("query param type = %q, want %q", gotType, "IssueUpdated")
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[]`)
	})

	ctx := context.Background()
	opts := &ChangelogOptions{ID: "last123", PerPage: 50, Field: "status", Type: "IssueUpdated"}
	_, _, err := client.Issues.GetChangelog(ctx, "QUEUE-1", opts)
	if err != nil {
		t.Fatalf("GetChangelog returned error: %v", err)
	}
}

func TestIssuesService_GetChangelog_Empty(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/issues/{key}/changelog", func(w http.ResponseWriter, r *http.Request) {
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
