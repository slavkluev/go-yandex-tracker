package tracker

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestFiltersService_List(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/filters", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[
			{
				"self": "https://api.tracker.yandex.net/v3/filters/1",
				"id": 1,
				"name": "My Open Issues"
			},
			{
				"self": "https://api.tracker.yandex.net/v3/filters/2",
				"id": 2,
				"name": "Assigned to Me"
			}
		]`)
	})

	ctx := context.Background()
	filters, _, err := client.Filters.List(ctx)
	if err != nil {
		t.Fatalf("Filters.List returned error: %v", err)
	}

	if len(filters) != 2 {
		t.Fatalf("Filters.List returned %d filters, want 2", len(filters))
	}

	want := []*Filter{
		{
			Self: Ptr("https://api.tracker.yandex.net/v3/filters/1"),
			ID:   Ptr(1),
			Name: Ptr("My Open Issues"),
		},
		{
			Self: Ptr("https://api.tracker.yandex.net/v3/filters/2"),
			ID:   Ptr(2),
			Name: Ptr("Assigned to Me"),
		},
	}

	if !reflect.DeepEqual(filters, want) {
		t.Errorf("Filters.List returned %+v, want %+v", filters, want)
	}
}

func TestFiltersService_Get(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/filters/{id}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/filters/1",
			"id": 1,
			"name": "My Open Issues",
			"query": "Status: Open",
			"groupBy": "priority",
			"filter": {
				"status": "open"
			}
		}`)
	})

	ctx := context.Background()
	filter, _, err := client.Filters.Get(ctx, 1)
	if err != nil {
		t.Fatalf("Filters.Get returned error: %v", err)
	}

	want := &Filter{
		Self:    Ptr("https://api.tracker.yandex.net/v3/filters/1"),
		ID:      Ptr(1),
		Name:    Ptr("My Open Issues"),
		Query:   Ptr("Status: Open"),
		GroupBy: Ptr("priority"),
		Filter:  map[string]any{"status": "open"},
	}

	if !reflect.DeepEqual(filter, want) {
		t.Errorf("Filters.Get returned %+v, want %+v", filter, want)
	}
}

func TestFiltersService_Create(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/filters", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"name":"New Filter","query":"Status: Open"}`)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/filters/3",
			"id": 3,
			"name": "New Filter",
			"query": "Status: Open"
		}`)
	})

	ctx := context.Background()
	filter, resp, err := client.Filters.Create(ctx, &FilterRequest{
		Name:  Ptr("New Filter"),
		Query: Ptr("Status: Open"),
	})
	if err != nil {
		t.Fatalf("Filters.Create returned error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusCreated)
	}

	want := &Filter{
		Self:  Ptr("https://api.tracker.yandex.net/v3/filters/3"),
		ID:    Ptr(3),
		Name:  Ptr("New Filter"),
		Query: Ptr("Status: Open"),
	}

	if !reflect.DeepEqual(filter, want) {
		t.Errorf("Filters.Create returned %+v, want %+v", filter, want)
	}
}

func TestFiltersService_Update(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("PATCH /v3/filters/{id}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testBody(t, r, `{"name":"Updated Filter"}`)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/filters/1",
			"id": 1,
			"name": "Updated Filter",
			"query": "Status: Open"
		}`)
	})

	ctx := context.Background()
	filter, _, err := client.Filters.Update(ctx, 1, &FilterRequest{
		Name: Ptr("Updated Filter"),
	})
	if err != nil {
		t.Fatalf("Filters.Update returned error: %v", err)
	}

	want := &Filter{
		Self:  Ptr("https://api.tracker.yandex.net/v3/filters/1"),
		ID:    Ptr(1),
		Name:  Ptr("Updated Filter"),
		Query: Ptr("Status: Open"),
	}

	if !reflect.DeepEqual(filter, want) {
		t.Errorf("Filters.Update returned %+v, want %+v", filter, want)
	}
}

func TestFiltersService_Delete(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("DELETE /v3/filters/{id}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	resp, err := client.Filters.Delete(ctx, 1)
	if err != nil {
		t.Fatalf("Filters.Delete returned error: %v", err)
	}
	if resp == nil {
		t.Fatal("Filters.Delete returned nil response")
	}
}
