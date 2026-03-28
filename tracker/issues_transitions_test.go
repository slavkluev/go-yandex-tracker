package tracker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestIssuesService_GetTransitions(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v2/issues/{key}/transitions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{"id":"2","display":"Close","to":{"key":"closed"}}]`)
	})

	ctx := context.Background()
	transitions, _, err := client.Issues.GetTransitions(ctx, "QUEUE-1")
	if err != nil {
		t.Fatalf("GetTransitions returned error: %v", err)
	}

	if len(transitions) != 1 {
		t.Fatalf("GetTransitions returned %d transitions, want 1", len(transitions))
	}
	want := &Transition{
		ID:      Ptr("2"),
		Display: Ptr("Close"),
		To:      &Status{Key: Ptr("closed")},
	}
	if !reflect.DeepEqual(transitions[0], want) {
		t.Errorf("GetTransitions[0] = %+v, want %+v", transitions[0], want)
	}
}

func TestIssuesService_GetTransitions_Empty(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v2/issues/{key}/transitions", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[]`)
	})

	ctx := context.Background()
	transitions, _, err := client.Issues.GetTransitions(ctx, "QUEUE-1")
	if err != nil {
		t.Fatalf("GetTransitions returned error: %v", err)
	}
	if len(transitions) != 0 {
		t.Errorf("GetTransitions returned %d transitions, want 0", len(transitions))
	}
}

func TestIssuesService_ExecuteTransition(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v2/issues/{key}/transitions/{id}/_execute", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("Error decoding request body: %v", err)
		}
		if body["comment"] != "done" {
			t.Errorf("request body comment = %v, want %q", body["comment"], "done")
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{"id":"3","display":"Reopen","to":{"key":"open"}}]`)
	})

	ctx := context.Background()
	req := &TransitionRequest{Comment: Ptr("done")}
	transitions, _, err := client.Issues.ExecuteTransition(ctx, "QUEUE-1", "2", req)
	if err != nil {
		t.Fatalf("ExecuteTransition returned error: %v", err)
	}

	if len(transitions) != 1 {
		t.Fatalf("ExecuteTransition returned %d transitions, want 1", len(transitions))
	}
	want := &Transition{
		ID:      Ptr("3"),
		Display: Ptr("Reopen"),
		To:      &Status{Key: Ptr("open")},
	}
	if !reflect.DeepEqual(transitions[0], want) {
		t.Errorf("ExecuteTransition[0] = %+v, want %+v", transitions[0], want)
	}
}

func TestIssuesService_ExecuteTransition_NilBody(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v2/issues/{key}/transitions/{id}/_execute", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{"id":"2","display":"Close","to":{"key":"closed"}}]`)
	})

	ctx := context.Background()
	transitions, _, err := client.Issues.ExecuteTransition(ctx, "QUEUE-1", "2", nil)
	if err != nil {
		t.Fatalf("ExecuteTransition returned error: %v", err)
	}

	if len(transitions) != 1 {
		t.Fatalf("ExecuteTransition returned %d transitions, want 1", len(transitions))
	}
	if *transitions[0].ID != "2" {
		t.Errorf("ExecuteTransition[0].ID = %q, want %q", *transitions[0].ID, "2")
	}
}
