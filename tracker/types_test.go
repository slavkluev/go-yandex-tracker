package tracker

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestTimestamp_UnmarshalJSON(t *testing.T) {
	t.Run("yandex format", func(t *testing.T) {
		var ts Timestamp
		data := []byte(`"2017-06-11T05:16:01.339+0000"`)
		if err := json.Unmarshal(data, &ts); err != nil {
			t.Fatalf("UnmarshalJSON returned error: %v", err)
		}
		if ts.Year() != 2017 {
			t.Errorf("Year = %d, want 2017", ts.Year())
		}
		if ts.Month() != time.June {
			t.Errorf("Month = %v, want June", ts.Month())
		}
		if ts.Day() != 11 {
			t.Errorf("Day = %d, want 11", ts.Day())
		}
		if ts.Hour() != 5 {
			t.Errorf("Hour = %d, want 5", ts.Hour())
		}
		if ts.Minute() != 16 {
			t.Errorf("Minute = %d, want 16", ts.Minute())
		}
		if ts.Second() != 1 {
			t.Errorf("Second = %d, want 1", ts.Second())
		}
	})

	t.Run("RFC3339 fallback", func(t *testing.T) {
		var ts Timestamp
		data := []byte(`"2017-06-11T05:16:01Z"`)
		if err := json.Unmarshal(data, &ts); err != nil {
			t.Fatalf("UnmarshalJSON returned error: %v", err)
		}
		if ts.Year() != 2017 {
			t.Errorf("Year = %d, want 2017", ts.Year())
		}
		if ts.Month() != time.June {
			t.Errorf("Month = %v, want June", ts.Month())
		}
		if ts.Day() != 11 {
			t.Errorf("Day = %d, want 11", ts.Day())
		}
	})

	t.Run("null value", func(t *testing.T) {
		var ts Timestamp
		data := []byte(`null`)
		if err := json.Unmarshal(data, &ts); err != nil {
			t.Fatalf("UnmarshalJSON returned error: %v", err)
		}
		if !ts.IsZero() {
			t.Errorf("expected zero time, got %v", ts.Time)
		}
	})

	t.Run("empty string", func(t *testing.T) {
		var ts Timestamp
		data := []byte(`""`)
		if err := json.Unmarshal(data, &ts); err != nil {
			t.Fatalf("UnmarshalJSON returned error: %v", err)
		}
		if !ts.IsZero() {
			t.Errorf("expected zero time, got %v", ts.Time)
		}
	})
}

func TestTimestamp_MarshalJSON(t *testing.T) {
	ts := Timestamp{Time: time.Date(2017, 6, 11, 5, 16, 1, 339000000, time.UTC)}
	data, err := json.Marshal(ts)
	if err != nil {
		t.Fatalf("MarshalJSON returned error: %v", err)
	}
	want := `"2017-06-11T05:16:01.339+0000"`
	if string(data) != want {
		t.Errorf("MarshalJSON = %s, want %s", string(data), want)
	}
}

func TestTimestamp_MarshalJSON_RoundTrip(t *testing.T) {
	original := `"2017-06-11T05:16:01.339+0000"`
	var ts Timestamp
	if err := json.Unmarshal([]byte(original), &ts); err != nil {
		t.Fatalf("UnmarshalJSON returned error: %v", err)
	}
	data, err := json.Marshal(ts)
	if err != nil {
		t.Fatalf("MarshalJSON returned error: %v", err)
	}
	if string(data) != original {
		t.Errorf("round-trip failed: got %s, want %s", string(data), original)
	}
}

func TestIssue_UnmarshalJSON(t *testing.T) {
	t.Run("known fields", func(t *testing.T) {
		data := []byte(`{
			"self": "https://api.tracker.yandex.net/v2/issues/QUEUE-1",
			"id": "123",
			"key": "QUEUE-1",
			"summary": "Test issue",
			"description": "A description",
			"status": {
				"self": "https://api.tracker.yandex.net/v2/statuses/1",
				"id": "1",
				"key": "open",
				"display": "Open"
			},
			"priority": {
				"self": "https://api.tracker.yandex.net/v2/priorities/1",
				"id": "1",
				"key": "normal",
				"display": "Normal"
			},
			"type": {
				"self": "https://api.tracker.yandex.net/v2/issuetypes/1",
				"id": "1",
				"key": "task",
				"display": "Task"
			},
			"queue": {
				"self": "https://api.tracker.yandex.net/v2/queues/QUEUE",
				"id": "1",
				"key": "QUEUE",
				"display": "Queue"
			},
			"assignee": {
				"self": "https://api.tracker.yandex.net/v2/users/1",
				"id": "1",
				"display": "User 1"
			}
		}`)

		var issue Issue
		if err := json.Unmarshal(data, &issue); err != nil {
			t.Fatalf("UnmarshalJSON returned error: %v", err)
		}

		if got := *issue.Key; got != "QUEUE-1" {
			t.Errorf("Key = %q, want %q", got, "QUEUE-1")
		}
		if got := *issue.Summary; got != "Test issue" {
			t.Errorf("Summary = %q, want %q", got, "Test issue")
		}
		if got := *issue.Description; got != "A description" {
			t.Errorf("Description = %q, want %q", got, "A description")
		}
		if got := *issue.Status.Key; got != "open" {
			t.Errorf("Status.Key = %q, want %q", got, "open")
		}
		if got := *issue.Priority.Key; got != "normal" {
			t.Errorf("Priority.Key = %q, want %q", got, "normal")
		}
		if got := *issue.Type.Key; got != "task" {
			t.Errorf("Type.Key = %q, want %q", got, "task")
		}
		if got := *issue.Queue.Key; got != "QUEUE" {
			t.Errorf("Queue.Key = %q, want %q", got, "QUEUE")
		}
		if got := *issue.Assignee.Display; got != "User 1" {
			t.Errorf("Assignee.Display = %q, want %q", got, "User 1")
		}
		if issue.CustomFields != nil {
			t.Errorf("CustomFields = %v, want nil", issue.CustomFields)
		}
	})

	t.Run("custom fields", func(t *testing.T) {
		data := []byte(`{
			"self": "https://api.tracker.yandex.net/v2/issues/QUEUE-1",
			"id": "123",
			"key": "QUEUE-1",
			"summary": "Test issue",
			"customField1": "value1",
			"customField2": 42
		}`)

		var issue Issue
		if err := json.Unmarshal(data, &issue); err != nil {
			t.Fatalf("UnmarshalJSON returned error: %v", err)
		}

		if got := *issue.Key; got != "QUEUE-1" {
			t.Errorf("Key = %q, want %q", got, "QUEUE-1")
		}
		if issue.CustomFields == nil {
			t.Fatal("CustomFields is nil, expected non-nil map")
		}
		if got, ok := issue.CustomFields["customField1"]; !ok || got != "value1" {
			t.Errorf("CustomFields[customField1] = %v, want %q", got, "value1")
		}
		if got, ok := issue.CustomFields["customField2"]; !ok {
			t.Error("CustomFields[customField2] not found")
		} else {
			// JSON numbers decode as float64
			if got != float64(42) {
				t.Errorf("CustomFields[customField2] = %v (%T), want 42", got, got)
			}
		}
	})

	t.Run("no custom fields", func(t *testing.T) {
		data := []byte(`{
			"self": "https://api.tracker.yandex.net/v2/issues/QUEUE-1",
			"id": "123",
			"key": "QUEUE-1",
			"summary": "Test issue"
		}`)

		var issue Issue
		if err := json.Unmarshal(data, &issue); err != nil {
			t.Fatalf("UnmarshalJSON returned error: %v", err)
		}

		if issue.CustomFields != nil {
			t.Errorf("CustomFields = %v, want nil", issue.CustomFields)
		}
	})
}

func TestPtr(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		p := Ptr("hello")
		if p == nil {
			t.Fatal("Ptr returned nil")
		}
		if *p != "hello" {
			t.Errorf("*Ptr = %q, want %q", *p, "hello")
		}
	})

	t.Run("int", func(t *testing.T) {
		p := Ptr(42)
		if p == nil {
			t.Fatal("Ptr returned nil")
		}
		if *p != 42 {
			t.Errorf("*Ptr = %d, want 42", *p)
		}
	})

	t.Run("bool", func(t *testing.T) {
		p := Ptr(true)
		if p == nil {
			t.Fatal("Ptr returned nil")
		}
		if *p != true {
			t.Errorf("*Ptr = %v, want true", *p)
		}
	})
}

func TestNewResponse(t *testing.T) {
	t.Run("all headers", func(t *testing.T) {
		httpResp := &http.Response{
			StatusCode: 200,
			Header:     http.Header{},
		}
		httpResp.Header.Set("X-Total-Count", "42")
		httpResp.Header.Set("X-Total-Pages", "5")
		httpResp.Header.Set("X-Scroll-Id", "scroll-id-123")
		httpResp.Header.Set("X-Scroll-Token", "scroll-token-456")
		httpResp.Header.Set("Retry-After", "30")

		resp := newResponse(httpResp)

		if resp.TotalCount != 42 {
			t.Errorf("TotalCount = %d, want 42", resp.TotalCount)
		}
		if resp.TotalPages != 5 {
			t.Errorf("TotalPages = %d, want 5", resp.TotalPages)
		}
		if resp.ScrollID != "scroll-id-123" {
			t.Errorf("ScrollID = %q, want %q", resp.ScrollID, "scroll-id-123")
		}
		if resp.ScrollToken != "scroll-token-456" {
			t.Errorf("ScrollToken = %q, want %q", resp.ScrollToken, "scroll-token-456")
		}
		if resp.RetryAfter != 30 {
			t.Errorf("RetryAfter = %d, want 30", resp.RetryAfter)
		}
	})

	t.Run("missing headers", func(t *testing.T) {
		httpResp := &http.Response{
			StatusCode: 200,
			Header:     http.Header{},
		}

		resp := newResponse(httpResp)

		if resp.TotalCount != 0 {
			t.Errorf("TotalCount = %d, want 0", resp.TotalCount)
		}
		if resp.TotalPages != 0 {
			t.Errorf("TotalPages = %d, want 0", resp.TotalPages)
		}
		if resp.ScrollID != "" {
			t.Errorf("ScrollID = %q, want empty", resp.ScrollID)
		}
		if resp.ScrollToken != "" {
			t.Errorf("ScrollToken = %q, want empty", resp.ScrollToken)
		}
		if resp.RetryAfter != 0 {
			t.Errorf("RetryAfter = %d, want 0", resp.RetryAfter)
		}
	})
}

func TestNewResponse_RetryAfter(t *testing.T) {
	t.Run("with Retry-After header", func(t *testing.T) {
		httpResp := &http.Response{
			StatusCode: 429,
			Header:     http.Header{},
		}
		httpResp.Header.Set("Retry-After", "30")

		resp := newResponse(httpResp)

		if resp.RetryAfter != 30 {
			t.Errorf("RetryAfter = %d, want 30", resp.RetryAfter)
		}
	})

	t.Run("without Retry-After header", func(t *testing.T) {
		httpResp := &http.Response{
			StatusCode: 200,
			Header:     http.Header{},
		}

		resp := newResponse(httpResp)

		if resp.RetryAfter != 0 {
			t.Errorf("RetryAfter = %d, want 0", resp.RetryAfter)
		}
	})
}

func TestDuration_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Duration
		wantErr bool
	}{
		{name: "hours and minutes", input: `"PT1H30M"`, want: 1*time.Hour + 30*time.Minute},
		{name: "weeks", input: `"P1W"`, want: 7 * 24 * time.Hour},
		{name: "days", input: `"P5D"`, want: 5 * 24 * time.Hour},
		{name: "days hours minutes", input: `"P1DT2H3M"`, want: 24*time.Hour + 2*time.Hour + 3*time.Minute},
		{name: "minutes only", input: `"PT300M"`, want: 5 * time.Hour},
		{name: "seconds only", input: `"PT30S"`, want: 30 * time.Second},
		{name: "full format", input: `"P0Y0M30DT2H10M25S"`, want: 30*24*time.Hour + 2*time.Hour + 10*time.Minute + 25*time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d Duration
			err := json.Unmarshal([]byte(tt.input), &d)
			if err != nil {
				t.Fatalf("UnmarshalJSON(%s) returned error: %v", tt.input, err)
			}
			if d.Duration != tt.want {
				t.Errorf("UnmarshalJSON(%s) = %v, want %v", tt.input, d.Duration, tt.want)
			}
		})
	}
}

func TestDuration_UnmarshalJSON_Invalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "invalid string", input: `"invalid"`},
		{name: "empty string", input: `""`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d Duration
			err := json.Unmarshal([]byte(tt.input), &d)
			if err == nil {
				t.Errorf("UnmarshalJSON(%s) expected error, got nil", tt.input)
			}
		})
	}
}

func TestDuration_MarshalJSON(t *testing.T) {
	tests := []struct {
		name  string
		input time.Duration
		want  string
	}{
		{name: "hours and minutes", input: 1*time.Hour + 30*time.Minute, want: `"PT1H30M"`},
		{name: "weeks", input: 7 * 24 * time.Hour, want: `"P1W"`},
		{name: "days", input: 5 * 24 * time.Hour, want: `"P5D"`},
		{name: "seconds only", input: 30 * time.Second, want: `"PT30S"`},
		{name: "days hours minutes", input: 24*time.Hour + 2*time.Hour + 3*time.Minute, want: `"P1DT2H3M"`},
		{name: "zero", input: 0, want: `"PT0S"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Duration{Duration: tt.input}
			data, err := json.Marshal(d)
			if err != nil {
				t.Fatalf("MarshalJSON returned error: %v", err)
			}
			if got := string(data); got != tt.want {
				t.Errorf("MarshalJSON(%v) = %s, want %s", tt.input, got, tt.want)
			}
		})
	}
}

func TestDuration_RoundTrip(t *testing.T) {
	durations := []time.Duration{
		1*time.Hour + 30*time.Minute,
		7 * 24 * time.Hour,
		5 * 24 * time.Hour,
		30 * time.Second,
		24*time.Hour + 2*time.Hour + 3*time.Minute,
	}

	for _, original := range durations {
		t.Run(original.String(), func(t *testing.T) {
			d := Duration{Duration: original}
			data, err := json.Marshal(d)
			if err != nil {
				t.Fatalf("MarshalJSON returned error: %v", err)
			}
			var roundTripped Duration
			if err := json.Unmarshal(data, &roundTripped); err != nil {
				t.Fatalf("UnmarshalJSON returned error: %v", err)
			}
			if roundTripped.Duration != original {
				t.Errorf("round-trip failed: got %v, want %v", roundTripped.Duration, original)
			}
		})
	}
}

func TestMacro_JSON(t *testing.T) {
	data := []byte(`{
		"id": 5,
		"self": "https://api.tracker.yandex.net/v3/queues/TEST/macros/5",
		"queue": {
			"self": "https://api.tracker.yandex.net/v2/queues/TEST",
			"id": "100",
			"key": "TEST",
			"display": "Test Queue"
		},
		"name": "Close issue macro",
		"body": "Closing this issue per policy.",
		"issueUpdate": [
			{
				"field": {
					"self": "https://api.tracker.yandex.net/v2/fields/status",
					"id": "status",
					"display": "Status"
				},
				"update": {"set": "closed"}
			}
		]
	}`)

	var macro Macro
	if err := json.Unmarshal(data, &macro); err != nil {
		t.Fatalf("Unmarshal Macro: %v", err)
	}

	// Macro.ID is *int (not *string)
	if macro.ID == nil || *macro.ID != 5 {
		t.Errorf("Macro.ID = %v, want Ptr(5)", macro.ID)
	}
	if macro.Self == nil || *macro.Self != "https://api.tracker.yandex.net/v3/queues/TEST/macros/5" {
		t.Errorf("Macro.Self = %v, want non-nil", macro.Self)
	}
	if macro.Queue == nil || *macro.Queue.Key != "TEST" {
		t.Errorf("Macro.Queue.Key = %v, want TEST", macro.Queue)
	}
	if macro.Name == nil || *macro.Name != "Close issue macro" {
		t.Errorf("Macro.Name = %v, want Close issue macro", macro.Name)
	}
	if macro.Body == nil || *macro.Body != "Closing this issue per policy." {
		t.Errorf("Macro.Body = %v, want non-nil", macro.Body)
	}
	if len(macro.IssueUpdate) != 1 {
		t.Fatalf("Macro.IssueUpdate length = %d, want 1", len(macro.IssueUpdate))
	}

	// Round-trip: marshal back and compare
	out, err := json.Marshal(macro)
	if err != nil {
		t.Fatalf("Marshal Macro: %v", err)
	}
	var macro2 Macro
	if err := json.Unmarshal(out, &macro2); err != nil {
		t.Fatalf("Unmarshal round-trip: %v", err)
	}
	if !reflect.DeepEqual(macro, macro2) {
		t.Errorf("Macro round-trip mismatch:\n got: %+v\nwant: %+v", macro2, macro)
	}
}

func TestMacroIssueUpdate_JSON(t *testing.T) {
	data := []byte(`{
		"field": {
			"self": "https://api.tracker.yandex.net/v2/fields/assignee",
			"id": "assignee",
			"display": "Assignee"
		},
		"update": {"set": "user123"}
	}`)

	var update MacroIssueUpdate
	if err := json.Unmarshal(data, &update); err != nil {
		t.Fatalf("Unmarshal MacroIssueUpdate: %v", err)
	}
	if update.Field == nil {
		t.Fatal("MacroIssueUpdate.Field is nil")
	}
	if *update.Field.ID != "assignee" {
		t.Errorf("Field.ID = %q, want assignee", *update.Field.ID)
	}
	if update.Update == nil {
		t.Fatal("MacroIssueUpdate.Update is nil")
	}

	// Round-trip
	out, err := json.Marshal(update)
	if err != nil {
		t.Fatalf("Marshal MacroIssueUpdate: %v", err)
	}
	var update2 MacroIssueUpdate
	if err := json.Unmarshal(out, &update2); err != nil {
		t.Fatalf("Unmarshal round-trip: %v", err)
	}
	if !reflect.DeepEqual(update, update2) {
		t.Errorf("MacroIssueUpdate round-trip mismatch")
	}
}

func TestMacroIssueUpdateField_JSON(t *testing.T) {
	data := []byte(`{
		"self": "https://api.tracker.yandex.net/v2/fields/priority",
		"id": "priority",
		"display": "Priority"
	}`)

	var field MacroIssueUpdateField
	if err := json.Unmarshal(data, &field); err != nil {
		t.Fatalf("Unmarshal MacroIssueUpdateField: %v", err)
	}
	if field.Self == nil || *field.Self != "https://api.tracker.yandex.net/v2/fields/priority" {
		t.Errorf("Self = %v, want non-nil", field.Self)
	}
	if field.ID == nil || *field.ID != "priority" {
		t.Errorf("ID = %v, want priority", field.ID)
	}
	if field.Display == nil || *field.Display != "Priority" {
		t.Errorf("Display = %v, want Priority", field.Display)
	}

	// Round-trip
	out, err := json.Marshal(field)
	if err != nil {
		t.Fatalf("Marshal MacroIssueUpdateField: %v", err)
	}
	var field2 MacroIssueUpdateField
	if err := json.Unmarshal(out, &field2); err != nil {
		t.Fatalf("Unmarshal round-trip: %v", err)
	}
	if !reflect.DeepEqual(field, field2) {
		t.Errorf("MacroIssueUpdateField round-trip mismatch")
	}
}

func TestMacroCreateRequest_JSON(t *testing.T) {
	req := &MacroCreateRequest{
		Name: Ptr("Test macro"),
		Body: Ptr("Comment body"),
		IssueUpdate: []map[string]any{
			{"field": "status", "set": "closed"},
		},
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Marshal MacroCreateRequest: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal to map: %v", err)
	}
	if got["name"] != "Test macro" {
		t.Errorf("name = %v, want Test macro", got["name"])
	}
	if got["body"] != "Comment body" {
		t.Errorf("body = %v, want Comment body", got["body"])
	}
	if got["issueUpdate"] == nil {
		t.Error("issueUpdate is nil, want non-nil")
	}
}

// Suppress unused import warning - reflect is used for DeepEqual in other test files.
var _ = reflect.DeepEqual
