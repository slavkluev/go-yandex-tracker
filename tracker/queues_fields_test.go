package tracker

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestQueuesService_ListFields(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/queues/{key}/fields", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		if got := r.PathValue("key"); got != "TEST" {
			t.Errorf("queue key = %q, want %q", got, "TEST")
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[
			{
				"self": "https://api.tracker.yandex.net/v3/fields/stand",
				"id": "stand",
				"name": "Bench",
				"version": 1361890459119,
				"schema": {
					"type": "string",
					"required": false
				},
				"readonly": false,
				"options": true,
				"suggest": false,
				"optionsProvider": {
					"type": "QueueFixedListOptionsProvider",
					"values": {
						"DIRECT": [
							"Not specified",
							"Test",
							"Developer",
							"Beta",
							"Production",
							"Trunk"
						]
					},
					"defaults": [
						"Not specified",
						"Test",
						"Developer",
						"Beta",
						"Production"
					]
				},
				"queryProvider": {
					"type": "StringOptionalQueryProvider"
				},
				"order": 222
			},
			{"self": "https://api.tracker.yandex.net/v3/fields/type", "id": "type", "name": "Тип", "key": "type", "version": 0,
			 "schema": {"type": "issuetype", "required": true}, "readonly": false, "options": true, "suggest": true,
			 "suggestProvider": {"type": "IssueTypeSuggestProvider"}, "optionsProvider": {"type": "IssueTypeOptionsProvider"},
			 "queryProvider": {"type": "IssueTypeQueryProvider"}, "order": 2,
			 "category": {"self": "https://api.tracker.yandex.net/v3/fields/categories/000000000000000000000001", "id": "000000000000000000000001", "display": "Системные"},
			 "type": "standard"},
			{"self": "https://api.tracker.yandex.net/v3/fields/followers", "id": "followers", "name": "Наблюдатели", "key": "followers", "version": 0,
			 "schema": {"type": "array", "items": "user"}, "readonly": false, "options": true, "suggest": true,
			 "suggestProvider": {"type": "UserSuggestProvider"}, "optionsProvider": {"type": "TeamOptionsProvider"},
			 "queryProvider": {"type": "UserOptionalQueryProvider"}, "order": 22,
			 "category": {"self": "https://api.tracker.yandex.net/v3/fields/categories/000000000000000000000001", "id": "000000000000000000000001", "display": "Системные"},
			 "type": "standard"}
		]`)
	})

	fields, _, err := client.Queues.ListFields(ctx, "TEST")
	if err != nil {
		t.Fatalf("Queues.ListFields returned error: %v", err)
	}

	systemCategory := &FieldCategory{
		Self:    Ptr("https://api.tracker.yandex.net/v3/fields/categories/000000000000000000000001"),
		ID:      Ptr(FlexString("000000000000000000000001")),
		Display: Ptr("Системные"),
	}
	want := []*Field{
		{
			Self:     Ptr("https://api.tracker.yandex.net/v3/fields/stand"),
			ID:       Ptr(FlexString("stand")),
			Name:     Ptr("Bench"),
			Version:  Ptr(FlexString("1361890459119")),
			Schema:   &FieldSchema{Type: Ptr("string"), Required: Ptr(false)},
			Readonly: Ptr(false),
			Options:  Ptr(true),
			Suggest:  Ptr(false),
			OptionsProvider: &OptionsProvider{
				Type: Ptr("QueueFixedListOptionsProvider"),
				QueueValues: map[string][]any{
					"DIRECT": {"Not specified", "Test", "Developer", "Beta", "Production", "Trunk"},
				},
				Defaults: []any{"Not specified", "Test", "Developer", "Beta", "Production"},
			},
			QueryProvider: &QueryProvider{Type: Ptr("StringOptionalQueryProvider")},
			Order:         Ptr(222),
		},
		{
			Self:            Ptr("https://api.tracker.yandex.net/v3/fields/type"),
			ID:              Ptr(FlexString("type")),
			Name:            Ptr("Тип"),
			Key:             Ptr("type"),
			Version:         Ptr(FlexString("0")),
			Schema:          &FieldSchema{Type: Ptr("issuetype"), Required: Ptr(true)},
			Readonly:        Ptr(false),
			Options:         Ptr(true),
			Suggest:         Ptr(true),
			OptionsProvider: &OptionsProvider{Type: Ptr("IssueTypeOptionsProvider")},
			QueryProvider:   &QueryProvider{Type: Ptr("IssueTypeQueryProvider")},
			Order:           Ptr(2),
			Category:        systemCategory,
			Type:            Ptr("standard"),
		},
		{
			Self:            Ptr("https://api.tracker.yandex.net/v3/fields/followers"),
			ID:              Ptr(FlexString("followers")),
			Name:            Ptr("Наблюдатели"),
			Key:             Ptr("followers"),
			Version:         Ptr(FlexString("0")),
			Schema:          &FieldSchema{Type: Ptr("array"), Items: Ptr("user")},
			Readonly:        Ptr(false),
			Options:         Ptr(true),
			Suggest:         Ptr(true),
			OptionsProvider: &OptionsProvider{Type: Ptr("TeamOptionsProvider")},
			QueryProvider:   &QueryProvider{Type: Ptr("UserOptionalQueryProvider")},
			Order:           Ptr(22),
			Category:        systemCategory,
			Type:            Ptr("standard"),
		},
	}

	if !reflect.DeepEqual(fields, want) {
		t.Errorf("Queues.ListFields returned %+v, want %+v", fields, want)
	}
}

func TestQueuesService_ListFields_Empty(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/queues/{key}/fields", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[]`)
	})

	fields, _, err := client.Queues.ListFields(ctx, "RPA")
	if err != nil {
		t.Fatalf("Queues.ListFields returned error: %v", err)
	}

	if want := []*Field{}; !reflect.DeepEqual(fields, want) {
		t.Errorf("Queues.ListFields returned %#v, want %#v", fields, want)
	}
}

func TestQueuesService_ListFields_NotFound(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/queues/{key}/fields", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"errorMessages":["Object not found"],"errors":{}}`)
	})

	_, _, err := client.Queues.ListFields(ctx, "NONEXIST")
	if err == nil {
		t.Fatal("Queues.ListFields expected error for 404, got nil")
	}
	if !IsNotFound(err) {
		t.Errorf("IsNotFound(err) = false, want true")
	}
}
