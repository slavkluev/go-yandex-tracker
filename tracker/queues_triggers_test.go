package tracker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestQueuesService_ListTriggers(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/queues/{key}/triggers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[
			{
				"id": "1",
				"self": "https://api.tracker.yandex.net/v3/queues/TEST/triggers/1",
				"queue": {
					"self": "https://api.tracker.yandex.net/v2/queues/TEST",
					"id": "100",
					"key": "TEST",
					"display": "Test Queue"
				},
				"name": "Auto-close trigger",
				"order": "0.0001",
				"actions": [
					{"type": "Transition", "id": "1", "status": {"self": "https://api.tracker.yandex.net/v2/statuses/closed", "key": "closed"}}
				],
				"conditions": [
					{"type": "Event.create"}
				],
				"version": 1,
				"active": true
			}
		]`)
	})

	triggers, _, err := client.Queues.ListTriggers(ctx, "TEST")
	if err != nil {
		t.Fatalf("Queues.ListTriggers returned error: %v", err)
	}

	if len(triggers) != 1 {
		t.Fatalf("Queues.ListTriggers returned %d triggers, want 1", len(triggers))
	}

	want := []*Trigger{
		{
			ID:   Ptr("1"),
			Self: Ptr("https://api.tracker.yandex.net/v3/queues/TEST/triggers/1"),
			Queue: &Queue{
				Self:    Ptr("https://api.tracker.yandex.net/v2/queues/TEST"),
				ID:      Ptr("100"),
				Key:     Ptr("TEST"),
				Display: Ptr("Test Queue"),
			},
			Name:  Ptr("Auto-close trigger"),
			Order: Ptr("0.0001"),
			Actions: []*AutomationAction{
				{
					Type:   Ptr("Transition"),
					ID:     Ptr("1"),
					Status: &Status{Self: Ptr("https://api.tracker.yandex.net/v2/statuses/closed"), Key: Ptr("closed")},
				},
			},
			Conditions: []*AutomationCondition{
				{
					Type: Ptr("Event.create"),
				},
			},
			Version: Ptr(1),
			Active:  Ptr(true),
		},
	}

	if !reflect.DeepEqual(triggers, want) {
		t.Errorf("Queues.ListTriggers returned %+v, want %+v", triggers, want)
	}
}

func TestQueuesService_GetTrigger(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/queues/{key}/triggers/{id}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"id": "1",
			"self": "https://api.tracker.yandex.net/v3/queues/TEST/triggers/1",
			"queue": {
				"self": "https://api.tracker.yandex.net/v2/queues/TEST",
				"id": "100",
				"key": "TEST",
				"display": "Test Queue"
			},
			"name": "Auto-close trigger",
			"order": "0.0001",
			"actions": [
				{"type": "Transition", "id": "1", "status": {"self": "https://api.tracker.yandex.net/v2/statuses/closed", "key": "closed"}}
			],
			"conditions": [
				{"type": "Event.create"}
			],
			"version": 1,
			"active": true
		}`)
	})

	trigger, _, err := client.Queues.GetTrigger(ctx, "TEST", 1)
	if err != nil {
		t.Fatalf("Queues.GetTrigger returned error: %v", err)
	}

	want := &Trigger{
		ID:   Ptr("1"),
		Self: Ptr("https://api.tracker.yandex.net/v3/queues/TEST/triggers/1"),
		Queue: &Queue{
			Self:    Ptr("https://api.tracker.yandex.net/v2/queues/TEST"),
			ID:      Ptr("100"),
			Key:     Ptr("TEST"),
			Display: Ptr("Test Queue"),
		},
		Name:  Ptr("Auto-close trigger"),
		Order: Ptr("0.0001"),
		Actions: []*AutomationAction{
			{
				Type:   Ptr("Transition"),
				ID:     Ptr("1"),
				Status: &Status{Self: Ptr("https://api.tracker.yandex.net/v2/statuses/closed"), Key: Ptr("closed")},
			},
		},
		Conditions: []*AutomationCondition{
			{
				Type: Ptr("Event.create"),
			},
		},
		Version: Ptr(1),
		Active:  Ptr(true),
	}

	if !reflect.DeepEqual(trigger, want) {
		t.Errorf("Queues.GetTrigger returned %+v, want %+v", trigger, want)
	}
}

func TestQueuesService_CreateTrigger(t *testing.T) {
	client, mux := setup(t)

	input := &TriggerCreateRequest{
		Name: Ptr("New trigger"),
		Actions: []*AutomationAction{
			{Type: Ptr("Transition"), Status: &Status{Key: Ptr("closed")}},
		},
		Conditions: []*AutomationCondition{
			{Type: Ptr("Event.create")},
		},
		Active: Ptr(true),
	}

	mux.HandleFunc("POST /v3/queues/{key}/triggers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		// Verify the request body is correct.
		var got TriggerCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("Decode request body: %v", err)
		}
		if !reflect.DeepEqual(&got, input) {
			t.Errorf("Request body = %+v, want %+v", got, input)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"id": "2",
			"self": "https://api.tracker.yandex.net/v3/queues/TEST/triggers/2",
			"queue": {
				"self": "https://api.tracker.yandex.net/v2/queues/TEST",
				"id": "100",
				"key": "TEST",
				"display": "Test Queue"
			},
			"name": "New trigger",
			"order": "0.0002",
			"actions": [
				{"type": "Transition", "status": {"key": "closed"}}
			],
			"conditions": [
				{"type": "Event.create"}
			],
			"version": 1,
			"active": true
		}`)
	})

	trigger, _, err := client.Queues.CreateTrigger(ctx, "TEST", input)
	if err != nil {
		t.Fatalf("Queues.CreateTrigger returned error: %v", err)
	}

	want := &Trigger{
		ID:   Ptr("2"),
		Self: Ptr("https://api.tracker.yandex.net/v3/queues/TEST/triggers/2"),
		Queue: &Queue{
			Self:    Ptr("https://api.tracker.yandex.net/v2/queues/TEST"),
			ID:      Ptr("100"),
			Key:     Ptr("TEST"),
			Display: Ptr("Test Queue"),
		},
		Name:  Ptr("New trigger"),
		Order: Ptr("0.0002"),
		Actions: []*AutomationAction{
			{
				Type:   Ptr("Transition"),
				Status: &Status{Key: Ptr("closed")},
			},
		},
		Conditions: []*AutomationCondition{
			{
				Type: Ptr("Event.create"),
			},
		},
		Version: Ptr(1),
		Active:  Ptr(true),
	}

	if !reflect.DeepEqual(trigger, want) {
		t.Errorf("Queues.CreateTrigger returned %+v, want %+v", trigger, want)
	}
}

func TestQueuesService_UpdateTrigger(t *testing.T) {
	client, mux := setup(t)

	input := &TriggerUpdateRequest{
		Name:   Ptr("Updated trigger"),
		Active: Ptr(false),
	}

	mux.HandleFunc("PATCH /v3/queues/{key}/triggers/{id}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")

		// Verify version query parameter is present.
		if got := r.URL.Query().Get("version"); got != "5" {
			t.Errorf("version query param = %q, want %q", got, "5")
		}

		// Verify request body.
		var got TriggerUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("Decode request body: %v", err)
		}
		if !reflect.DeepEqual(&got, input) {
			t.Errorf("Request body = %+v, want %+v", got, input)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"id": "1",
			"self": "https://api.tracker.yandex.net/v3/queues/TEST/triggers/1",
			"queue": {
				"self": "https://api.tracker.yandex.net/v2/queues/TEST",
				"id": "100",
				"key": "TEST",
				"display": "Test Queue"
			},
			"name": "Updated trigger",
			"order": "0.0001",
			"actions": [
				{"type": "Transition", "id": "1", "status": {"self": "https://api.tracker.yandex.net/v2/statuses/closed", "key": "closed"}}
			],
			"conditions": [
				{"type": "Event.create"}
			],
			"version": 6,
			"active": false
		}`)
	})

	trigger, _, err := client.Queues.UpdateTrigger(ctx, "TEST", 1, &TriggerUpdateOptions{Version: 5}, input)
	if err != nil {
		t.Fatalf("Queues.UpdateTrigger returned error: %v", err)
	}

	want := &Trigger{
		ID:   Ptr("1"),
		Self: Ptr("https://api.tracker.yandex.net/v3/queues/TEST/triggers/1"),
		Queue: &Queue{
			Self:    Ptr("https://api.tracker.yandex.net/v2/queues/TEST"),
			ID:      Ptr("100"),
			Key:     Ptr("TEST"),
			Display: Ptr("Test Queue"),
		},
		Name:  Ptr("Updated trigger"),
		Order: Ptr("0.0001"),
		Actions: []*AutomationAction{
			{
				Type:   Ptr("Transition"),
				ID:     Ptr("1"),
				Status: &Status{Self: Ptr("https://api.tracker.yandex.net/v2/statuses/closed"), Key: Ptr("closed")},
			},
		},
		Conditions: []*AutomationCondition{
			{
				Type: Ptr("Event.create"),
			},
		},
		Version: Ptr(6),
		Active:  Ptr(false),
	}

	if !reflect.DeepEqual(trigger, want) {
		t.Errorf("Queues.UpdateTrigger returned %+v, want %+v", trigger, want)
	}
}
