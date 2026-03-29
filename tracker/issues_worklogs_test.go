package tracker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestIssuesService_ListWorklogs(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/issues/{key}/worklog", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{
			"self": "https://api.tracker.yandex.net/v3/issues/QUEUE-1/worklog/1",
			"id": "1",
			"comment": "Working on task",
			"start": "2021-09-21T15:22:00.000+0000",
			"duration": "P3W",
			"createdBy": {"self": "https://api.tracker.yandex.net/v3/users/1", "id": "1", "display": "John Doe"}
		}]`)
	})

	ctx := context.Background()
	worklogs, _, err := client.Issues.ListWorklogs(ctx, "QUEUE-1")
	if err != nil {
		t.Fatalf("ListWorklogs returned error: %v", err)
	}

	if len(worklogs) != 1 {
		t.Fatalf("ListWorklogs returned %d worklogs, want 1", len(worklogs))
	}

	wl := worklogs[0]
	if got := *wl.ID; got != "1" {
		t.Errorf("ID = %q, want %q", got, "1")
	}
	if got := *wl.Comment; got != "Working on task" {
		t.Errorf("Comment = %q, want %q", got, "Working on task")
	}

	// P3W = 3 weeks = 504 hours
	wantDuration := 3 * 7 * 24 * time.Hour
	if wl.Duration == nil {
		t.Fatal("Duration is nil, want non-nil")
	}
	if got := wl.Duration.Duration; got != wantDuration {
		t.Errorf("Duration = %v, want %v", got, wantDuration)
	}

	if wl.Start == nil {
		t.Fatal("Start is nil, want non-nil")
	}
	if wl.Start.IsZero() {
		t.Error("Start is zero, want non-zero timestamp")
	}
}

func TestIssuesService_CreateWorklog(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/issues/{key}/worklog", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if _, ok := body["start"]; !ok {
			t.Error("Request body missing 'start' field")
		}
		if _, ok := body["duration"]; !ok {
			t.Error("Request body missing 'duration' field")
		}
		if _, ok := body["comment"]; !ok {
			t.Error("Request body missing 'comment' field")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/issues/QUEUE-1/worklog/1",
			"id": "1",
			"comment": "Coding session",
			"start": "2021-09-21T15:00:00.000+0000",
			"duration": "PT1H30M",
			"createdBy": {"self": "https://api.tracker.yandex.net/v3/users/1", "id": "1", "display": "John Doe"}
		}`)
	})

	ctx := context.Background()
	startTime := Timestamp{Time: time.Date(2021, 9, 21, 15, 0, 0, 0, time.UTC)}
	worklog, _, err := client.Issues.CreateWorklog(ctx, "QUEUE-1", &WorklogRequest{
		Start:    &startTime,
		Duration: &Duration{Duration: 1*time.Hour + 30*time.Minute},
		Comment:  Ptr("Coding session"),
	})
	if err != nil {
		t.Fatalf("CreateWorklog returned error: %v", err)
	}

	if got := *worklog.ID; got != "1" {
		t.Errorf("ID = %q, want %q", got, "1")
	}
	if got := *worklog.Comment; got != "Coding session" {
		t.Errorf("Comment = %q, want %q", got, "Coding session")
	}

	// PT1H30M = 1 hour 30 minutes
	wantDuration := 1*time.Hour + 30*time.Minute
	if worklog.Duration == nil {
		t.Fatal("Duration is nil, want non-nil")
	}
	if got := worklog.Duration.Duration; got != wantDuration {
		t.Errorf("Duration = %v, want %v", got, wantDuration)
	}
}

func TestIssuesService_EditWorklog(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("PATCH /v3/issues/{key}/worklog/{worklogID}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if _, ok := body["duration"]; !ok {
			t.Error("Request body missing 'duration' field")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/issues/QUEUE-1/worklog/abc123",
			"id": "abc123",
			"comment": "Updated session",
			"start": "2021-09-21T15:00:00.000+0000",
			"duration": "PT2H",
			"createdBy": {"self": "https://api.tracker.yandex.net/v3/users/1", "id": "1", "display": "John Doe"}
		}`)
	})

	ctx := context.Background()
	worklog, _, err := client.Issues.EditWorklog(ctx, "QUEUE-1", "abc123", &WorklogRequest{
		Duration: &Duration{Duration: 2 * time.Hour},
	})
	if err != nil {
		t.Fatalf("EditWorklog returned error: %v", err)
	}

	if got := *worklog.ID; got != "abc123" {
		t.Errorf("ID = %q, want %q", got, "abc123")
	}

	wantDuration := 2 * time.Hour
	if worklog.Duration == nil {
		t.Fatal("Duration is nil, want non-nil")
	}
	if got := worklog.Duration.Duration; got != wantDuration {
		t.Errorf("Duration = %v, want %v", got, wantDuration)
	}
}

func TestIssuesService_DeleteWorklog(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("DELETE /v3/issues/{key}/worklog/{worklogID}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	resp, err := client.Issues.DeleteWorklog(ctx, "QUEUE-1", "abc123")
	if err != nil {
		t.Fatalf("DeleteWorklog returned error: %v", err)
	}
	if got := resp.StatusCode; got != http.StatusNoContent {
		t.Errorf("StatusCode = %d, want %d", got, http.StatusNoContent)
	}
}
