package tracker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestQueuesService_Create(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/queues/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "application/json")

		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["key"] != "TEST" {
			t.Errorf("Request body key = %v, want %v", body["key"], "TEST")
		}
		if body["name"] != "Test Queue" {
			t.Errorf("Request body name = %v, want %v", body["name"], "Test Queue")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Queue{
			Self:       Ptr("https://api.tracker.yandex.net/v3/queues/TEST"),
			ID:         Ptr(FlexString("1")),
			Key:        Ptr("TEST"),
			Display:    Ptr("Test Queue"),
			Version:    Ptr(FlexString("1")),
			Name:       Ptr("Test Queue"),
			Lead:       &User{Self: Ptr("https://api.tracker.yandex.net/v3/users/123"), ID: Ptr(FlexString("123")), Display: Ptr("John Doe")},
			AssignAuto: Ptr(false),
			DefaultType: &IssueType{
				Self:    Ptr("https://api.tracker.yandex.net/v3/issuetypes/2"),
				ID:      Ptr(FlexString("2")),
				Key:     Ptr("task"),
				Display: Ptr("Task"),
			},
			DefaultPriority: &Priority{
				Self:    Ptr("https://api.tracker.yandex.net/v3/priorities/3"),
				ID:      Ptr(FlexString("3")),
				Key:     Ptr("normal"),
				Display: Ptr("Normal"),
			},
		})
	})

	queue, _, err := client.Queues.Create(ctx, &QueueCreateRequest{
		Key:  Ptr("TEST"),
		Name: Ptr("Test Queue"),
		Lead: Ptr("john.doe"),
	})
	if err != nil {
		t.Fatalf("Queues.Create returned error: %v", err)
	}

	want := &Queue{
		Self:       Ptr("https://api.tracker.yandex.net/v3/queues/TEST"),
		ID:         Ptr(FlexString("1")),
		Key:        Ptr("TEST"),
		Display:    Ptr("Test Queue"),
		Version:    Ptr(FlexString("1")),
		Name:       Ptr("Test Queue"),
		Lead:       &User{Self: Ptr("https://api.tracker.yandex.net/v3/users/123"), ID: Ptr(FlexString("123")), Display: Ptr("John Doe")},
		AssignAuto: Ptr(false),
		DefaultType: &IssueType{
			Self:    Ptr("https://api.tracker.yandex.net/v3/issuetypes/2"),
			ID:      Ptr(FlexString("2")),
			Key:     Ptr("task"),
			Display: Ptr("Task"),
		},
		DefaultPriority: &Priority{
			Self:    Ptr("https://api.tracker.yandex.net/v3/priorities/3"),
			ID:      Ptr(FlexString("3")),
			Key:     Ptr("normal"),
			Display: Ptr("Normal"),
		},
	}

	if !reflect.DeepEqual(queue, want) {
		t.Errorf("Queues.Create returned %+v, want %+v", queue, want)
	}
}

func TestQueuesService_Get(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/queues/{key}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		if got := r.URL.Query().Get("expand"); got != "all" {
			t.Errorf("expand = %q, want %q", got, "all")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Queue{
			Self:    Ptr("https://api.tracker.yandex.net/v3/queues/TEST"),
			ID:      Ptr(FlexString("1")),
			Key:     Ptr("TEST"),
			Display: Ptr("Test Queue"),
			Version: Ptr(FlexString("1")),
			Name:    Ptr("Test Queue"),
			Lead:    &User{Self: Ptr("https://api.tracker.yandex.net/v3/users/123"), ID: Ptr(FlexString("123")), Display: Ptr("John Doe")},
		})
	})

	queue, _, err := client.Queues.Get(ctx, "TEST", &QueueGetOptions{Expand: "all"})
	if err != nil {
		t.Fatalf("Queues.Get returned error: %v", err)
	}

	want := &Queue{
		Self:    Ptr("https://api.tracker.yandex.net/v3/queues/TEST"),
		ID:      Ptr(FlexString("1")),
		Key:     Ptr("TEST"),
		Display: Ptr("Test Queue"),
		Version: Ptr(FlexString("1")),
		Name:    Ptr("Test Queue"),
		Lead:    &User{Self: Ptr("https://api.tracker.yandex.net/v3/users/123"), ID: Ptr(FlexString("123")), Display: Ptr("John Doe")},
	}

	if !reflect.DeepEqual(queue, want) {
		t.Errorf("Queues.Get returned %+v, want %+v", queue, want)
	}
}

func TestQueuesService_Get_Expand(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/queues/{key}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		if got := r.URL.Query().Get("expand"); got != "all" {
			t.Errorf("expand = %q, want %q", got, "all")
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/queues/TEST",
			"id": "3",
			"key": "TEST",
			"version": 5,
			"name": "Test",
			"description": "Queue created for testing purposes",
			"lead": {
				"self": "https://api.tracker.yandex.net/v3/users/11********",
				"id": "11********",
				"display": "Ivan Ivanov"
			},
			"assignAuto": false,
			"defaultType": {
				"self": "https://api.tracker.yandex.net/v3/issuetypes/1",
				"id": "1",
				"key": "bug",
				"display": "Error"
			},
			"defaultPriority": {
				"self": "https://api.tracker.yandex.net/v3/priorities/3",
				"id": "3",
				"key": "normal",
				"display": "Normal"
			},
			"teamUsers": [
				{
					"self": "https://api.tracker.yandex.net/v3/users/11********",
					"id": "11********",
					"display": "Ivan Ivanov"
				}
			],
			"issueTypes": [
				{
					"self": "https://api.tracker.yandex.net/v3/issuetypes/1",
					"id": "1",
					"key": "bug",
					"display": "Error"
				}
			],
			"versions": [
				{
					"self": "https://api.tracker.yandex.net/v3/versions/4",
					"id": "4",
					"display": "Peek-a-boo"
				},
				{"self": "https://api.tracker.yandex.net/v3/versions/2", "id": "2", "display": "2023 Q1"}
			],
			"workflows": {
				"dev": [
					{
						"self": "https://api.tracker.yandex.net/v3/issuetypes/1",
						"id": "1",
						"key": "bug",
						"display": "Errror"
					}
				]
			},
			"denyVoting": false,
			"issueTypesConfig": [
				{
					"issueType": {
						"self": "https://api.tracker.yandex.net/v3/issuetypes/1",
						"id": "1",
						"key": "bug",
						"display": "Error"
					},
					"workflow": {
						"self": "https://api.tracker.yandex.net/v3/workflows/dev",
						"id": "dev",
						"display": "dev"
					},
					"resolutions": [
						{
							"self": "https://api.tracker.yandex.net/v3/resolutions/2",
							"id": "2",
							"key": "wontFix",
							"display": "Won't fix"
						}
					]
				}
			],
			"components": [{"self": "https://api.tracker.yandex.net/v3/components/56", "id": "56", "display": "Standard"}],
			"fields": [{"self": "https://api.tracker.yandex.net/v3/queues/APP/localFields/size", "id": "5d0e4f1a2b3c4d5e6f708192--size", "display": "Размер"}]
		}`)
	})

	queue, _, err := client.Queues.Get(ctx, "TEST", &QueueGetOptions{Expand: "all"})
	if err != nil {
		t.Fatalf("Queues.Get returned error: %v", err)
	}

	ivan := &User{
		Self:    Ptr("https://api.tracker.yandex.net/v3/users/11********"),
		ID:      Ptr(FlexString("11********")),
		Display: Ptr("Ivan Ivanov"),
	}
	bug := &IssueType{
		Self:    Ptr("https://api.tracker.yandex.net/v3/issuetypes/1"),
		ID:      Ptr(FlexString("1")),
		Key:     Ptr("bug"),
		Display: Ptr("Error"),
	}
	want := &Queue{
		Self:        Ptr("https://api.tracker.yandex.net/v3/queues/TEST"),
		ID:          Ptr(FlexString("3")),
		Key:         Ptr("TEST"),
		Version:     Ptr(FlexString("5")),
		Name:        Ptr("Test"),
		Description: Ptr("Queue created for testing purposes"),
		Lead:        ivan,
		AssignAuto:  Ptr(false),
		DenyVoting:  Ptr(false),
		DefaultType: bug,
		DefaultPriority: &Priority{
			Self:    Ptr("https://api.tracker.yandex.net/v3/priorities/3"),
			ID:      Ptr(FlexString("3")),
			Key:     Ptr("normal"),
			Display: Ptr("Normal"),
		},
		TeamUsers:  []*User{ivan},
		IssueTypes: []*IssueType{bug},
		Versions: []*QueueVersion{
			{Self: Ptr("https://api.tracker.yandex.net/v3/versions/4"), ID: Ptr(FlexString("4")), Display: Ptr("Peek-a-boo")},
			{Self: Ptr("https://api.tracker.yandex.net/v3/versions/2"), ID: Ptr(FlexString("2")), Display: Ptr("2023 Q1")},
		},
		Components: []*Component{
			{Self: Ptr("https://api.tracker.yandex.net/v3/components/56"), ID: Ptr(FlexString("56")), Display: Ptr("Standard")},
		},
		Workflows: map[string][]*IssueType{
			"dev": {
				{
					Self:    Ptr("https://api.tracker.yandex.net/v3/issuetypes/1"),
					ID:      Ptr(FlexString("1")),
					Key:     Ptr("bug"),
					Display: Ptr("Errror"),
				},
			},
		},
		Fields: []*Field{
			{
				Self:    Ptr("https://api.tracker.yandex.net/v3/queues/APP/localFields/size"),
				ID:      Ptr(FlexString("5d0e4f1a2b3c4d5e6f708192--size")),
				Display: Ptr("Размер"),
			},
		},
		IssueTypesConfig: []*QueueIssueTypeConfig{
			{
				IssueType: bug,
				Workflow: &Workflow{
					Self:    Ptr("https://api.tracker.yandex.net/v3/workflows/dev"),
					ID:      Ptr(FlexString("dev")),
					Display: Ptr("dev"),
				},
				Resolutions: []*Resolution{
					{
						Self:    Ptr("https://api.tracker.yandex.net/v3/resolutions/2"),
						ID:      Ptr(FlexString("2")),
						Key:     Ptr("wontFix"),
						Display: Ptr("Won't fix"),
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(queue, want) {
		t.Errorf("Queues.Get returned %+v, want %+v", queue, want)
	}
}

func TestQueuesService_Get_IssueTypesConfigWithoutResolutions(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/queues/{key}", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"key": "RESTRICTED",
			"issueTypesConfig": [
				{"issueType": {"self": "https://api.tracker.yandex.net/v3/issuetypes/21", "id": "21", "key": "milestone", "display": "Веха"},
				 "workflow": {"self": "https://api.tracker.yandex.net/v3/workflows/W100", "id": "W100", "display": "W100"}}
			]
		}`)
	})

	queue, _, err := client.Queues.Get(ctx, "RESTRICTED", &QueueGetOptions{Expand: "issueTypesConfig"})
	if err != nil {
		t.Fatalf("Queues.Get returned error: %v", err)
	}

	want := []*QueueIssueTypeConfig{
		{
			IssueType: &IssueType{
				Self:    Ptr("https://api.tracker.yandex.net/v3/issuetypes/21"),
				ID:      Ptr(FlexString("21")),
				Key:     Ptr("milestone"),
				Display: Ptr("Веха"),
			},
			Workflow: &Workflow{
				Self:    Ptr("https://api.tracker.yandex.net/v3/workflows/W100"),
				ID:      Ptr(FlexString("W100")),
				Display: Ptr("W100"),
			},
		},
	}

	if !reflect.DeepEqual(queue.IssueTypesConfig, want) {
		t.Fatalf("Queue.IssueTypesConfig = %+v, want %+v", queue.IssueTypesConfig, want)
	}
	if got := queue.IssueTypesConfig[0].Resolutions; got != nil {
		t.Errorf("Queue.IssueTypesConfig[0].Resolutions = %+v, want nil", got)
	}
}

func TestQueuesService_List(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/queues", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		if got := r.URL.Query().Get("page"); got != "1" {
			t.Errorf("page = %q, want %q", got, "1")
		}
		if got := r.URL.Query().Get("perPage"); got != "10" {
			t.Errorf("perPage = %q, want %q", got, "10")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]*Queue{
			{
				Self:    Ptr("https://api.tracker.yandex.net/v3/queues/TEST"),
				ID:      Ptr(FlexString("1")),
				Key:     Ptr("TEST"),
				Display: Ptr("Test Queue"),
				Name:    Ptr("Test Queue"),
			},
			{
				Self:    Ptr("https://api.tracker.yandex.net/v3/queues/DEV"),
				ID:      Ptr(FlexString("2")),
				Key:     Ptr("DEV"),
				Display: Ptr("Dev Queue"),
				Name:    Ptr("Dev Queue"),
			},
		})
	})

	queues, _, err := client.Queues.List(ctx, &QueueListOptions{
		ListOptions: ListOptions{Page: 1, PerPage: 10},
	})
	if err != nil {
		t.Fatalf("Queues.List returned error: %v", err)
	}

	if len(queues) != 2 {
		t.Fatalf("Queues.List returned %d queues, want 2", len(queues))
	}

	want := []*Queue{
		{
			Self:    Ptr("https://api.tracker.yandex.net/v3/queues/TEST"),
			ID:      Ptr(FlexString("1")),
			Key:     Ptr("TEST"),
			Display: Ptr("Test Queue"),
			Name:    Ptr("Test Queue"),
		},
		{
			Self:    Ptr("https://api.tracker.yandex.net/v3/queues/DEV"),
			ID:      Ptr(FlexString("2")),
			Key:     Ptr("DEV"),
			Display: Ptr("Dev Queue"),
			Name:    Ptr("Dev Queue"),
		},
	}

	if !reflect.DeepEqual(queues, want) {
		t.Errorf("Queues.List returned %+v, want %+v", queues, want)
	}
}

func TestQueuesService_Delete(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("DELETE /v3/queues/{key}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.Queues.Delete(ctx, "TEST")
	if err != nil {
		t.Fatalf("Queues.Delete returned error: %v", err)
	}
	if resp == nil {
		t.Fatal("Queues.Delete returned nil response")
	}
}

func TestQueuesService_Restore(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/queues/{key}/_restore", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Queue{
			Self:    Ptr("https://api.tracker.yandex.net/v3/queues/TEST"),
			ID:      Ptr(FlexString("1")),
			Key:     Ptr("TEST"),
			Display: Ptr("Test Queue"),
			Version: Ptr(FlexString("2")),
			Name:    Ptr("Test Queue"),
		})
	})

	queue, _, err := client.Queues.Restore(ctx, "TEST")
	if err != nil {
		t.Fatalf("Queues.Restore returned error: %v", err)
	}

	want := &Queue{
		Self:    Ptr("https://api.tracker.yandex.net/v3/queues/TEST"),
		ID:      Ptr(FlexString("1")),
		Key:     Ptr("TEST"),
		Display: Ptr("Test Queue"),
		Version: Ptr(FlexString("2")),
		Name:    Ptr("Test Queue"),
	}

	if !reflect.DeepEqual(queue, want) {
		t.Errorf("Queues.Restore returned %+v, want %+v", queue, want)
	}
}

func TestQueuesService_Get_NotFound(t *testing.T) {
	client, mux := setup(t)
	ctx := context.Background()

	mux.HandleFunc("GET /v3/queues/{key}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"errorMessages":["Object not found"],"errors":{}}`)
	})

	_, _, err := client.Queues.Get(ctx, "NONEXIST", nil)
	if err == nil {
		t.Fatal("Queues.Get expected error for 404, got nil")
	}
	if !IsNotFound(err) {
		t.Errorf("IsNotFound(err) = false, want true")
	}
}
