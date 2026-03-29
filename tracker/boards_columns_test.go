package tracker

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestBoardsService_ListColumns(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/boards/{boardID}/columns", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[
			{
				"self": "https://api.tracker.yandex.net/v3/boards/1/columns/1",
				"id": 1,
				"name": "Open",
				"statuses": [{"self": "https://api.tracker.yandex.net/v3/statuses/1", "id": "1", "key": "open", "display": "Open"}]
			},
			{
				"self": "https://api.tracker.yandex.net/v3/boards/1/columns/2",
				"id": 2,
				"name": "In Progress",
				"statuses": [{"self": "https://api.tracker.yandex.net/v3/statuses/2", "id": "2", "key": "inProgress", "display": "In Progress"}]
			}
		]`)
	})

	ctx := context.Background()
	columns, _, err := client.Boards.ListColumns(ctx, 1)
	if err != nil {
		t.Fatalf("ListColumns returned error: %v", err)
	}

	if len(columns) != 2 {
		t.Fatalf("ListColumns returned %d columns, want 2", len(columns))
	}

	want := []*Column{
		{
			Self: Ptr("https://api.tracker.yandex.net/v3/boards/1/columns/1"),
			ID:   Ptr(1),
			Name: Ptr("Open"),
			Statuses: []*Status{
				{Self: Ptr("https://api.tracker.yandex.net/v3/statuses/1"), ID: Ptr("1"), Key: Ptr("open"), Display: Ptr("Open")},
			},
		},
		{
			Self: Ptr("https://api.tracker.yandex.net/v3/boards/1/columns/2"),
			ID:   Ptr(2),
			Name: Ptr("In Progress"),
			Statuses: []*Status{
				{Self: Ptr("https://api.tracker.yandex.net/v3/statuses/2"), ID: Ptr("2"), Key: Ptr("inProgress"), Display: Ptr("In Progress")},
			},
		},
	}

	if !reflect.DeepEqual(columns, want) {
		t.Errorf("ListColumns returned %+v, want %+v", columns, want)
	}
}

func TestBoardsService_CreateColumn(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/boards/{boardID}/columns/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "If-Match", "5")
		testBody(t, r, `{"name":"Done","statuses":["closed","resolved"]}`)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/boards/1/columns/3",
			"id": 3,
			"name": "Done",
			"statuses": [
				{"self": "https://api.tracker.yandex.net/v3/statuses/3", "id": "3", "key": "closed", "display": "Closed"},
				{"self": "https://api.tracker.yandex.net/v3/statuses/4", "id": "4", "key": "resolved", "display": "Resolved"}
			]
		}`)
	})

	ctx := context.Background()
	column, resp, err := client.Boards.CreateColumn(ctx, 1, 5, &ColumnCreateRequest{
		Name:     Ptr("Done"),
		Statuses: []string{"closed", "resolved"},
	})
	if err != nil {
		t.Fatalf("CreateColumn returned error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusCreated)
	}
	if got := *column.ID; got != 3 {
		t.Errorf("ID = %d, want %d", got, 3)
	}
	if got := *column.Name; got != "Done" {
		t.Errorf("Name = %q, want %q", got, "Done")
	}
	if len(column.Statuses) != 2 {
		t.Fatalf("len(Statuses) = %d, want 2", len(column.Statuses))
	}
}

func TestBoardsService_EditColumn(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("PATCH /v3/boards/{boardID}/columns/{columnID}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testHeader(t, r, "If-Match", "5")
		testBody(t, r, `{"name":"Updated"}`)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/boards/1/columns/2",
			"id": 2,
			"name": "Updated",
			"statuses": [{"self": "https://api.tracker.yandex.net/v3/statuses/2", "id": "2", "key": "inProgress", "display": "In Progress"}]
		}`)
	})

	ctx := context.Background()
	column, _, err := client.Boards.EditColumn(ctx, 1, 2, 5, &ColumnEditRequest{
		Name: Ptr("Updated"),
	})
	if err != nil {
		t.Fatalf("EditColumn returned error: %v", err)
	}
	if got := *column.ID; got != 2 {
		t.Errorf("ID = %d, want %d", got, 2)
	}
	if got := *column.Name; got != "Updated" {
		t.Errorf("Name = %q, want %q", got, "Updated")
	}
}

func TestBoardsService_DeleteColumn(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("DELETE /v3/boards/{boardID}/columns/{columnID}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "If-Match", "5")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	resp, err := client.Boards.DeleteColumn(ctx, 1, 2, 5)
	if err != nil {
		t.Fatalf("DeleteColumn returned error: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusNoContent)
	}
}
