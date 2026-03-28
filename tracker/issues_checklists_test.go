package tracker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestIssuesService_ListChecklistItems(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v2/issues/{key}/checklistItems", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{
			"id": "1",
			"text": "Buy groceries",
			"checked": false,
			"checklistItemType": "standard"
		}, {
			"id": "2",
			"text": "Clean house",
			"checked": true,
			"checklistItemType": "standard"
		}]`)
	})

	ctx := context.Background()
	items, _, err := client.Issues.ListChecklistItems(ctx, "QUEUE-1")
	if err != nil {
		t.Fatalf("ListChecklistItems returned error: %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("ListChecklistItems returned %d items, want 2", len(items))
	}

	want := &ChecklistItem{
		ID:                Ptr("1"),
		Text:              Ptr("Buy groceries"),
		Checked:           Ptr(false),
		ChecklistItemType: Ptr("standard"),
	}
	if !reflect.DeepEqual(items[0], want) {
		t.Errorf("ListChecklistItems[0] = %+v, want %+v", items[0], want)
	}
}

func TestIssuesService_CreateChecklistItem(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v2/issues/{key}/checklistItems", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		var body ChecklistItemRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if got := *body.Text; got != "New item" {
			t.Errorf("Request body text = %q, want %q", got, "New item")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v2/issues/QUEUE-1",
			"id": "100",
			"key": "QUEUE-1",
			"summary": "Test issue"
		}`)
	})

	ctx := context.Background()
	issue, _, err := client.Issues.CreateChecklistItem(ctx, "QUEUE-1", &ChecklistItemRequest{
		Text: Ptr("New item"),
	})
	if err != nil {
		t.Fatalf("CreateChecklistItem returned error: %v", err)
	}

	if got := *issue.Key; got != "QUEUE-1" {
		t.Errorf("Issue.Key = %q, want %q", got, "QUEUE-1")
	}
	if got := *issue.Summary; got != "Test issue" {
		t.Errorf("Issue.Summary = %q, want %q", got, "Test issue")
	}
}

func TestIssuesService_EditChecklistItem(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("PATCH /v2/issues/{key}/checklistItems/{itemID}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")

		var body ChecklistItemRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if got := *body.Text; got != "Updated item" {
			t.Errorf("Request body text = %q, want %q", got, "Updated item")
		}
		if got := *body.Checked; got != true {
			t.Errorf("Request body checked = %v, want %v", got, true)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v2/issues/QUEUE-1",
			"id": "100",
			"key": "QUEUE-1",
			"summary": "Test issue"
		}`)
	})

	ctx := context.Background()
	issue, _, err := client.Issues.EditChecklistItem(ctx, "QUEUE-1", "item-1", &ChecklistItemRequest{
		Text:    Ptr("Updated item"),
		Checked: Ptr(true),
	})
	if err != nil {
		t.Fatalf("EditChecklistItem returned error: %v", err)
	}

	if got := *issue.Key; got != "QUEUE-1" {
		t.Errorf("Issue.Key = %q, want %q", got, "QUEUE-1")
	}
}

func TestIssuesService_DeleteChecklistItem(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("DELETE /v2/issues/{key}/checklistItems/{itemID}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v2/issues/QUEUE-1",
			"id": "100",
			"key": "QUEUE-1",
			"summary": "Test issue"
		}`)
	})

	ctx := context.Background()
	issue, _, err := client.Issues.DeleteChecklistItem(ctx, "QUEUE-1", "item-1")
	if err != nil {
		t.Fatalf("DeleteChecklistItem returned error: %v", err)
	}

	if issue == nil {
		t.Fatal("DeleteChecklistItem returned nil issue, want non-nil")
	}
	if got := *issue.Key; got != "QUEUE-1" {
		t.Errorf("Issue.Key = %q, want %q", got, "QUEUE-1")
	}
}
