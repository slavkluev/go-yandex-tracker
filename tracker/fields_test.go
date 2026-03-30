package tracker

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestFieldsService_List(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/fields", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[
			{
				"self": "https://api.tracker.yandex.net/v3/fields/summary",
				"id": "summary",
				"name": "Summary",
				"key": "summary",
				"version": 1,
				"schema": {"type": "string", "required": true},
				"readonly": false,
				"options": false,
				"suggest": false,
				"order": 1,
				"category": {"self": "https://api.tracker.yandex.net/v3/fields/categories/system", "id": "system", "display": "System"}
			},
			{
				"self": "https://api.tracker.yandex.net/v3/fields/description",
				"id": "description",
				"name": "Description",
				"key": "description",
				"version": 1,
				"schema": {"type": "string", "required": false},
				"readonly": false,
				"options": false,
				"suggest": false,
				"order": 2,
				"category": {"self": "https://api.tracker.yandex.net/v3/fields/categories/system", "id": "system", "display": "System"}
			}
		]`)
	})

	ctx := context.Background()
	fields, _, err := client.Fields.List(ctx)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	if len(fields) != 2 {
		t.Fatalf("List returned %d fields, want 2", len(fields))
	}

	want := []*Field{
		{
			Self:     Ptr("https://api.tracker.yandex.net/v3/fields/summary"),
			ID:       Ptr(FlexString("summary")),
			Name:     Ptr("Summary"),
			Key:      Ptr("summary"),
			Version:  Ptr(FlexString("1")),
			Schema:   &FieldSchema{Type: Ptr("string"), Required: Ptr(true)},
			Readonly: Ptr(false),
			Options:  Ptr(false),
			Suggest:  Ptr(false),
			Order:    Ptr(1),
			Category: &FieldCategory{Self: Ptr("https://api.tracker.yandex.net/v3/fields/categories/system"), ID: Ptr(FlexString("system")), Display: Ptr("System")},
		},
		{
			Self:     Ptr("https://api.tracker.yandex.net/v3/fields/description"),
			ID:       Ptr(FlexString("description")),
			Name:     Ptr("Description"),
			Key:      Ptr("description"),
			Version:  Ptr(FlexString("1")),
			Schema:   &FieldSchema{Type: Ptr("string"), Required: Ptr(false)},
			Readonly: Ptr(false),
			Options:  Ptr(false),
			Suggest:  Ptr(false),
			Order:    Ptr(2),
			Category: &FieldCategory{Self: Ptr("https://api.tracker.yandex.net/v3/fields/categories/system"), ID: Ptr(FlexString("system")), Display: Ptr("System")},
		},
	}

	if !reflect.DeepEqual(fields, want) {
		t.Errorf("List returned %+v, want %+v", fields, want)
	}
}

func TestFieldsService_Get(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/fields/{id}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/fields/summary",
			"id": "summary",
			"name": "Summary",
			"description": "Issue summary field",
			"key": "summary",
			"version": 1,
			"schema": {"type": "string", "required": true},
			"readonly": false,
			"options": false,
			"suggest": false,
			"order": 1,
			"category": {"self": "https://api.tracker.yandex.net/v3/fields/categories/system", "id": "system", "display": "System"}
		}`)
	})

	ctx := context.Background()
	field, _, err := client.Fields.Get(ctx, "summary")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}

	want := &Field{
		Self:        Ptr("https://api.tracker.yandex.net/v3/fields/summary"),
		ID:          Ptr(FlexString("summary")),
		Name:        Ptr("Summary"),
		Description: Ptr("Issue summary field"),
		Key:         Ptr("summary"),
		Version:     Ptr(FlexString("1")),
		Schema:      &FieldSchema{Type: Ptr("string"), Required: Ptr(true)},
		Readonly:    Ptr(false),
		Options:     Ptr(false),
		Suggest:     Ptr(false),
		Order:       Ptr(1),
		Category:    &FieldCategory{Self: Ptr("https://api.tracker.yandex.net/v3/fields/categories/system"), ID: Ptr(FlexString("system")), Display: Ptr("System")},
	}

	if !reflect.DeepEqual(field, want) {
		t.Errorf("Get returned %+v, want %+v", field, want)
	}
}

func TestFieldsService_Create(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/fields", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"id":"myField","name":{"en":"Field name","ru":"Имя поля"},"category":"0000000000000002********","type":"ru.yandex.startrek.core.fields.StringFieldType"}`)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/fields/myField",
			"id": "myField",
			"name": "Field name",
			"key": "myField",
			"version": 1,
			"schema": {"type": "string", "required": false},
			"readonly": false,
			"options": false,
			"suggest": false,
			"order": 100,
			"category": {"self": "https://api.tracker.yandex.net/v3/fields/categories/custom", "id": "custom", "display": "Custom"}
		}`)
	})

	ctx := context.Background()
	field, resp, err := client.Fields.Create(ctx, &FieldCreateRequest{
		ID:       Ptr("myField"),
		Name:     &FieldName{EN: Ptr("Field name"), RU: Ptr("Имя поля")},
		Category: Ptr("0000000000000002********"),
		Type:     Ptr("ru.yandex.startrek.core.fields.StringFieldType"),
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusCreated)
	}
	if got := *field.ID; got != "myField" {
		t.Errorf("ID = %q, want %q", got, "myField")
	}
	if got := *field.Name; got != "Field name" {
		t.Errorf("Name = %q, want %q", got, "Field name")
	}
	if got := *field.Version; got != FlexString("1") {
		t.Errorf("Version = %q, want %q", got, FlexString("1"))
	}
}

func TestFieldsService_Edit(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("PATCH /v3/fields/{id}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		// Verify version query parameter is present
		if got := r.URL.Query().Get("version"); got != "2" {
			t.Errorf("version query param = %q, want %q", got, "2")
		}
		testBody(t, r, `{"name":{"en":"Updated name","ru":"Обновлённое имя"},"description":"Updated description"}`)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/fields/myField",
			"id": "myField",
			"name": "Updated name",
			"description": "Updated description",
			"key": "myField",
			"version": 3,
			"schema": {"type": "string", "required": false},
			"readonly": false,
			"options": false,
			"suggest": false,
			"order": 100,
			"category": {"self": "https://api.tracker.yandex.net/v3/fields/categories/custom", "id": "custom", "display": "Custom"}
		}`)
	})

	ctx := context.Background()
	field, _, err := client.Fields.Edit(ctx, "myField", &FieldEditRequest{
		Name:        &FieldName{EN: Ptr("Updated name"), RU: Ptr("Обновлённое имя")},
		Description: Ptr("Updated description"),
	}, &FieldEditOptions{Version: 2})
	if err != nil {
		t.Fatalf("Edit returned error: %v", err)
	}
	if got := *field.Name; got != "Updated name" {
		t.Errorf("Name = %q, want %q", got, "Updated name")
	}
	if got := *field.Description; got != "Updated description" {
		t.Errorf("Description = %q, want %q", got, "Updated description")
	}
	if got := *field.Version; got != FlexString("3") {
		t.Errorf("Version = %q, want %q", got, FlexString("3"))
	}
}

func TestFieldsService_List_Error(t *testing.T) {
	client, mux := setup(t)
	ctx := context.Background()

	mux.HandleFunc("GET /v3/fields", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, `{"errorMessages":["Access denied"],"errors":{}}`)
	})

	_, _, err := client.Fields.List(ctx)
	if err == nil {
		t.Fatal("Fields.List expected error for 403, got nil")
	}
	if !IsForbidden(err) {
		t.Errorf("IsForbidden(err) = false, want true")
	}
}
