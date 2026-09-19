package tracker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
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

// TestFieldsService_List_OptionValues decodes option values of both JSON types
// in one response: epicstatus lists strings, and possibleSpam lists numbers.
// Both elements are recorded verbatim from GET /v3/fields on a real org, where
// possibleSpam's numeric values once failed the whole list.
func TestFieldsService_List_OptionValues(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/fields", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[`+
			`{"self":"https://api.tracker.yandex.net/v3/fields/epicstatus","id":"epicstatus","name":"Epic Status","key":"epicstatus","version":1,"schema":{"type":"string","required":false},"readonly":false,"options":true,"suggest":false,"optionsProvider":{"type":"FixedListOptionsProvider","needValidation":true,"values":["To Do","In Progress","Done"]},"queryProvider":{"type":"StringOptionalQueryProvider"},"order":29,"category":{"self":"https://api.tracker.yandex.net/v3/fields/categories/000000000000000000000001","id":"000000000000000000000001","display":"Системные"},"type":"standard"},`+
			`{"self":"https://api.tracker.yandex.net/v3/fields/possibleSpam","id":"possibleSpam","name":"Возможно спам","key":"possibleSpam","version":0,"schema":{"type":"integer","required":false},"readonly":false,"options":true,"suggest":false,"optionsProvider":{"type":"FixedListOptionsProvider","needValidation":true,"values":[0,1]},"queryProvider":{"type":"NumberOptionalQueryProvider"},"order":66,"category":{"self":"https://api.tracker.yandex.net/v3/fields/categories/000000000000000000000001","id":"000000000000000000000001","display":"Системные"},"type":"standard"}`+
			`]`)
	})

	fields, _, err := client.Fields.List(context.Background())
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	category := &FieldCategory{
		Self:    Ptr("https://api.tracker.yandex.net/v3/fields/categories/000000000000000000000001"),
		ID:      Ptr(FlexString("000000000000000000000001")),
		Display: Ptr("Системные"),
	}
	want := []*Field{
		{
			Self:     Ptr("https://api.tracker.yandex.net/v3/fields/epicstatus"),
			ID:       Ptr(FlexString("epicstatus")),
			Name:     Ptr("Epic Status"),
			Key:      Ptr("epicstatus"),
			Version:  Ptr(FlexString("1")),
			Schema:   &FieldSchema{Type: Ptr("string"), Required: Ptr(false)},
			Readonly: Ptr(false),
			Options:  Ptr(true),
			Suggest:  Ptr(false),
			OptionsProvider: &OptionsProvider{
				Type:           Ptr("FixedListOptionsProvider"),
				NeedValidation: Ptr(true),
				Values:         []any{"To Do", "In Progress", "Done"},
			},
			QueryProvider: &QueryProvider{Type: Ptr("StringOptionalQueryProvider")},
			Order:         Ptr(29),
			Category:      category,
			Type:          Ptr("standard"),
		},
		{
			Self:     Ptr("https://api.tracker.yandex.net/v3/fields/possibleSpam"),
			ID:       Ptr(FlexString("possibleSpam")),
			Name:     Ptr("Возможно спам"),
			Key:      Ptr("possibleSpam"),
			Version:  Ptr(FlexString("0")),
			Schema:   &FieldSchema{Type: Ptr("integer"), Required: Ptr(false)},
			Readonly: Ptr(false),
			Options:  Ptr(true),
			Suggest:  Ptr(false),
			OptionsProvider: &OptionsProvider{
				Type:           Ptr("FixedListOptionsProvider"),
				NeedValidation: Ptr(true),
				Values:         []any{json.Number("0"), json.Number("1")},
			},
			QueryProvider: &QueryProvider{Type: Ptr("NumberOptionalQueryProvider")},
			Order:         Ptr(66),
			Category:      category,
			Type:          Ptr("standard"),
		},
	}

	if !reflect.DeepEqual(fields, want) {
		t.Errorf("List returned %+v, want %+v", fields, want)
	}
}

func TestFieldsService_Get_OptionValuesWrongType(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/fields/{id}", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"id":"possibleSpam","optionsProvider":{"type":"FixedListOptionsProvider","values":"0"}}`)
	})

	_, _, err := client.Fields.Get(context.Background(), "possibleSpam")
	if err == nil {
		t.Fatal("Get returned nil error for a string optionsProvider.values")
	}
	if !strings.Contains(err.Error(), "optionsProvider.values") {
		t.Errorf("Get error = %q, want it to name optionsProvider.values", err)
	}
}

func TestFieldsService_Create_OptionValues(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/fields", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"id":"myglobalfield","type":"ru.yandex.startrek.core.fields.StringFieldType","optionsProvider":{"type":"FixedListOptionsProvider","values":["list item 1","list item 2","list item 3"]}}`)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"id":"myglobalfield"}`)
	})

	_, _, err := client.Fields.Create(context.Background(), &FieldCreateRequest{
		ID:   Ptr("myglobalfield"),
		Type: Ptr("ru.yandex.startrek.core.fields.StringFieldType"),
		OptionsProvider: &OptionsProvider{
			Type:   Ptr("FixedListOptionsProvider"),
			Values: []any{"list item 1", "list item 2", "list item 3"},
		},
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
}

func TestOptionsProvider_UnmarshalJSON_Strings(t *testing.T) {
	var got OptionsProvider
	if err := json.Unmarshal([]byte(`{"type":"FixedListOptionsProvider","values":["S","M"]}`), &got); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	want := OptionsProvider{Type: Ptr("FixedListOptionsProvider"), Values: []any{"S", "M"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Unmarshal = %+v, want %+v", got, want)
	}
}

func TestOptionsProvider_UnmarshalJSON_Numbers(t *testing.T) {
	var got OptionsProvider
	if err := json.Unmarshal([]byte(`{"values":[0,1]}`), &got); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	want := OptionsProvider{Values: []any{json.Number("0"), json.Number("1")}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Unmarshal = %+v, want %+v", got, want)
	}
}

// TestOptionsProvider_UnmarshalJSON_PerQueue decodes the per-queue shape from
// the GET /v3/queues/{queue}/fields reference example: "values" maps a queue
// key to its list, and "defaults" holds the fallback list.
func TestOptionsProvider_UnmarshalJSON_PerQueue(t *testing.T) {
	var got Field
	err := json.Unmarshal([]byte(`{
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
	}`), &got)
	if err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	want := &OptionsProvider{
		Type: Ptr("QueueFixedListOptionsProvider"),
		QueueValues: map[string][]any{
			"DIRECT": {"Not specified", "Test", "Developer", "Beta", "Production", "Trunk"},
		},
		Defaults: []any{"Not specified", "Test", "Developer", "Beta", "Production"},
	}
	if !reflect.DeepEqual(got.OptionsProvider, want) {
		t.Errorf("OptionsProvider = %+v, want %+v", got.OptionsProvider, want)
	}
}

func TestOptionsProvider_UnmarshalJSON_NoValues(t *testing.T) {
	for _, input := range []string{
		`{"type":"TeamOptionsProvider"}`,
		`{"type":"TeamOptionsProvider","values":null,"defaults":null}`,
	} {
		var got OptionsProvider
		if err := json.Unmarshal([]byte(input), &got); err != nil {
			t.Fatalf("Unmarshal(%s) returned error: %v", input, err)
		}

		want := OptionsProvider{Type: Ptr("TeamOptionsProvider")}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Unmarshal(%s) = %+v, want %+v", input, got, want)
		}
	}
}

func TestOptionsProvider_UnmarshalJSON_UnknownValuesShape(t *testing.T) {
	for _, input := range []string{
		`{"type":"FixedListOptionsProvider","values":"S"}`,
		`{"type":"FixedListOptionsProvider","values":1}`,
		`{"type":"FixedListOptionsProvider","values":true}`,
		`{"type":"QueueFixedListOptionsProvider","values":{"DIRECT":"S"}}`,
	} {
		var got OptionsProvider
		err := json.Unmarshal([]byte(input), &got)
		if err == nil {
			t.Errorf("Unmarshal(%s) returned nil error", input)
			continue
		}
		if !strings.Contains(err.Error(), "optionsProvider.values") {
			t.Errorf("Unmarshal(%s) error = %q, want it to name optionsProvider.values", input, err)
		}
	}
}

func TestOptionsProvider_UnmarshalJSON_DefaultsNotArray(t *testing.T) {
	var got OptionsProvider
	err := json.Unmarshal([]byte(`{"type":"QueueFixedListOptionsProvider","defaults":"Test"}`), &got)
	if err == nil {
		t.Fatal("Unmarshal returned nil error for a string defaults")
	}
	if !strings.Contains(err.Error(), "optionsProvider.defaults") {
		t.Errorf("Unmarshal error = %q, want it to name optionsProvider.defaults", err)
	}
}

// TestOptionsProvider_MarshalJSON_RoundTrip re-encodes the recorded
// possibleSpam optionsProvider: its numbers must stay numbers.
func TestOptionsProvider_MarshalJSON_RoundTrip(t *testing.T) {
	const recorded = `{"type":"FixedListOptionsProvider","needValidation":true,"values":[0,1]}`

	var p OptionsProvider
	if err := json.Unmarshal([]byte(recorded), &p); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	data, err := json.Marshal(&p)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}
	if string(data) != recorded {
		t.Errorf("Marshal = %s, want %s", data, recorded)
	}
}

func TestOptionsProvider_MarshalJSON_QueueValues(t *testing.T) {
	p := OptionsProvider{
		Type:        Ptr("QueueFixedListOptionsProvider"),
		QueueValues: map[string][]any{"DIRECT": {"Test", "Beta"}},
		Defaults:    []any{"Test"},
	}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}

	want := `{"type":"QueueFixedListOptionsProvider","values":{"DIRECT":["Test","Beta"]},"defaults":["Test"]}`
	if string(data) != want {
		t.Errorf("Marshal = %s, want %s", data, want)
	}
}

func TestOptionsProvider_MarshalJSON_NoValues(t *testing.T) {
	data, err := json.Marshal(OptionsProvider{Type: Ptr("TeamOptionsProvider")})
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}

	if want := `{"type":"TeamOptionsProvider"}`; string(data) != want {
		t.Errorf("Marshal = %s, want %s", data, want)
	}
}

func TestOptionsProvider_MarshalJSON_BothValueShapes(t *testing.T) {
	req := &FieldEditRequest{OptionsProvider: &OptionsProvider{
		Type:        Ptr("FixedListOptionsProvider"),
		Values:      []any{"S"},
		QueueValues: map[string][]any{"DIRECT": {"S"}},
	}}

	_, err := json.Marshal(req)
	if err == nil {
		t.Fatal("Marshal returned nil error with both Values and QueueValues set")
	}
	if !strings.Contains(err.Error(), "Values and QueueValues") {
		t.Errorf("Marshal error = %q, want it to name Values and QueueValues", err)
	}
}

func TestOptionsProvider_UnmarshalJSON_PerQueueNumbers(t *testing.T) {
	var got OptionsProvider
	if err := json.Unmarshal([]byte(`{"values":{"Q":[0,1]},"defaults":[2]}`), &got); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	want := OptionsProvider{
		QueueValues: map[string][]any{"Q": {json.Number("0"), json.Number("1")}},
		Defaults:    []any{json.Number("2")},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Unmarshal = %+v, want %+v", got, want)
	}
}
