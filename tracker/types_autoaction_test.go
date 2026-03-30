package tracker

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestAutoAction_JSONRoundTrip(t *testing.T) {
	want := &AutoAction{
		ID:   Ptr(FlexString("1")),
		Self: Ptr("https://api.tracker.yandex.net/v3/queues/TEST/autoactions/1"),
		Queue: &Queue{
			Self:    Ptr("https://api.tracker.yandex.net/v3/queues/TEST"),
			ID:      Ptr(FlexString("100")),
			Key:     Ptr("TEST"),
			Display: Ptr("Test Queue"),
		},
		Name:    Ptr("Auto close stale"),
		Version: Ptr(FlexString("3")),
		Active:  Ptr(true),
		Filter: map[string]any{
			"assignee": "user1",
		},
		Query: Ptr("status: open"),
		Actions: []*AutomationAction{
			{Type: Ptr("Transition"), Status: &Status{Key: Ptr("closed")}},
		},
		EnableNotifications:  Ptr(false),
		TotalIssuesProcessed: Ptr(42),
		IntervalMillis:       Ptr(int64(3600000)),
		Calendar:             &AutoActionCalendar{ID: Ptr(FlexString("1"))},
	}

	data, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}

	got := new(AutoAction)
	if err := json.Unmarshal(data, got); err != nil {
		t.Fatalf("json.Unmarshal error: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("AutoAction JSON round-trip:\ngot  %+v\nwant %+v", got, want)
	}
}

func TestAutoActionCalendar_JSONRoundTrip(t *testing.T) {
	want := &AutoActionCalendar{
		ID: Ptr(FlexString("1")),
	}

	data, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}

	got := new(AutoActionCalendar)
	if err := json.Unmarshal(data, got); err != nil {
		t.Fatalf("json.Unmarshal error: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("AutoActionCalendar JSON round-trip:\ngot  %+v\nwant %+v", got, want)
	}
}

func TestAutoActionCreateRequest_JSON(t *testing.T) {
	req := &AutoActionCreateRequest{
		Name: Ptr("Close stale issues"),
		Filter: map[string]any{
			"assignee": "user1",
		},
		Actions: []*AutomationAction{
			{Type: Ptr("Transition"), Status: &Status{Key: Ptr("closed")}},
		},
		Active:              Ptr(true),
		EnableNotifications: Ptr(false),
		IntervalMillis:      Ptr(int64(3600000)),
		Calendar:            &AutoActionCalendar{ID: Ptr(FlexString("1"))},
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("json.Unmarshal error: %v", err)
	}

	if m["name"] != "Close stale issues" {
		t.Errorf("name = %v, want %q", m["name"], "Close stale issues")
	}
	if m["intervalMillis"] != float64(3600000) {
		t.Errorf("intervalMillis = %v, want %v", m["intervalMillis"], 3600000)
	}
	filter, ok := m["filter"].(map[string]any)
	if !ok {
		t.Fatalf("filter is not a map, got %T", m["filter"])
	}
	if filter["assignee"] != "user1" {
		t.Errorf("filter.assignee = %v, want %q", filter["assignee"], "user1")
	}
}
