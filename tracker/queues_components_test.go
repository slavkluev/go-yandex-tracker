package tracker

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestQueuesService_ListComponents(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/queues/{key}/components", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		if got := r.PathValue("key"); got != "TEST" {
			t.Errorf("queue key = %q, want %q", got, "TEST")
		}
		if got := r.URL.Query().Get("fields"); got != "name" {
			t.Errorf("fields = %q, want %q", got, "name")
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[
			{
				"self": "https://api.tracker.yandex.net/v3/components/123",
				"id": 123,
				"name": "Component",
				"queue": {"self": "https://api.tracker.yandex.net/v3/queues/TEST", "id": "1", "key": "TEST", "display": "Test queue"}
			}
		]`)
	})

	components, _, err := client.Queues.ListComponents(ctx, "TEST", &QueueComponentsListOptions{Fields: "name"})
	if err != nil {
		t.Fatalf("Queues.ListComponents returned error: %v", err)
	}

	want := []*Component{
		{
			Self: Ptr("https://api.tracker.yandex.net/v3/components/123"),
			ID:   Ptr(FlexString("123")),
			Name: Ptr("Component"),
			Queue: &Queue{
				Self:    Ptr("https://api.tracker.yandex.net/v3/queues/TEST"),
				ID:      Ptr(FlexString("1")),
				Key:     Ptr("TEST"),
				Display: Ptr("Test queue"),
			},
		},
	}

	if !reflect.DeepEqual(components, want) {
		t.Errorf("Queues.ListComponents returned %+v, want %+v", components, want)
	}
}

func TestQueuesService_ListComponents_NoOptions(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/queues/{key}/components", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		if got := r.URL.RawQuery; got != "" {
			t.Errorf("query = %q, want none", got)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{"self": "https://api.tracker.yandex.net/v3/components/55", "id": 55, "version": 1, "name": "Expedite",
		  "queue": {"self": "https://api.tracker.yandex.net/v3/queues/MTP", "id": "140", "key": "MTP", "display": "Metal trading platform"},
		  "assignAuto": false}]`)
	})

	components, _, err := client.Queues.ListComponents(ctx, "MTP", nil)
	if err != nil {
		t.Fatalf("Queues.ListComponents returned error: %v", err)
	}

	want := []*Component{
		{
			Self:    Ptr("https://api.tracker.yandex.net/v3/components/55"),
			ID:      Ptr(FlexString("55")),
			Version: Ptr(FlexString("1")),
			Name:    Ptr("Expedite"),
			Queue: &Queue{
				Self:    Ptr("https://api.tracker.yandex.net/v3/queues/MTP"),
				ID:      Ptr(FlexString("140")),
				Key:     Ptr("MTP"),
				Display: Ptr("Metal trading platform"),
			},
			AssignAuto: Ptr(false),
		},
	}

	if !reflect.DeepEqual(components, want) {
		t.Errorf("Queues.ListComponents returned %+v, want %+v", components, want)
	}
}

func TestQueuesService_ListComponents_NotFound(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/queues/{key}/components", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"errorMessages":["Object not found"],"errors":{}}`)
	})

	_, _, err := client.Queues.ListComponents(ctx, "NONEXIST", nil)
	if err == nil {
		t.Fatal("Queues.ListComponents expected error for 404, got nil")
	}
	if !IsNotFound(err) {
		t.Errorf("IsNotFound(err) = false, want true")
	}
}
