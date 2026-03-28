package tracker

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestFieldsService_ListLocal(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/queues/{queueKey}/localFields", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[
			{
				"self": "https://api.tracker.yandex.net/v3/queues/TEST/localFields/localField1",
				"id": "localField1",
				"name": "Local Field One",
				"key": "localField1",
				"version": 1,
				"schema": {"type": "string", "required": false},
				"readonly": false,
				"options": false,
				"suggest": false,
				"order": 10,
				"category": {"self": "https://api.tracker.yandex.net/v3/fields/categories/custom", "id": "custom", "display": "Custom"},
				"type": "local",
				"queue": {"self": "https://api.tracker.yandex.net/v2/queues/TEST", "id": "1", "key": "TEST", "display": "Test Queue"}
			},
			{
				"self": "https://api.tracker.yandex.net/v3/queues/TEST/localFields/localField2",
				"id": "localField2",
				"name": "Local Field Two",
				"key": "localField2",
				"version": 1,
				"schema": {"type": "integer", "required": false},
				"readonly": false,
				"options": false,
				"suggest": false,
				"order": 20,
				"category": {"self": "https://api.tracker.yandex.net/v3/fields/categories/custom", "id": "custom", "display": "Custom"},
				"type": "local",
				"queue": {"self": "https://api.tracker.yandex.net/v2/queues/TEST", "id": "1", "key": "TEST", "display": "Test Queue"}
			}
		]`)
	})

	ctx := context.Background()
	fields, _, err := client.Fields.ListLocal(ctx, "TEST")
	if err != nil {
		t.Fatalf("ListLocal returned error: %v", err)
	}

	if len(fields) != 2 {
		t.Fatalf("ListLocal returned %d fields, want 2", len(fields))
	}

	want := []*Field{
		{
			Self:     Ptr("https://api.tracker.yandex.net/v3/queues/TEST/localFields/localField1"),
			ID:       Ptr("localField1"),
			Name:     Ptr("Local Field One"),
			Key:      Ptr("localField1"),
			Version:  Ptr(1),
			Schema:   &FieldSchema{Type: Ptr("string"), Required: Ptr(false)},
			Readonly: Ptr(false),
			Options:  Ptr(false),
			Suggest:  Ptr(false),
			Order:    Ptr(10),
			Category: &FieldCategory{Self: Ptr("https://api.tracker.yandex.net/v3/fields/categories/custom"), ID: Ptr("custom"), Display: Ptr("Custom")},
			Type:     Ptr("local"),
			Queue:    &Queue{Self: Ptr("https://api.tracker.yandex.net/v2/queues/TEST"), ID: Ptr("1"), Key: Ptr("TEST"), Display: Ptr("Test Queue")},
		},
		{
			Self:     Ptr("https://api.tracker.yandex.net/v3/queues/TEST/localFields/localField2"),
			ID:       Ptr("localField2"),
			Name:     Ptr("Local Field Two"),
			Key:      Ptr("localField2"),
			Version:  Ptr(1),
			Schema:   &FieldSchema{Type: Ptr("integer"), Required: Ptr(false)},
			Readonly: Ptr(false),
			Options:  Ptr(false),
			Suggest:  Ptr(false),
			Order:    Ptr(20),
			Category: &FieldCategory{Self: Ptr("https://api.tracker.yandex.net/v3/fields/categories/custom"), ID: Ptr("custom"), Display: Ptr("Custom")},
			Type:     Ptr("local"),
			Queue:    &Queue{Self: Ptr("https://api.tracker.yandex.net/v2/queues/TEST"), ID: Ptr("1"), Key: Ptr("TEST"), Display: Ptr("Test Queue")},
		},
	}

	if !reflect.DeepEqual(fields, want) {
		t.Errorf("ListLocal returned %+v, want %+v", fields, want)
	}

	// Verify Queue field is populated on local fields
	if got := *fields[0].Queue.Key; got != "TEST" {
		t.Errorf("fields[0].Queue.Key = %q, want %q", got, "TEST")
	}
	if got := *fields[0].Type; got != "local" {
		t.Errorf("fields[0].Type = %q, want %q", got, "local")
	}
}

func TestFieldsService_GetLocal(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/queues/{queueKey}/localFields/{fieldKey}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/queues/TEST/localFields/localField1",
			"id": "localField1",
			"name": "Local Field One",
			"description": "A local field for testing",
			"key": "localField1",
			"version": 1,
			"schema": {"type": "string", "required": false},
			"readonly": false,
			"options": false,
			"suggest": false,
			"order": 10,
			"category": {"self": "https://api.tracker.yandex.net/v3/fields/categories/custom", "id": "custom", "display": "Custom"},
			"type": "local",
			"queue": {"self": "https://api.tracker.yandex.net/v2/queues/TEST", "id": "1", "key": "TEST", "display": "Test Queue"}
		}`)
	})

	ctx := context.Background()
	field, _, err := client.Fields.GetLocal(ctx, "TEST", "localField1")
	if err != nil {
		t.Fatalf("GetLocal returned error: %v", err)
	}

	want := &Field{
		Self:        Ptr("https://api.tracker.yandex.net/v3/queues/TEST/localFields/localField1"),
		ID:          Ptr("localField1"),
		Name:        Ptr("Local Field One"),
		Description: Ptr("A local field for testing"),
		Key:         Ptr("localField1"),
		Version:     Ptr(1),
		Schema:      &FieldSchema{Type: Ptr("string"), Required: Ptr(false)},
		Readonly:    Ptr(false),
		Options:     Ptr(false),
		Suggest:     Ptr(false),
		Order:       Ptr(10),
		Category:    &FieldCategory{Self: Ptr("https://api.tracker.yandex.net/v3/fields/categories/custom"), ID: Ptr("custom"), Display: Ptr("Custom")},
		Type:        Ptr("local"),
		Queue:       &Queue{Self: Ptr("https://api.tracker.yandex.net/v2/queues/TEST"), ID: Ptr("1"), Key: Ptr("TEST"), Display: Ptr("Test Queue")},
	}

	if !reflect.DeepEqual(field, want) {
		t.Errorf("GetLocal returned %+v, want %+v", field, want)
	}
}

func TestFieldsService_CreateLocal(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/queues/{queueKey}/localFields", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"id":"localNew","name":{"en":"New Local","ru":"Новое локальное"},"category":"0000000000000002********","type":"ru.yandex.startrek.core.fields.StringFieldType"}`)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/queues/TEST/localFields/localNew",
			"id": "localNew",
			"name": "New Local",
			"key": "localNew",
			"version": 1,
			"schema": {"type": "string", "required": false},
			"readonly": false,
			"options": false,
			"suggest": false,
			"order": 50,
			"category": {"self": "https://api.tracker.yandex.net/v3/fields/categories/custom", "id": "custom", "display": "Custom"},
			"type": "local",
			"queue": {"self": "https://api.tracker.yandex.net/v2/queues/TEST", "id": "1", "key": "TEST", "display": "Test Queue"}
		}`)
	})

	ctx := context.Background()
	field, resp, err := client.Fields.CreateLocal(ctx, "TEST", &FieldCreateRequest{
		ID:       Ptr("localNew"),
		Name:     &FieldName{EN: Ptr("New Local"), RU: Ptr("Новое локальное")},
		Category: Ptr("0000000000000002********"),
		Type:     Ptr("ru.yandex.startrek.core.fields.StringFieldType"),
	})
	if err != nil {
		t.Fatalf("CreateLocal returned error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusCreated)
	}
	if got := *field.ID; got != "localNew" {
		t.Errorf("ID = %q, want %q", got, "localNew")
	}
	if got := *field.Type; got != "local" {
		t.Errorf("Type = %q, want %q", got, "local")
	}
	if got := *field.Queue.Key; got != "TEST" {
		t.Errorf("Queue.Key = %q, want %q", got, "TEST")
	}
}

func TestFieldsService_EditLocal(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("PATCH /v3/queues/{queueKey}/localFields/{fieldKey}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		// CRITICAL: Verify NO version query parameter is present (unlike global Edit)
		if got := r.URL.Query().Get("version"); got != "" {
			t.Errorf("version query param = %q, want empty (local fields do not use version param)", got)
		}
		testBody(t, r, `{"name":{"en":"Updated Local","ru":"Обновлённое локальное"},"description":"Updated local field"}`)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/queues/TEST/localFields/localField1",
			"id": "localField1",
			"name": "Updated Local",
			"description": "Updated local field",
			"key": "localField1",
			"version": 2,
			"schema": {"type": "string", "required": false},
			"readonly": false,
			"options": false,
			"suggest": false,
			"order": 10,
			"category": {"self": "https://api.tracker.yandex.net/v3/fields/categories/custom", "id": "custom", "display": "Custom"},
			"type": "local",
			"queue": {"self": "https://api.tracker.yandex.net/v2/queues/TEST", "id": "1", "key": "TEST", "display": "Test Queue"}
		}`)
	})

	ctx := context.Background()
	field, _, err := client.Fields.EditLocal(ctx, "TEST", "localField1", &FieldEditRequest{
		Name:        &FieldName{EN: Ptr("Updated Local"), RU: Ptr("Обновлённое локальное")},
		Description: Ptr("Updated local field"),
	})
	if err != nil {
		t.Fatalf("EditLocal returned error: %v", err)
	}
	if got := *field.Name; got != "Updated Local" {
		t.Errorf("Name = %q, want %q", got, "Updated Local")
	}
	if got := *field.Description; got != "Updated local field" {
		t.Errorf("Description = %q, want %q", got, "Updated local field")
	}
	if got := *field.Version; got != 2 {
		t.Errorf("Version = %d, want %d", got, 2)
	}
	if got := *field.Queue.Key; got != "TEST" {
		t.Errorf("Queue.Key = %q, want %q", got, "TEST")
	}
}
