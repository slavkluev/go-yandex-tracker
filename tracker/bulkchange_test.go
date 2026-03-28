package tracker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestBulkChangeService_Move(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/bulkchange/_move", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		// Verify request body contains expected fields.
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if body["queue"] != "NEWQUEUE" {
			t.Errorf("body queue = %v, want %q", body["queue"], "NEWQUEUE")
		}
		issues, ok := body["issues"].([]any)
		if !ok || len(issues) != 2 {
			t.Errorf("body issues = %v, want 2 elements", body["issues"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/bulkchange/593cd211ef7e8a0000000001",
			"id": "593cd211ef7e8a0000000001",
			"status": "CREATED",
			"statusText": "Bulk operation created"
		}`)
	})

	ctx := context.Background()
	bc, _, err := client.BulkChange.Move(ctx, &BulkMoveRequest{
		Queue:         Ptr("NEWQUEUE"),
		Issues:        []string{"QUEUE-1", "QUEUE-2"},
		MoveAllFields: Ptr(true),
	})
	if err != nil {
		t.Fatalf("Move returned error: %v", err)
	}

	want := &BulkChange{
		Self:       Ptr("https://api.tracker.yandex.net/v3/bulkchange/593cd211ef7e8a0000000001"),
		ID:         Ptr("593cd211ef7e8a0000000001"),
		Status:     Ptr("CREATED"),
		StatusText: Ptr("Bulk operation created"),
	}
	if !reflect.DeepEqual(bc, want) {
		t.Errorf("Move returned %+v, want %+v", bc, want)
	}
}

func TestBulkChangeService_Update(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/bulkchange/_update", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		// Verify request body contains expected fields.
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		issues, ok := body["issues"].([]any)
		if !ok || len(issues) != 1 {
			t.Errorf("body issues = %v, want 1 element", body["issues"])
		}
		values, ok := body["values"].(map[string]any)
		if !ok || values["priority"] != "critical" {
			t.Errorf("body values = %v, want priority=critical", body["values"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/bulkchange/593cd211ef7e8a0000000002",
			"id": "593cd211ef7e8a0000000002",
			"status": "CREATED",
			"statusText": "Bulk operation created"
		}`)
	})

	ctx := context.Background()
	bc, _, err := client.BulkChange.Update(ctx, &BulkUpdateRequest{
		Issues: []string{"QUEUE-1"},
		Values: map[string]any{"priority": "critical"},
	})
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	want := &BulkChange{
		Self:       Ptr("https://api.tracker.yandex.net/v3/bulkchange/593cd211ef7e8a0000000002"),
		ID:         Ptr("593cd211ef7e8a0000000002"),
		Status:     Ptr("CREATED"),
		StatusText: Ptr("Bulk operation created"),
	}
	if !reflect.DeepEqual(bc, want) {
		t.Errorf("Update returned %+v, want %+v", bc, want)
	}
}

func TestBulkChangeService_Transition(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/bulkchange/_transition", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		// Verify request body contains expected fields.
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if body["transition"] != "close" {
			t.Errorf("body transition = %v, want %q", body["transition"], "close")
		}
		issues, ok := body["issues"].([]any)
		if !ok || len(issues) != 2 {
			t.Errorf("body issues = %v, want 2 elements", body["issues"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/bulkchange/593cd211ef7e8a0000000003",
			"id": "593cd211ef7e8a0000000003",
			"status": "CREATED",
			"statusText": "Bulk operation created"
		}`)
	})

	ctx := context.Background()
	bc, _, err := client.BulkChange.Transition(ctx, &BulkTransitionRequest{
		Transition: Ptr("close"),
		Issues:     []string{"QUEUE-1", "QUEUE-2"},
	})
	if err != nil {
		t.Fatalf("Transition returned error: %v", err)
	}

	want := &BulkChange{
		Self:       Ptr("https://api.tracker.yandex.net/v3/bulkchange/593cd211ef7e8a0000000003"),
		ID:         Ptr("593cd211ef7e8a0000000003"),
		Status:     Ptr("CREATED"),
		StatusText: Ptr("Bulk operation created"),
	}
	if !reflect.DeepEqual(bc, want) {
		t.Errorf("Transition returned %+v, want %+v", bc, want)
	}
}

func TestBulkChangeService_GetStatus(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/bulkchange/{id}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/bulkchange/593cd211ef7e8a0000000001",
			"id": "593cd211ef7e8a0000000001",
			"status": "COMPLETED",
			"statusText": "Bulk edit completed",
			"executionChunkPercent": 100,
			"executionIssuePercent": 100,
			"totalIssues": 2,
			"totalCompletedIssues": 2
		}`)
	})

	ctx := context.Background()
	bc, _, err := client.BulkChange.GetStatus(ctx, "593cd211ef7e8a0000000001")
	if err != nil {
		t.Fatalf("GetStatus returned error: %v", err)
	}

	want := &BulkChange{
		Self:                  Ptr("https://api.tracker.yandex.net/v3/bulkchange/593cd211ef7e8a0000000001"),
		ID:                    Ptr("593cd211ef7e8a0000000001"),
		Status:                Ptr("COMPLETED"),
		StatusText:            Ptr("Bulk edit completed"),
		ExecutionChunkPercent: Ptr(100),
		ExecutionIssuePercent: Ptr(100),
		TotalIssues:           Ptr(2),
		TotalCompletedIssues:  Ptr(2),
	}
	if !reflect.DeepEqual(bc, want) {
		t.Errorf("GetStatus returned %+v, want %+v", bc, want)
	}
}

func TestBulkChangeService_GetStatus_NotFound(t *testing.T) {
	client, mux := setup(t)
	ctx := context.Background()

	mux.HandleFunc("GET /v3/bulkchange/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"errorMessages":["Object not found"],"errors":{}}`)
	})

	_, _, err := client.BulkChange.GetStatus(ctx, "nonexistent123")
	if err == nil {
		t.Fatal("BulkChange.GetStatus expected error for 404, got nil")
	}
	if !IsNotFound(err) {
		t.Errorf("IsNotFound(err) = false, want true")
	}
}
