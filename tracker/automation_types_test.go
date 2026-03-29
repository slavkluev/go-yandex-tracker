package tracker

import (
	"encoding/json"
	"net/url"
	"reflect"
	"testing"
)

func TestTrigger_JSONRoundTrip(t *testing.T) {
	want := &Trigger{
		ID:   Ptr("1"),
		Self: Ptr("https://api.tracker.yandex.net/v3/queues/TEST/triggers/1"),
		Queue: &Queue{
			Self:    Ptr("https://api.tracker.yandex.net/v3/queues/TEST"),
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
				Status: &Status{Self: Ptr("https://api.tracker.yandex.net/v3/statuses/closed"), Key: Ptr("closed")},
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

	data, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("Marshal Trigger: %v", err)
	}

	got := new(Trigger)
	if err := json.Unmarshal(data, got); err != nil {
		t.Fatalf("Unmarshal Trigger: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Trigger round-trip:\ngot  %+v\nwant %+v", got, want)
	}
}

func TestAutomationAction_JSONRoundTrip(t *testing.T) {
	want := &AutomationAction{
		Type:   Ptr("Transition"),
		ID:     Ptr("1"),
		Status: &Status{Self: Ptr("https://api.tracker.yandex.net/v3/statuses/closed"), Key: Ptr("closed")},
		Parameters: map[string]any{
			"transitionScreen": "screen1",
		},
	}

	data, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("Marshal AutomationAction: %v", err)
	}

	got := new(AutomationAction)
	if err := json.Unmarshal(data, got); err != nil {
		t.Fatalf("Unmarshal AutomationAction: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("AutomationAction round-trip:\ngot  %+v\nwant %+v", got, want)
	}

	// Verify that the JSON output contains both known and overflow fields.
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Unmarshal raw: %v", err)
	}
	if raw["type"] != "Transition" {
		t.Errorf("expected type=Transition in JSON, got %v", raw["type"])
	}
	if raw["transitionScreen"] != "screen1" {
		t.Errorf("expected transitionScreen=screen1 in JSON, got %v", raw["transitionScreen"])
	}
}

func TestAutomationCondition_JSONRoundTrip(t *testing.T) {
	want := &AutomationCondition{
		Type: Ptr("Or"),
		Conditions: []*AutomationCondition{
			{
				Type: Ptr("Event.create"),
				Parameters: map[string]any{
					"field": "status",
				},
			},
		},
		Parameters: map[string]any{
			"ignoreCase": true,
		},
	}

	data, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("Marshal AutomationCondition: %v", err)
	}

	got := new(AutomationCondition)
	if err := json.Unmarshal(data, got); err != nil {
		t.Fatalf("Unmarshal AutomationCondition: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("AutomationCondition round-trip:\ngot  %+v\nwant %+v", got, want)
	}

	// Verify the JSON output contains known and overflow fields.
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Unmarshal raw: %v", err)
	}
	if raw["type"] != "Or" {
		t.Errorf("expected type=Or in JSON, got %v", raw["type"])
	}
	if raw["ignoreCase"] != true {
		t.Errorf("expected ignoreCase=true in JSON, got %v", raw["ignoreCase"])
	}
}

func TestTriggerCreateRequest_JSON(t *testing.T) {
	req := &TriggerCreateRequest{
		Name: Ptr("Test trigger"),
		Actions: []*AutomationAction{
			{Type: Ptr("Transition")},
		},
		Conditions: []*AutomationCondition{
			{Type: Ptr("Event.create")},
		},
		Active: Ptr(true),
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Marshal TriggerCreateRequest: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Unmarshal raw: %v", err)
	}

	if raw["name"] != "Test trigger" {
		t.Errorf("expected name=Test trigger, got %v", raw["name"])
	}
	if raw["active"] != true {
		t.Errorf("expected active=true, got %v", raw["active"])
	}
}

func TestTriggerUpdateOptions_URL(t *testing.T) {
	opts := &TriggerUpdateOptions{Version: 5}
	got, err := addOptions("v3/queues/TEST/triggers/1", opts)
	if err != nil {
		t.Fatalf("addOptions: %v", err)
	}

	u, err := url.Parse(got)
	if err != nil {
		t.Fatalf("Parse URL: %v", err)
	}

	if v := u.Query().Get("version"); v != "5" {
		t.Errorf("version query param = %q, want %q", v, "5")
	}
}
