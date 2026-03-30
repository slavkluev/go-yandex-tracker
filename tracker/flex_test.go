package tracker

import (
	"encoding/json"
	"testing"
)

func TestFlexString_UnmarshalJSON(t *testing.T) {
	t.Run("from string", func(t *testing.T) {
		var f FlexString
		if err := json.Unmarshal([]byte(`"abc"`), &f); err != nil {
			t.Fatalf("UnmarshalJSON returned error: %v", err)
		}
		if f != "abc" {
			t.Errorf("FlexString = %q, want %q", f, "abc")
		}
	})

	t.Run("from numeric string", func(t *testing.T) {
		var f FlexString
		if err := json.Unmarshal([]byte(`"42"`), &f); err != nil {
			t.Fatalf("UnmarshalJSON returned error: %v", err)
		}
		if f != "42" {
			t.Errorf("FlexString = %q, want %q", f, "42")
		}
	})

	t.Run("from number", func(t *testing.T) {
		var f FlexString
		if err := json.Unmarshal([]byte(`42`), &f); err != nil {
			t.Fatalf("UnmarshalJSON returned error: %v", err)
		}
		if f != "42" {
			t.Errorf("FlexString = %q, want %q", f, "42")
		}
	})

	t.Run("from zero number", func(t *testing.T) {
		var f FlexString
		if err := json.Unmarshal([]byte(`0`), &f); err != nil {
			t.Fatalf("UnmarshalJSON returned error: %v", err)
		}
		if f != "0" {
			t.Errorf("FlexString = %q, want %q", f, "0")
		}
	})

	t.Run("error on boolean", func(t *testing.T) {
		var f FlexString
		if err := json.Unmarshal([]byte(`true`), &f); err == nil {
			t.Error("expected error for boolean, got nil")
		}
	})

	t.Run("null unmarshals to empty", func(t *testing.T) {
		var f FlexString
		if err := json.Unmarshal([]byte(`null`), &f); err != nil {
			t.Fatalf("UnmarshalJSON returned error: %v", err)
		}
		if f != "" {
			t.Errorf("FlexString = %q, want empty", f)
		}
	})
}

func TestFlexString_MarshalJSON(t *testing.T) {
	t.Run("marshal string value", func(t *testing.T) {
		f := FlexString("42")
		data, err := json.Marshal(f)
		if err != nil {
			t.Fatalf("MarshalJSON returned error: %v", err)
		}
		if string(data) != `"42"` {
			t.Errorf("MarshalJSON = %s, want %s", string(data), `"42"`)
		}
	})

	t.Run("marshal non-numeric", func(t *testing.T) {
		f := FlexString("abc")
		data, err := json.Marshal(f)
		if err != nil {
			t.Fatalf("MarshalJSON returned error: %v", err)
		}
		if string(data) != `"abc"` {
			t.Errorf("MarshalJSON = %s, want %s", string(data), `"abc"`)
		}
	})
}

func TestFlexString_RoundTrip(t *testing.T) {
	original := FlexString("123")
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var decoded FlexString
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if decoded != original {
		t.Errorf("round-trip: got %q, want %q", decoded, original)
	}
}

func TestFlexString_PointerOmitempty(t *testing.T) {
	type S struct {
		ID *FlexString `json:"id,omitempty"`
	}

	t.Run("nil pointer omitted", func(t *testing.T) {
		data, _ := json.Marshal(S{})
		if string(data) != `{}` {
			t.Errorf("got %s, want {}", string(data))
		}
	})

	t.Run("non-nil pointer included", func(t *testing.T) {
		v := FlexString("5")
		data, _ := json.Marshal(S{ID: &v})
		if string(data) != `{"id":"5"}` {
			t.Errorf("got %s, want %s", string(data), `{"id":"5"}`)
		}
	})
}

// TestComponent_UnmarshalJSON_StringID tests deserialization of a component
// with a string ID, as returned when embedded in issue/import responses.
// JSON from API docs: https://yandex.ru/support/tracker/en/api-ref/import/import-ticket.
func TestComponent_UnmarshalJSON_StringID(t *testing.T) {
	data := []byte(`{
		"self": "https://api.tracker.yandex.net/v3/components/7",
		"id": "7",
		"display": "Component 7"
	}`)

	var c Component
	if err := json.Unmarshal(data, &c); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	if c.ID == nil || *c.ID != "7" {
		t.Errorf("Component.ID = %v, want %q", c.ID, "7")
	}
}

// TestComponent_UnmarshalJSON_NumberID tests deserialization of a component
// with a numeric ID, as returned by GET /v3/components.
// JSON from API docs: https://yandex.ru/support/tracker/en/api-ref/queues/get-components.
func TestComponent_UnmarshalJSON_NumberID(t *testing.T) {
	data := []byte(`{
		"self": "https://api.tracker.yandex.net/v3/components/1",
		"id": 1,
		"version": 3,
		"name": "Test",
		"queue": {
			"self": "https://api.tracker.yandex.net/v3/queues/ORG",
			"id": "1",
			"key": "ORG",
			"display": "Organization"
		},
		"description": "<component_description>",
		"lead": {
			"self": "https://api.tracker.yandex.net/v3/users/11********",
			"id": "11********",
			"display": "Ivan Ivanov"
		},
		"assignAuto": false
	}`)

	var c Component
	if err := json.Unmarshal(data, &c); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	if c.ID == nil || *c.ID != "1" {
		t.Errorf("Component.ID = %v, want %q", c.ID, "1")
	}
}

// TestWorklog_UnmarshalJSON_NumberID tests deserialization of a worklog
// with a numeric ID, as returned by GET /v3/issues/{id}/worklog.
// JSON from API docs: https://yandex.ru/support/tracker/en/api-ref/issues/issue-worklog.
func TestWorklog_UnmarshalJSON_NumberID(t *testing.T) {
	data := []byte(`{
		"self": "https://api.tracker.yandex.net/v3/issues/TEST-324/worklog/1",
		"id": 1,
		"version": "1402121720882",
		"issue": {
			"self": "https://api.tracker.yandex.net/v3/issues/TEST-324",
			"id": "515ec9eae4b09cfa********",
			"key": "TEST-324",
			"display": "important issue"
		},
		"comment": "important comment",
		"createdBy": {
			"self": "https://api.tracker.yandex.net/v3/users/66********",
			"id": "veikus",
			"display": "Artem Veikus"
		},
		"updatedBy": {
			"self": "https://api.tracker.yandex.net/v3/users/66********",
			"id": "veikus",
			"display": "Artem Veikus"
		},
		"createdAt": "2021-09-28T08:42:06.258+0000",
		"updatedAt": "2021-09-28T08:42:06.258+0000",
		"start": "2021-09-21T10:30:00.000+0000",
		"duration": "P3W"
	}`)

	var w Worklog
	if err := json.Unmarshal(data, &w); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	if w.ID == nil || *w.ID != "1" {
		t.Errorf("Worklog.ID = %v, want %q", w.ID, "1")
	}
}
