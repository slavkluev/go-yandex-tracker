package tracker

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestPrioritiesService_List(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/priorities", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[
			{
				"self": "https://api.tracker.yandex.net/v3/priorities/1",
				"id": "1",
				"key": "critical",
				"name": "Critical",
				"version": 1,
				"order": 1
			},
			{
				"self": "https://api.tracker.yandex.net/v3/priorities/2",
				"id": "2",
				"key": "normal",
				"name": "Normal",
				"version": 1,
				"order": 2
			}
		]`)
	})

	ctx := context.Background()
	priorities, _, err := client.Priorities.List(ctx, nil)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	if len(priorities) != 2 {
		t.Fatalf("List returned %d priorities, want 2", len(priorities))
	}

	want := &Priority{
		Self:    Ptr("https://api.tracker.yandex.net/v3/priorities/1"),
		ID:      Ptr("1"),
		Key:     Ptr("critical"),
		Name:    Ptr("Critical"),
		Version: Ptr(1),
		Order:   Ptr(1),
	}

	if !reflect.DeepEqual(priorities[0], want) {
		t.Errorf("List returned %+v, want %+v", priorities[0], want)
	}
}

func TestPrioritiesService_List_Localized(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/priorities", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		if got := r.URL.Query().Get("localized"); got != "true" {
			t.Errorf("Query param localized = %q, want %q", got, "true")
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[
			{
				"self": "https://api.tracker.yandex.net/v3/priorities/1",
				"id": "1",
				"key": "critical",
				"name": "Critical",
				"version": 1,
				"order": 1
			}
		]`)
	})

	ctx := context.Background()
	priorities, _, err := client.Priorities.List(ctx, &PriorityListOptions{
		Localized: Ptr(true),
	})
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	if len(priorities) != 1 {
		t.Fatalf("List returned %d priorities, want 1", len(priorities))
	}

	want := &Priority{
		Self:    Ptr("https://api.tracker.yandex.net/v3/priorities/1"),
		ID:      Ptr("1"),
		Key:     Ptr("critical"),
		Name:    Ptr("Critical"),
		Version: Ptr(1),
		Order:   Ptr(1),
	}

	if !reflect.DeepEqual(priorities[0], want) {
		t.Errorf("List returned %+v, want %+v", priorities[0], want)
	}
}
