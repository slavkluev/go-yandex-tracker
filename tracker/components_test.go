package tracker

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestComponentsService_List(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/components", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[
			{
				"self": "https://api.tracker.yandex.net/v3/components/1",
				"id": 1,
				"version": 1,
				"name": "Backend",
				"queue": {"self": "https://api.tracker.yandex.net/v3/queues/TEST", "id": "1", "key": "TEST", "display": "Test Queue"},
				"description": "Backend component",
				"lead": {"self": "https://api.tracker.yandex.net/v3/users/1", "id": "1", "display": "John"},
				"assignAuto": false
			},
			{
				"self": "https://api.tracker.yandex.net/v3/components/2",
				"id": 2,
				"version": 1,
				"name": "Frontend",
				"queue": {"self": "https://api.tracker.yandex.net/v3/queues/TEST", "id": "1", "key": "TEST", "display": "Test Queue"},
				"description": "Frontend component",
				"lead": {"self": "https://api.tracker.yandex.net/v3/users/2", "id": "2", "display": "Jane"},
				"assignAuto": true
			}
		]`)
	})

	ctx := context.Background()
	components, _, err := client.Components.List(ctx)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	if len(components) != 2 {
		t.Fatalf("List returned %d components, want 2", len(components))
	}

	want := []*Component{
		{
			Self:        Ptr("https://api.tracker.yandex.net/v3/components/1"),
			ID:          Ptr(1),
			Version:     Ptr(1),
			Name:        Ptr("Backend"),
			Queue:       &Queue{Self: Ptr("https://api.tracker.yandex.net/v3/queues/TEST"), ID: Ptr("1"), Key: Ptr("TEST"), Display: Ptr("Test Queue")},
			Description: Ptr("Backend component"),
			Lead:        &User{Self: Ptr("https://api.tracker.yandex.net/v3/users/1"), ID: Ptr("1"), Display: Ptr("John")},
			AssignAuto:  Ptr(false),
		},
		{
			Self:        Ptr("https://api.tracker.yandex.net/v3/components/2"),
			ID:          Ptr(2),
			Version:     Ptr(1),
			Name:        Ptr("Frontend"),
			Queue:       &Queue{Self: Ptr("https://api.tracker.yandex.net/v3/queues/TEST"), ID: Ptr("1"), Key: Ptr("TEST"), Display: Ptr("Test Queue")},
			Description: Ptr("Frontend component"),
			Lead:        &User{Self: Ptr("https://api.tracker.yandex.net/v3/users/2"), ID: Ptr("2"), Display: Ptr("Jane")},
			AssignAuto:  Ptr(true),
		},
	}

	if !reflect.DeepEqual(components, want) {
		t.Errorf("List returned %+v, want %+v", components, want)
	}
}

func TestComponentsService_Create(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/components", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"name":"Backend","queue":"TEST","description":"Backend component","lead":"user1","assignAuto":false}`)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/components/1",
			"id": 1,
			"version": 1,
			"name": "Backend",
			"queue": {"self": "https://api.tracker.yandex.net/v3/queues/TEST", "id": "1", "key": "TEST", "display": "Test Queue"},
			"description": "Backend component",
			"lead": {"self": "https://api.tracker.yandex.net/v3/users/1", "id": "1", "display": "User One"},
			"assignAuto": false
		}`)
	})

	ctx := context.Background()
	component, resp, err := client.Components.Create(ctx, &ComponentRequest{
		Name:        Ptr("Backend"),
		Queue:       Ptr("TEST"),
		Description: Ptr("Backend component"),
		Lead:        Ptr("user1"),
		AssignAuto:  Ptr(false),
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusCreated)
	}
	if got := *component.ID; got != 1 {
		t.Errorf("ID = %d, want %d", got, 1)
	}
	if got := *component.Name; got != "Backend" {
		t.Errorf("Name = %q, want %q", got, "Backend")
	}
	if got := *component.Description; got != "Backend component" {
		t.Errorf("Description = %q, want %q", got, "Backend component")
	}
}

func TestComponentsService_Get(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/components/{id}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/components/1",
			"id": 1,
			"version": 1,
			"name": "Backend",
			"queue": {"self": "https://api.tracker.yandex.net/v3/queues/TEST", "id": "1", "key": "TEST", "display": "Test Queue"},
			"description": "Backend component",
			"lead": {"self": "https://api.tracker.yandex.net/v3/users/1", "id": "1", "display": "John"},
			"assignAuto": false
		}`)
	})

	ctx := context.Background()
	component, _, err := client.Components.Get(ctx, 1)
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if got := *component.ID; got != 1 {
		t.Errorf("ID = %d, want %d", got, 1)
	}
	if got := *component.Name; got != "Backend" {
		t.Errorf("Name = %q, want %q", got, "Backend")
	}
	if got := *component.Queue.Key; got != "TEST" {
		t.Errorf("Queue.Key = %q, want %q", got, "TEST")
	}
}

func TestComponentsService_Edit(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("PATCH /v3/components/{id}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testBody(t, r, `{"name":"Updated"}`)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/components/1",
			"id": 1,
			"version": 2,
			"name": "Updated",
			"queue": {"self": "https://api.tracker.yandex.net/v3/queues/TEST", "id": "1", "key": "TEST", "display": "Test Queue"},
			"description": "Backend component",
			"lead": {"self": "https://api.tracker.yandex.net/v3/users/1", "id": "1", "display": "John"},
			"assignAuto": false
		}`)
	})

	ctx := context.Background()
	component, _, err := client.Components.Edit(ctx, 1, &ComponentRequest{
		Name: Ptr("Updated"),
	})
	if err != nil {
		t.Fatalf("Edit returned error: %v", err)
	}
	if got := *component.ID; got != 1 {
		t.Errorf("ID = %d, want %d", got, 1)
	}
	if got := *component.Name; got != "Updated" {
		t.Errorf("Name = %q, want %q", got, "Updated")
	}
	if got := *component.Version; got != 2 {
		t.Errorf("Version = %d, want %d", got, 2)
	}
}

func TestComponentsService_Delete(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("DELETE /v3/components/{id}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	resp, err := client.Components.Delete(ctx, 1)
	if err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusNoContent)
	}
}

func TestComponentsService_Get_NotFound(t *testing.T) {
	client, mux := setup(t)
	ctx := context.Background()

	mux.HandleFunc("GET /v3/components/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"errorMessages":["Object not found"],"errors":{}}`)
	})

	_, _, err := client.Components.Get(ctx, 99999)
	if err == nil {
		t.Fatal("Components.Get expected error for 404, got nil")
	}
	if !IsNotFound(err) {
		t.Errorf("IsNotFound(err) = false, want true")
	}
}
