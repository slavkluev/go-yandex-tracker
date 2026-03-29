package tracker

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestSprintsService_List(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/boards/{boardID}/sprints", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[
			{
				"self": "https://api.tracker.yandex.net/v3/sprints/1",
				"id": 1,
				"version": 1,
				"name": "Sprint 1",
				"board": {"self": "https://api.tracker.yandex.net/v3/boards/1", "id": "1", "display": "Test Board"},
				"status": "in_progress",
				"archived": false,
				"startDate": "2024-01-01",
				"endDate": "2024-01-14"
			},
			{
				"self": "https://api.tracker.yandex.net/v3/sprints/2",
				"id": 2,
				"version": 1,
				"name": "Sprint 2",
				"board": {"self": "https://api.tracker.yandex.net/v3/boards/1", "id": "1", "display": "Test Board"},
				"status": "draft",
				"archived": false,
				"startDate": "2024-01-15",
				"endDate": "2024-01-28"
			}
		]`)
	})

	ctx := context.Background()
	sprints, _, err := client.Sprints.List(ctx, 1)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	if len(sprints) != 2 {
		t.Fatalf("List returned %d sprints, want 2", len(sprints))
	}

	want := []*Sprint{
		{
			Self:      Ptr("https://api.tracker.yandex.net/v3/sprints/1"),
			ID:        Ptr(1),
			Version:   Ptr(1),
			Name:      Ptr("Sprint 1"),
			Board:     &BoardRef{Self: Ptr("https://api.tracker.yandex.net/v3/boards/1"), ID: Ptr("1"), Display: Ptr("Test Board")},
			Status:    Ptr("in_progress"),
			Archived:  Ptr(false),
			StartDate: Ptr("2024-01-01"),
			EndDate:   Ptr("2024-01-14"),
		},
		{
			Self:      Ptr("https://api.tracker.yandex.net/v3/sprints/2"),
			ID:        Ptr(2),
			Version:   Ptr(1),
			Name:      Ptr("Sprint 2"),
			Board:     &BoardRef{Self: Ptr("https://api.tracker.yandex.net/v3/boards/1"), ID: Ptr("1"), Display: Ptr("Test Board")},
			Status:    Ptr("draft"),
			Archived:  Ptr(false),
			StartDate: Ptr("2024-01-15"),
			EndDate:   Ptr("2024-01-28"),
		},
	}

	if !reflect.DeepEqual(sprints, want) {
		t.Errorf("List returned %+v, want %+v", sprints, want)
	}
}

func TestSprintsService_Create(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/sprints", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"name":"Sprint 3","board":{"id":"1"},"startDate":"2024-02-01","endDate":"2024-02-14"}`)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/sprints/3",
			"id": 3,
			"version": 1,
			"name": "Sprint 3",
			"board": {"self": "https://api.tracker.yandex.net/v3/boards/1", "id": "1", "display": "Test Board"},
			"status": "draft",
			"archived": false,
			"startDate": "2024-02-01",
			"endDate": "2024-02-14"
		}`)
	})

	ctx := context.Background()
	sprint, resp, err := client.Sprints.Create(ctx, &SprintCreateRequest{
		Name:      Ptr("Sprint 3"),
		Board:     &BoardRef{ID: Ptr("1")},
		StartDate: Ptr("2024-02-01"),
		EndDate:   Ptr("2024-02-14"),
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusCreated)
	}
	if got := *sprint.ID; got != 3 {
		t.Errorf("ID = %d, want %d", got, 3)
	}
	if got := *sprint.Name; got != "Sprint 3" {
		t.Errorf("Name = %q, want %q", got, "Sprint 3")
	}
	if got := *sprint.StartDate; got != "2024-02-01" {
		t.Errorf("StartDate = %q, want %q", got, "2024-02-01")
	}
}

func TestSprintsService_Get(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/sprints/{sprintID}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/sprints/1",
			"id": 1,
			"version": 1,
			"name": "Sprint 1",
			"board": {"self": "https://api.tracker.yandex.net/v3/boards/1", "id": "1", "display": "Test Board"},
			"status": "in_progress",
			"archived": false,
			"createdBy": {"self": "https://api.tracker.yandex.net/v3/users/1", "id": "1", "display": "John"},
			"startDate": "2024-01-01",
			"endDate": "2024-01-14"
		}`)
	})

	ctx := context.Background()
	sprint, _, err := client.Sprints.Get(ctx, 1)
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if got := *sprint.ID; got != 1 {
		t.Errorf("ID = %d, want %d", got, 1)
	}
	if got := *sprint.Name; got != "Sprint 1" {
		t.Errorf("Name = %q, want %q", got, "Sprint 1")
	}
	if got := *sprint.Board.ID; got != "1" {
		t.Errorf("Board.ID = %q, want %q", got, "1")
	}
	if got := *sprint.Status; got != "in_progress" {
		t.Errorf("Status = %q, want %q", got, "in_progress")
	}
	if got := *sprint.StartDate; got != "2024-01-01" {
		t.Errorf("StartDate = %q, want %q", got, "2024-01-01")
	}
}

func TestSprintsService_Get_NotFound(t *testing.T) {
	client, mux := setup(t)
	ctx := context.Background()

	mux.HandleFunc("GET /v3/sprints/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"errorMessages":["Object not found"],"errors":{}}`)
	})

	_, _, err := client.Sprints.Get(ctx, 99999)
	if err == nil {
		t.Fatal("Sprints.Get expected error for 404, got nil")
	}
	if !IsNotFound(err) {
		t.Errorf("IsNotFound(err) = false, want true")
	}
}
