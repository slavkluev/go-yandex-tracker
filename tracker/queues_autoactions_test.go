package tracker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestQueuesService_ListAutoActions(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/queues/{key}/autoactions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[
			{
				"id": "1",
				"self": "https://api.tracker.yandex.net/v3/queues/TEST/autoactions/1",
				"queue": {
					"self": "https://api.tracker.yandex.net/v2/queues/TEST",
					"id": "100",
					"key": "TEST",
					"display": "Test Queue"
				},
				"name": "Close stale issues",
				"version": 3,
				"active": true,
				"actions": [
					{"type": "Transition", "id": "1", "status": {"self": "https://api.tracker.yandex.net/v2/statuses/closed", "key": "closed"}}
				],
				"enableNotifications": false,
				"intervalMillis": 3600000,
				"calendar": {"id": 1}
			}
		]`)
	})

	autoActions, _, err := client.Queues.ListAutoActions(ctx, "TEST")
	if err != nil {
		t.Fatalf("Queues.ListAutoActions returned error: %v", err)
	}

	if len(autoActions) != 1 {
		t.Fatalf("Queues.ListAutoActions returned %d auto-actions, want 1", len(autoActions))
	}

	want := []*AutoAction{
		{
			ID:   Ptr("1"),
			Self: Ptr("https://api.tracker.yandex.net/v3/queues/TEST/autoactions/1"),
			Queue: &Queue{
				Self:    Ptr("https://api.tracker.yandex.net/v2/queues/TEST"),
				ID:      Ptr("100"),
				Key:     Ptr("TEST"),
				Display: Ptr("Test Queue"),
			},
			Name:    Ptr("Close stale issues"),
			Version: Ptr(3),
			Active:  Ptr(true),
			Actions: []*AutomationAction{
				{
					Type:   Ptr("Transition"),
					ID:     Ptr("1"),
					Status: &Status{Self: Ptr("https://api.tracker.yandex.net/v2/statuses/closed"), Key: Ptr("closed")},
				},
			},
			EnableNotifications: Ptr(false),
			IntervalMillis:      Ptr(int64(3600000)),
			Calendar:            &AutoActionCalendar{ID: Ptr(1)},
		},
	}

	if !reflect.DeepEqual(autoActions, want) {
		t.Errorf("Queues.ListAutoActions returned %+v, want %+v", autoActions, want)
	}
}

func TestQueuesService_GetAutoAction(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/queues/{key}/autoactions/{id}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"id": "1",
			"self": "https://api.tracker.yandex.net/v3/queues/TEST/autoactions/1",
			"queue": {
				"self": "https://api.tracker.yandex.net/v2/queues/TEST",
				"id": "100",
				"key": "TEST",
				"display": "Test Queue"
			},
			"name": "Close stale issues",
			"version": 3,
			"active": true,
			"actions": [
				{"type": "Transition", "id": "1", "status": {"self": "https://api.tracker.yandex.net/v2/statuses/closed", "key": "closed"}}
			],
			"enableNotifications": false,
			"intervalMillis": 3600000,
			"calendar": {"id": 1},
			"totalIssuesProcessed": 42
		}`)
	})

	autoAction, _, err := client.Queues.GetAutoAction(ctx, "TEST", 1)
	if err != nil {
		t.Fatalf("Queues.GetAutoAction returned error: %v", err)
	}

	want := &AutoAction{
		ID:   Ptr("1"),
		Self: Ptr("https://api.tracker.yandex.net/v3/queues/TEST/autoactions/1"),
		Queue: &Queue{
			Self:    Ptr("https://api.tracker.yandex.net/v2/queues/TEST"),
			ID:      Ptr("100"),
			Key:     Ptr("TEST"),
			Display: Ptr("Test Queue"),
		},
		Name:    Ptr("Close stale issues"),
		Version: Ptr(3),
		Active:  Ptr(true),
		Actions: []*AutomationAction{
			{
				Type:   Ptr("Transition"),
				ID:     Ptr("1"),
				Status: &Status{Self: Ptr("https://api.tracker.yandex.net/v2/statuses/closed"), Key: Ptr("closed")},
			},
		},
		EnableNotifications:  Ptr(false),
		IntervalMillis:       Ptr(int64(3600000)),
		Calendar:             &AutoActionCalendar{ID: Ptr(1)},
		TotalIssuesProcessed: Ptr(42),
	}

	if !reflect.DeepEqual(autoAction, want) {
		t.Errorf("Queues.GetAutoAction returned %+v, want %+v", autoAction, want)
	}
}

func TestQueuesService_CreateAutoAction(t *testing.T) {
	client, mux := setup(t)

	input := &AutoActionCreateRequest{
		Name: Ptr("Close stale issues"),
		Filter: map[string]any{
			"assignee": "user1",
		},
		Actions: []*AutomationAction{
			{Type: Ptr("Transition"), Status: &Status{Key: Ptr("closed")}},
		},
		Active:         Ptr(true),
		IntervalMillis: Ptr(int64(3600000)),
	}

	mux.HandleFunc("POST /v3/queues/{key}/autoactions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		// Verify request body.
		var got map[string]any
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("Decode request body: %v", err)
		}
		if got["name"] != "Close stale issues" {
			t.Errorf("Request body name = %v, want %q", got["name"], "Close stale issues")
		}
		filter, ok := got["filter"].(map[string]any)
		if !ok {
			t.Fatalf("Request body filter is not map, got %T", got["filter"])
		}
		if filter["assignee"] != "user1" {
			t.Errorf("Request body filter.assignee = %v, want %q", filter["assignee"], "user1")
		}
		if got["intervalMillis"] != float64(3600000) {
			t.Errorf("Request body intervalMillis = %v, want %v", got["intervalMillis"], 3600000)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"id": "2",
			"self": "https://api.tracker.yandex.net/v3/queues/TEST/autoactions/2",
			"queue": {
				"self": "https://api.tracker.yandex.net/v2/queues/TEST",
				"id": "100",
				"key": "TEST",
				"display": "Test Queue"
			},
			"name": "Close stale issues",
			"version": 1,
			"active": true,
			"filter": {"assignee": "user1"},
			"actions": [
				{"type": "Transition", "status": {"key": "closed"}}
			],
			"intervalMillis": 3600000
		}`)
	})

	autoAction, _, err := client.Queues.CreateAutoAction(ctx, "TEST", input)
	if err != nil {
		t.Fatalf("Queues.CreateAutoAction returned error: %v", err)
	}

	want := &AutoAction{
		ID:   Ptr("2"),
		Self: Ptr("https://api.tracker.yandex.net/v3/queues/TEST/autoactions/2"),
		Queue: &Queue{
			Self:    Ptr("https://api.tracker.yandex.net/v2/queues/TEST"),
			ID:      Ptr("100"),
			Key:     Ptr("TEST"),
			Display: Ptr("Test Queue"),
		},
		Name:    Ptr("Close stale issues"),
		Version: Ptr(1),
		Active:  Ptr(true),
		Filter: map[string]any{
			"assignee": "user1",
		},
		Actions: []*AutomationAction{
			{
				Type:   Ptr("Transition"),
				Status: &Status{Key: Ptr("closed")},
			},
		},
		IntervalMillis: Ptr(int64(3600000)),
	}

	if !reflect.DeepEqual(autoAction, want) {
		t.Errorf("Queues.CreateAutoAction returned %+v, want %+v", autoAction, want)
	}
}

func TestQueuesService_UpdateAutoAction(t *testing.T) {
	client, mux := setup(t)

	input := &AutoActionUpdateRequest{
		Name:   Ptr("Updated auto-action"),
		Active: Ptr(false),
	}

	mux.HandleFunc("PATCH /v3/queues/{key}/autoactions/{id}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")

		// Verify request body.
		var got AutoActionUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("Decode request body: %v", err)
		}
		if !reflect.DeepEqual(&got, input) {
			t.Errorf("Request body = %+v, want %+v", got, input)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"id": "1",
			"self": "https://api.tracker.yandex.net/v3/queues/TEST/autoactions/1",
			"queue": {
				"self": "https://api.tracker.yandex.net/v2/queues/TEST",
				"id": "100",
				"key": "TEST",
				"display": "Test Queue"
			},
			"name": "Updated auto-action",
			"version": 4,
			"active": false,
			"actions": [
				{"type": "Transition", "id": "1", "status": {"self": "https://api.tracker.yandex.net/v2/statuses/closed", "key": "closed"}}
			],
			"intervalMillis": 3600000
		}`)
	})

	autoAction, _, err := client.Queues.UpdateAutoAction(ctx, "TEST", 1, input)
	if err != nil {
		t.Fatalf("Queues.UpdateAutoAction returned error: %v", err)
	}

	want := &AutoAction{
		ID:   Ptr("1"),
		Self: Ptr("https://api.tracker.yandex.net/v3/queues/TEST/autoactions/1"),
		Queue: &Queue{
			Self:    Ptr("https://api.tracker.yandex.net/v2/queues/TEST"),
			ID:      Ptr("100"),
			Key:     Ptr("TEST"),
			Display: Ptr("Test Queue"),
		},
		Name:    Ptr("Updated auto-action"),
		Version: Ptr(4),
		Active:  Ptr(false),
		Actions: []*AutomationAction{
			{
				Type:   Ptr("Transition"),
				ID:     Ptr("1"),
				Status: &Status{Self: Ptr("https://api.tracker.yandex.net/v2/statuses/closed"), Key: Ptr("closed")},
			},
		},
		IntervalMillis: Ptr(int64(3600000)),
	}

	if !reflect.DeepEqual(autoAction, want) {
		t.Errorf("Queues.UpdateAutoAction returned %+v, want %+v", autoAction, want)
	}
}
