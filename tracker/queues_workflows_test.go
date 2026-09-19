package tracker

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestQueuesService_ListWorkflows(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/queues/{key}/workflows", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		if got := r.PathValue("key"); got != "TEST" {
			t.Errorf("queue key = %q, want %q", got, "TEST")
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"dev": [
				{
					"self": "https://api.tracker.yandex.net/v3/issuetypes/1",
					"id": "1",
					"key": "task",
					"display": "Issue"
				}
			]
		}`)
	})

	workflows, _, err := client.Queues.ListWorkflows(ctx, "TEST")
	if err != nil {
		t.Fatalf("Queues.ListWorkflows returned error: %v", err)
	}

	want := map[string][]*IssueType{
		"dev": {
			{
				Self:    Ptr("https://api.tracker.yandex.net/v3/issuetypes/1"),
				ID:      Ptr(FlexString("1")),
				Key:     Ptr("task"),
				Display: Ptr("Issue"),
			},
		},
	}

	if !reflect.DeepEqual(workflows, want) {
		t.Errorf("Queues.ListWorkflows returned %+v, want %+v", workflows, want)
	}
}

func TestQueuesService_ListWorkflows_NotFound(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/queues/{key}/workflows", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"errorMessages":["Object not found"],"errors":{}}`)
	})

	_, _, err := client.Queues.ListWorkflows(ctx, "NONEXIST")
	if err == nil {
		t.Fatal("Queues.ListWorkflows expected error for 404, got nil")
	}
	if !IsNotFound(err) {
		t.Errorf("IsNotFound(err) = false, want true")
	}
}
