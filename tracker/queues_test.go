package tracker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestQueuesService_Create(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/queues/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "application/json")

		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["key"] != "TEST" {
			t.Errorf("Request body key = %v, want %v", body["key"], "TEST")
		}
		if body["name"] != "Test Queue" {
			t.Errorf("Request body name = %v, want %v", body["name"], "Test Queue")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Queue{
			Self:       Ptr("https://api.tracker.yandex.net/v3/queues/TEST"),
			ID:         Ptr("1"),
			Key:        Ptr("TEST"),
			Display:    Ptr("Test Queue"),
			Version:    Ptr(1),
			Name:       Ptr("Test Queue"),
			Lead:       &User{Self: Ptr("https://api.tracker.yandex.net/v3/users/123"), ID: Ptr("123"), Display: Ptr("John Doe")},
			AssignAuto: Ptr(false),
			DefaultType: &IssueType{
				Self:    Ptr("https://api.tracker.yandex.net/v3/issuetypes/2"),
				ID:      Ptr("2"),
				Key:     Ptr("task"),
				Display: Ptr("Task"),
			},
			DefaultPriority: &Priority{
				Self:    Ptr("https://api.tracker.yandex.net/v3/priorities/3"),
				ID:      Ptr("3"),
				Key:     Ptr("normal"),
				Display: Ptr("Normal"),
			},
		})
	})

	queue, _, err := client.Queues.Create(ctx, &QueueCreateRequest{
		Key:  Ptr("TEST"),
		Name: Ptr("Test Queue"),
		Lead: Ptr("john.doe"),
	})
	if err != nil {
		t.Fatalf("Queues.Create returned error: %v", err)
	}

	want := &Queue{
		Self:       Ptr("https://api.tracker.yandex.net/v3/queues/TEST"),
		ID:         Ptr("1"),
		Key:        Ptr("TEST"),
		Display:    Ptr("Test Queue"),
		Version:    Ptr(1),
		Name:       Ptr("Test Queue"),
		Lead:       &User{Self: Ptr("https://api.tracker.yandex.net/v3/users/123"), ID: Ptr("123"), Display: Ptr("John Doe")},
		AssignAuto: Ptr(false),
		DefaultType: &IssueType{
			Self:    Ptr("https://api.tracker.yandex.net/v3/issuetypes/2"),
			ID:      Ptr("2"),
			Key:     Ptr("task"),
			Display: Ptr("Task"),
		},
		DefaultPriority: &Priority{
			Self:    Ptr("https://api.tracker.yandex.net/v3/priorities/3"),
			ID:      Ptr("3"),
			Key:     Ptr("normal"),
			Display: Ptr("Normal"),
		},
	}

	if !reflect.DeepEqual(queue, want) {
		t.Errorf("Queues.Create returned %+v, want %+v", queue, want)
	}
}

func TestQueuesService_Get(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/queues/{key}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		if got := r.URL.Query().Get("expand"); got != "all" {
			t.Errorf("expand = %q, want %q", got, "all")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Queue{
			Self:    Ptr("https://api.tracker.yandex.net/v3/queues/TEST"),
			ID:      Ptr("1"),
			Key:     Ptr("TEST"),
			Display: Ptr("Test Queue"),
			Version: Ptr(1),
			Name:    Ptr("Test Queue"),
			Lead:    &User{Self: Ptr("https://api.tracker.yandex.net/v3/users/123"), ID: Ptr("123"), Display: Ptr("John Doe")},
		})
	})

	queue, _, err := client.Queues.Get(ctx, "TEST", &QueueGetOptions{Expand: "all"})
	if err != nil {
		t.Fatalf("Queues.Get returned error: %v", err)
	}

	want := &Queue{
		Self:    Ptr("https://api.tracker.yandex.net/v3/queues/TEST"),
		ID:      Ptr("1"),
		Key:     Ptr("TEST"),
		Display: Ptr("Test Queue"),
		Version: Ptr(1),
		Name:    Ptr("Test Queue"),
		Lead:    &User{Self: Ptr("https://api.tracker.yandex.net/v3/users/123"), ID: Ptr("123"), Display: Ptr("John Doe")},
	}

	if !reflect.DeepEqual(queue, want) {
		t.Errorf("Queues.Get returned %+v, want %+v", queue, want)
	}
}

func TestQueuesService_List(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/queues", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		if got := r.URL.Query().Get("page"); got != "1" {
			t.Errorf("page = %q, want %q", got, "1")
		}
		if got := r.URL.Query().Get("perPage"); got != "10" {
			t.Errorf("perPage = %q, want %q", got, "10")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]*Queue{
			{
				Self:    Ptr("https://api.tracker.yandex.net/v3/queues/TEST"),
				ID:      Ptr("1"),
				Key:     Ptr("TEST"),
				Display: Ptr("Test Queue"),
				Name:    Ptr("Test Queue"),
			},
			{
				Self:    Ptr("https://api.tracker.yandex.net/v3/queues/DEV"),
				ID:      Ptr("2"),
				Key:     Ptr("DEV"),
				Display: Ptr("Dev Queue"),
				Name:    Ptr("Dev Queue"),
			},
		})
	})

	queues, _, err := client.Queues.List(ctx, &QueueListOptions{
		ListOptions: ListOptions{Page: 1, PerPage: 10},
	})
	if err != nil {
		t.Fatalf("Queues.List returned error: %v", err)
	}

	if len(queues) != 2 {
		t.Fatalf("Queues.List returned %d queues, want 2", len(queues))
	}

	want := []*Queue{
		{
			Self:    Ptr("https://api.tracker.yandex.net/v3/queues/TEST"),
			ID:      Ptr("1"),
			Key:     Ptr("TEST"),
			Display: Ptr("Test Queue"),
			Name:    Ptr("Test Queue"),
		},
		{
			Self:    Ptr("https://api.tracker.yandex.net/v3/queues/DEV"),
			ID:      Ptr("2"),
			Key:     Ptr("DEV"),
			Display: Ptr("Dev Queue"),
			Name:    Ptr("Dev Queue"),
		},
	}

	if !reflect.DeepEqual(queues, want) {
		t.Errorf("Queues.List returned %+v, want %+v", queues, want)
	}
}

func TestQueuesService_Delete(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("DELETE /v3/queues/{key}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.Queues.Delete(ctx, "TEST")
	if err != nil {
		t.Fatalf("Queues.Delete returned error: %v", err)
	}
	if resp == nil {
		t.Fatal("Queues.Delete returned nil response")
	}
}

func TestQueuesService_Restore(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/queues/{key}/_restore", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Queue{
			Self:    Ptr("https://api.tracker.yandex.net/v3/queues/TEST"),
			ID:      Ptr("1"),
			Key:     Ptr("TEST"),
			Display: Ptr("Test Queue"),
			Version: Ptr(2),
			Name:    Ptr("Test Queue"),
		})
	})

	queue, _, err := client.Queues.Restore(ctx, "TEST")
	if err != nil {
		t.Fatalf("Queues.Restore returned error: %v", err)
	}

	want := &Queue{
		Self:    Ptr("https://api.tracker.yandex.net/v3/queues/TEST"),
		ID:      Ptr("1"),
		Key:     Ptr("TEST"),
		Display: Ptr("Test Queue"),
		Version: Ptr(2),
		Name:    Ptr("Test Queue"),
	}

	if !reflect.DeepEqual(queue, want) {
		t.Errorf("Queues.Restore returned %+v, want %+v", queue, want)
	}
}

func TestQueuesService_Get_NotFound(t *testing.T) {
	client, mux := setup(t)
	ctx := context.Background()

	mux.HandleFunc("GET /v3/queues/{key}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"errorMessages":["Object not found"],"errors":{}}`)
	})

	_, _, err := client.Queues.Get(ctx, "NONEXIST", nil)
	if err == nil {
		t.Fatal("Queues.Get expected error for 404, got nil")
	}
	if !IsNotFound(err) {
		t.Errorf("IsNotFound(err) = false, want true")
	}
}
