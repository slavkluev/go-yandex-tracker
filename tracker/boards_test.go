package tracker

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestBoardsService_Create(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/boards/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"name":"Sprint Board","boardType":"scrum","defaultQueue":{"key":"TEST"}}`)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/boards/1",
			"id": 1,
			"version": 1,
			"name": "Sprint Board",
			"columns": [{"self": "https://api.tracker.yandex.net/v3/boards/1/columns/1", "id": "1", "display": "Open"}],
			"defaultQueue": {"self": "https://api.tracker.yandex.net/v2/queues/TEST", "id": "1", "key": "TEST", "display": "Test Queue"},
			"createdBy": {"self": "https://api.tracker.yandex.net/v2/users/1", "id": "1", "display": "John"}
		}`)
	})

	ctx := context.Background()
	board, resp, err := client.Boards.Create(ctx, &BoardCreateRequest{
		Name:         Ptr("Sprint Board"),
		BoardType:    Ptr("scrum"),
		DefaultQueue: &DefaultQueue{Key: Ptr("TEST")},
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusCreated)
	}
	if got := *board.ID; got != 1 {
		t.Errorf("ID = %d, want %d", got, 1)
	}
	if got := *board.Name; got != "Sprint Board" {
		t.Errorf("Name = %q, want %q", got, "Sprint Board")
	}
	if got := *board.Version; got != 1 {
		t.Errorf("Version = %d, want %d", got, 1)
	}
	if got := *board.DefaultQueue.Key; got != "TEST" {
		t.Errorf("DefaultQueue.Key = %q, want %q", got, "TEST")
	}
}

func TestBoardsService_Get(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/boards/{id}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/boards/1",
			"id": 1,
			"version": 5,
			"name": "Test Board",
			"columns": [{"self": "https://api.tracker.yandex.net/v3/boards/1/columns/1", "id": "1", "display": "Open"}],
			"defaultQueue": {"self": "https://api.tracker.yandex.net/v2/queues/TEST", "id": "1", "key": "TEST", "display": "Test Queue"},
			"createdBy": {"self": "https://api.tracker.yandex.net/v2/users/1", "id": "1", "display": "John"}
		}`)
	})

	ctx := context.Background()
	board, _, err := client.Boards.Get(ctx, 1)
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if got := *board.ID; got != 1 {
		t.Errorf("ID = %d, want %d", got, 1)
	}
	if got := *board.Name; got != "Test Board" {
		t.Errorf("Name = %q, want %q", got, "Test Board")
	}
	if got := *board.Version; got != 5 {
		t.Errorf("Version = %d, want %d", got, 5)
	}
	if got := *board.DefaultQueue.Key; got != "TEST" {
		t.Errorf("DefaultQueue.Key = %q, want %q", got, "TEST")
	}
	if got := *board.CreatedBy.Display; got != "John" {
		t.Errorf("CreatedBy.Display = %q, want %q", got, "John")
	}
	if len(board.Columns) != 1 {
		t.Fatalf("len(Columns) = %d, want 1", len(board.Columns))
	}
	if got := *board.Columns[0].Display; got != "Open" {
		t.Errorf("Columns[0].Display = %q, want %q", got, "Open")
	}
}

func TestBoardsService_List(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/boards", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[
			{
				"self": "https://api.tracker.yandex.net/v3/boards/1",
				"id": 1,
				"version": 1,
				"name": "Board One",
				"defaultQueue": {"self": "https://api.tracker.yandex.net/v2/queues/TEST", "id": "1", "key": "TEST", "display": "Test Queue"},
				"createdBy": {"self": "https://api.tracker.yandex.net/v2/users/1", "id": "1", "display": "John"}
			},
			{
				"self": "https://api.tracker.yandex.net/v3/boards/2",
				"id": 2,
				"version": 3,
				"name": "Board Two",
				"defaultQueue": {"self": "https://api.tracker.yandex.net/v2/queues/DEV", "id": "2", "key": "DEV", "display": "Dev Queue"},
				"createdBy": {"self": "https://api.tracker.yandex.net/v2/users/2", "id": "2", "display": "Jane"}
			}
		]`)
	})

	ctx := context.Background()
	boards, _, err := client.Boards.List(ctx)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	if len(boards) != 2 {
		t.Fatalf("List returned %d boards, want 2", len(boards))
	}

	want := []*Board{
		{
			Self:         Ptr("https://api.tracker.yandex.net/v3/boards/1"),
			ID:           Ptr(1),
			Version:      Ptr(1),
			Name:         Ptr("Board One"),
			DefaultQueue: &Queue{Self: Ptr("https://api.tracker.yandex.net/v2/queues/TEST"), ID: Ptr("1"), Key: Ptr("TEST"), Display: Ptr("Test Queue")},
			CreatedBy:    &User{Self: Ptr("https://api.tracker.yandex.net/v2/users/1"), ID: Ptr("1"), Display: Ptr("John")},
		},
		{
			Self:         Ptr("https://api.tracker.yandex.net/v3/boards/2"),
			ID:           Ptr(2),
			Version:      Ptr(3),
			Name:         Ptr("Board Two"),
			DefaultQueue: &Queue{Self: Ptr("https://api.tracker.yandex.net/v2/queues/DEV"), ID: Ptr("2"), Key: Ptr("DEV"), Display: Ptr("Dev Queue")},
			CreatedBy:    &User{Self: Ptr("https://api.tracker.yandex.net/v2/users/2"), ID: Ptr("2"), Display: Ptr("Jane")},
		},
	}

	if !reflect.DeepEqual(boards, want) {
		t.Errorf("List returned %+v, want %+v", boards, want)
	}
}

func TestBoardsService_Edit(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("PATCH /v3/boards/{id}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testHeader(t, r, "If-Match", "5")
		testBody(t, r, `{"name":"Updated Board"}`)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/boards/1",
			"id": 1,
			"version": 6,
			"name": "Updated Board",
			"columns": [{"self": "https://api.tracker.yandex.net/v3/boards/1/columns/1", "id": "1", "display": "Open"}],
			"defaultQueue": {"self": "https://api.tracker.yandex.net/v2/queues/TEST", "id": "1", "key": "TEST", "display": "Test Queue"},
			"createdBy": {"self": "https://api.tracker.yandex.net/v2/users/1", "id": "1", "display": "John"}
		}`)
	})

	ctx := context.Background()
	board, _, err := client.Boards.Edit(ctx, 1, 5, &BoardEditRequest{
		Name: Ptr("Updated Board"),
	})
	if err != nil {
		t.Fatalf("Edit returned error: %v", err)
	}
	if got := *board.ID; got != 1 {
		t.Errorf("ID = %d, want %d", got, 1)
	}
	if got := *board.Name; got != "Updated Board" {
		t.Errorf("Name = %q, want %q", got, "Updated Board")
	}
	if got := *board.Version; got != 6 {
		t.Errorf("Version = %d, want %d", got, 6)
	}
}

func TestBoardsService_Delete(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("DELETE /v3/boards/{id}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	resp, err := client.Boards.Delete(ctx, 1)
	if err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusNoContent)
	}
}

func TestBoardsService_Get_NotFound(t *testing.T) {
	client, mux := setup(t)
	ctx := context.Background()

	mux.HandleFunc("GET /v3/boards/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"errorMessages":["Object not found"],"errors":{}}`)
	})

	_, _, err := client.Boards.Get(ctx, 99999)
	if err == nil {
		t.Fatal("Boards.Get expected error for 404, got nil")
	}
	if !IsNotFound(err) {
		t.Errorf("IsNotFound(err) = false, want true")
	}
}
