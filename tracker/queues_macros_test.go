package tracker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestQueuesService_ListMacros(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/queues/{key}/macros", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[
			{
				"id": 1,
				"self": "https://api.tracker.yandex.net/v3/queues/TEST/macros/1",
				"queue": {
					"self": "https://api.tracker.yandex.net/v3/queues/TEST",
					"id": "100",
					"key": "TEST",
					"display": "Test Queue"
				},
				"name": "Close issue",
				"body": "Closing this issue."
			},
			{
				"id": 2,
				"self": "https://api.tracker.yandex.net/v3/queues/TEST/macros/2",
				"queue": {
					"self": "https://api.tracker.yandex.net/v3/queues/TEST",
					"id": "100",
					"key": "TEST",
					"display": "Test Queue"
				},
				"name": "Assign to me",
				"body": "Taking ownership."
			}
		]`)
	})

	macros, _, err := client.Queues.ListMacros(ctx, "TEST")
	if err != nil {
		t.Fatalf("Queues.ListMacros returned error: %v", err)
	}

	if len(macros) != 2 {
		t.Fatalf("Queues.ListMacros returned %d macros, want 2", len(macros))
	}

	want := []*Macro{
		{
			ID:   Ptr(1),
			Self: Ptr("https://api.tracker.yandex.net/v3/queues/TEST/macros/1"),
			Queue: &Queue{
				Self:    Ptr("https://api.tracker.yandex.net/v3/queues/TEST"),
				ID:      Ptr("100"),
				Key:     Ptr("TEST"),
				Display: Ptr("Test Queue"),
			},
			Name: Ptr("Close issue"),
			Body: Ptr("Closing this issue."),
		},
		{
			ID:   Ptr(2),
			Self: Ptr("https://api.tracker.yandex.net/v3/queues/TEST/macros/2"),
			Queue: &Queue{
				Self:    Ptr("https://api.tracker.yandex.net/v3/queues/TEST"),
				ID:      Ptr("100"),
				Key:     Ptr("TEST"),
				Display: Ptr("Test Queue"),
			},
			Name: Ptr("Assign to me"),
			Body: Ptr("Taking ownership."),
		},
	}

	if !reflect.DeepEqual(macros, want) {
		t.Errorf("Queues.ListMacros returned %+v, want %+v", macros, want)
	}
}

func TestQueuesService_GetMacro(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/queues/{key}/macros/{id}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"id": 5,
			"self": "https://api.tracker.yandex.net/v3/queues/TEST/macros/5",
			"queue": {
				"self": "https://api.tracker.yandex.net/v3/queues/TEST",
				"id": "100",
				"key": "TEST",
				"display": "Test Queue"
			},
			"name": "Close and comment",
			"body": "Resolved per policy.",
			"issueUpdate": [
				{
					"field": {
						"self": "https://api.tracker.yandex.net/v3/fields/status",
						"id": "status",
						"display": "Status"
					},
					"update": {"set": "closed"}
				}
			]
		}`)
	})

	macro, _, err := client.Queues.GetMacro(ctx, "TEST", 5)
	if err != nil {
		t.Fatalf("Queues.GetMacro returned error: %v", err)
	}

	want := &Macro{
		ID:   Ptr(5),
		Self: Ptr("https://api.tracker.yandex.net/v3/queues/TEST/macros/5"),
		Queue: &Queue{
			Self:    Ptr("https://api.tracker.yandex.net/v3/queues/TEST"),
			ID:      Ptr("100"),
			Key:     Ptr("TEST"),
			Display: Ptr("Test Queue"),
		},
		Name: Ptr("Close and comment"),
		Body: Ptr("Resolved per policy."),
		IssueUpdate: []*MacroIssueUpdate{
			{
				Field: &MacroIssueUpdateField{
					Self:    Ptr("https://api.tracker.yandex.net/v3/fields/status"),
					ID:      Ptr("status"),
					Display: Ptr("Status"),
				},
				Update: map[string]any{"set": "closed"},
			},
		},
	}

	if !reflect.DeepEqual(macro, want) {
		t.Errorf("Queues.GetMacro returned %+v, want %+v", macro, want)
	}
}

func TestQueuesService_CreateMacro(t *testing.T) {
	client, mux := setup(t)

	input := &MacroCreateRequest{
		Name: Ptr("New macro"),
		Body: Ptr("Auto comment"),
		IssueUpdate: []map[string]any{
			{"field": "status", "set": "closed"},
		},
	}

	mux.HandleFunc("POST /v3/queues/{key}/macros", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		// Verify request body.
		var got MacroCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("Decode request body: %v", err)
		}
		if *got.Name != "New macro" {
			t.Errorf("Request body name = %q, want %q", *got.Name, "New macro")
		}

		// Return 201 Created (not 200).
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"id": 10,
			"self": "https://api.tracker.yandex.net/v3/queues/TEST/macros/10",
			"queue": {
				"self": "https://api.tracker.yandex.net/v3/queues/TEST",
				"id": "100",
				"key": "TEST",
				"display": "Test Queue"
			},
			"name": "New macro",
			"body": "Auto comment"
		}`)
	})

	macro, _, err := client.Queues.CreateMacro(ctx, "TEST", input)
	if err != nil {
		t.Fatalf("Queues.CreateMacro returned error: %v", err)
	}

	want := &Macro{
		ID:   Ptr(10),
		Self: Ptr("https://api.tracker.yandex.net/v3/queues/TEST/macros/10"),
		Queue: &Queue{
			Self:    Ptr("https://api.tracker.yandex.net/v3/queues/TEST"),
			ID:      Ptr("100"),
			Key:     Ptr("TEST"),
			Display: Ptr("Test Queue"),
		},
		Name: Ptr("New macro"),
		Body: Ptr("Auto comment"),
	}

	if !reflect.DeepEqual(macro, want) {
		t.Errorf("Queues.CreateMacro returned %+v, want %+v", macro, want)
	}
}

func TestQueuesService_EditMacro(t *testing.T) {
	client, mux := setup(t)

	input := &MacroEditRequest{
		Name: Ptr("Updated macro"),
		Body: Ptr("Updated comment"),
	}

	mux.HandleFunc("PATCH /v3/queues/{key}/macros/{id}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")

		// Verify request body.
		var got MacroEditRequest
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("Decode request body: %v", err)
		}
		if !reflect.DeepEqual(&got, input) {
			t.Errorf("Request body = %+v, want %+v", got, input)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"id": 5,
			"self": "https://api.tracker.yandex.net/v3/queues/TEST/macros/5",
			"queue": {
				"self": "https://api.tracker.yandex.net/v3/queues/TEST",
				"id": "100",
				"key": "TEST",
				"display": "Test Queue"
			},
			"name": "Updated macro",
			"body": "Updated comment"
		}`)
	})

	macro, _, err := client.Queues.EditMacro(ctx, "TEST", 5, input)
	if err != nil {
		t.Fatalf("Queues.EditMacro returned error: %v", err)
	}

	want := &Macro{
		ID:   Ptr(5),
		Self: Ptr("https://api.tracker.yandex.net/v3/queues/TEST/macros/5"),
		Queue: &Queue{
			Self:    Ptr("https://api.tracker.yandex.net/v3/queues/TEST"),
			ID:      Ptr("100"),
			Key:     Ptr("TEST"),
			Display: Ptr("Test Queue"),
		},
		Name: Ptr("Updated macro"),
		Body: Ptr("Updated comment"),
	}

	if !reflect.DeepEqual(macro, want) {
		t.Errorf("Queues.EditMacro returned %+v, want %+v", macro, want)
	}
}

func TestQueuesService_DeleteMacro(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("DELETE /v3/queues/{key}/macros/{id}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.Queues.DeleteMacro(ctx, "TEST", 5)
	if err != nil {
		t.Fatalf("Queues.DeleteMacro returned error: %v", err)
	}
	if resp == nil {
		t.Fatal("Queues.DeleteMacro returned nil response")
	}
}
