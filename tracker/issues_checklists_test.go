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

	mux.HandleFunc("GET /v3/issues/{key}/checklistItems", func(w http.ResponseWriter, r *http.Request) {
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

	mux.HandleFunc("POST /v3/issues/{key}/checklistItems", func(w http.ResponseWriter, r *http.Request) {
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
			"self": "https://api.tracker.yandex.net/v3/issues/QUEUE-1",
			"id": "100",
			"key": "QUEUE-1",
			"summary": "Test issue",
			"checklistItems": [
				{
					"id": "item-1",
					"text": "New item",
					"textHtml": "<b>New item</b>",
					"checked": false,
					"assignee": {
						"self": "https://api.tracker.yandex.net/v3/users/1",
						"id": "1",
						"display": "John Doe"
					},
					"deadline": {
						"date": "2025-01-15T00:00:00.000+0000",
						"deadlineType": "date",
						"isExceeded": false
					},
					"checklistItemType": "standard"
				}
			],
			"checklistDone": 0,
			"checklistTotal": 1
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
	if len(issue.ChecklistItems) != 1 {
		t.Fatalf("Issue.ChecklistItems length = %d, want 1", len(issue.ChecklistItems))
	}
	if got := *issue.ChecklistItems[0].ID; got != "item-1" {
		t.Errorf("ChecklistItems[0].ID = %q, want %q", got, "item-1")
	}
	if got := *issue.ChecklistItems[0].Text; got != "New item" {
		t.Errorf("ChecklistItems[0].Text = %q, want %q", got, "New item")
	}
	if got := *issue.ChecklistItems[0].Checked; got != false {
		t.Errorf("ChecklistItems[0].Checked = %v, want false", got)
	}
	if got := *issue.ChecklistItems[0].Assignee.Display; got != "John Doe" {
		t.Errorf("ChecklistItems[0].Assignee.Display = %q, want %q", got, "John Doe")
	}
	if got := *issue.ChecklistItems[0].Deadline.IsExceeded; got != false {
		t.Errorf("ChecklistItems[0].Deadline.IsExceeded = %v, want false", got)
	}
	if got := *issue.ChecklistDone; got != 0 {
		t.Errorf("Issue.ChecklistDone = %d, want 0", got)
	}
	if got := *issue.ChecklistTotal; got != 1 {
		t.Errorf("Issue.ChecklistTotal = %d, want 1", got)
	}
}

func TestIssuesService_EditChecklistItem(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("PATCH /v3/issues/{key}/checklistItems/{itemID}", func(w http.ResponseWriter, r *http.Request) {
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
			"self": "https://api.tracker.yandex.net/v3/issues/QUEUE-1",
			"id": "100",
			"key": "QUEUE-1",
			"summary": "Test issue",
			"checklistItems": [
				{
					"id": "item-1",
					"text": "Updated item",
					"checked": true,
					"checklistItemType": "standard"
				}
			],
			"checklistDone": 1,
			"checklistTotal": 1
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
	if len(issue.ChecklistItems) != 1 {
		t.Fatalf("Issue.ChecklistItems length = %d, want 1", len(issue.ChecklistItems))
	}
	if got := *issue.ChecklistItems[0].Checked; got != true {
		t.Errorf("ChecklistItems[0].Checked = %v, want true", got)
	}
	if got := *issue.ChecklistDone; got != 1 {
		t.Errorf("Issue.ChecklistDone = %d, want 1", got)
	}
	if got := *issue.ChecklistTotal; got != 1 {
		t.Errorf("Issue.ChecklistTotal = %d, want 1", got)
	}
}

func TestIssuesService_DeleteChecklistItem(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("DELETE /v3/issues/{key}/checklistItems/{itemID}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/issues/QUEUE-1",
			"id": "100",
			"key": "QUEUE-1",
			"summary": "Test issue",
			"checklistItems": [
				{
					"id": "item-2",
					"text": "Remaining item",
					"checked": false,
					"checklistItemType": "standard"
				}
			],
			"checklistDone": 0,
			"checklistTotal": 1
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
	if len(issue.ChecklistItems) != 1 {
		t.Fatalf("Issue.ChecklistItems length = %d, want 1", len(issue.ChecklistItems))
	}
	if got := *issue.ChecklistItems[0].ID; got != "item-2" {
		t.Errorf("ChecklistItems[0].ID = %q, want %q", got, "item-2")
	}
	if got := *issue.ChecklistTotal; got != 1 {
		t.Errorf("Issue.ChecklistTotal = %d, want 1", got)
	}
}
