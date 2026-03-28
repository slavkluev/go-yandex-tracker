package tracker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestEntitiesService_Create(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/entities/project", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)

		// Verify the fields wrapper is present
		fields, ok := body["fields"].(map[string]any)
		if !ok {
			t.Fatal("Request body missing 'fields' wrapper")
		}
		if fields["summary"] != "Test project" {
			t.Errorf("Request body fields.summary = %v, want %q", fields["summary"], "Test project")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"self":"https://api.tracker.yandex.net/v3/entities/project/1","id":"1","version":1,"entityType":"project","fields":{"summary":"Test project"}}`)
	})

	entity, resp, err := client.Entities.Create(ctx, EntityTypeProject, &EntityCreateRequest{
		Fields: &EntityCreateFields{
			Summary: Ptr("Test project"),
		},
	})
	if err != nil {
		t.Fatalf("Entities.Create returned error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusCreated)
	}

	want := &Entity{
		Self:       Ptr("https://api.tracker.yandex.net/v3/entities/project/1"),
		ID:         Ptr("1"),
		Version:    Ptr(1),
		EntityType: Ptr("project"),
		Fields: &EntityFields{
			Summary: Ptr("Test project"),
		},
	}
	if !reflect.DeepEqual(entity, want) {
		t.Errorf("Entities.Create = %+v, want %+v", entity, want)
	}
}

func TestEntitiesService_Create_InvalidJSON(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/entities/project", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"errorMessages":["Invalid request"],"errors":{}}`)
	})

	_, _, err := client.Entities.Create(ctx, EntityTypeProject, &EntityCreateRequest{
		Fields: &EntityCreateFields{
			Summary: Ptr("Test"),
		},
	})
	if err == nil {
		t.Fatal("Entities.Create expected error, got nil")
	}
}

func TestEntitiesService_Get(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/entities/project/{id}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"self":"https://api.tracker.yandex.net/v3/entities/project/1","id":"1","version":1,"shortId":1,"entityType":"project","fields":{"summary":"Test project","description":"A test project","teamAccess":true}}`)
	})

	entity, _, err := client.Entities.Get(ctx, EntityTypeProject, "1", nil)
	if err != nil {
		t.Fatalf("Entities.Get returned error: %v", err)
	}

	want := &Entity{
		Self:       Ptr("https://api.tracker.yandex.net/v3/entities/project/1"),
		ID:         Ptr("1"),
		Version:    Ptr(1),
		ShortID:    Ptr(1),
		EntityType: Ptr("project"),
		Fields: &EntityFields{
			Summary:     Ptr("Test project"),
			Description: Ptr("A test project"),
			TeamAccess:  Ptr(true),
		},
	}
	if !reflect.DeepEqual(entity, want) {
		t.Errorf("Entities.Get = %+v, want %+v", entity, want)
	}
}

func TestEntitiesService_Get_WithOptions(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/entities/project/{id}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		fields := r.URL.Query().Get("fields")
		if fields != "summary,description" {
			t.Errorf("fields query param = %q, want %q", fields, "summary,description")
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"id":"1","fields":{"summary":"Test project"}}`)
	})

	_, _, err := client.Entities.Get(ctx, EntityTypeProject, "1", &EntityGetOptions{
		Fields: "summary,description",
	})
	if err != nil {
		t.Fatalf("Entities.Get returned error: %v", err)
	}
}

func TestEntitiesService_Get_NotFound(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/entities/project/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"errorMessages":["Object not found"],"errors":{}}`)
	})

	_, _, err := client.Entities.Get(ctx, EntityTypeProject, "999", nil)
	if err == nil {
		t.Fatal("Entities.Get expected error for 404, got nil")
	}
	if !IsNotFound(err) {
		t.Errorf("IsNotFound(err) = false, want true")
	}
}

func TestEntitiesService_Update(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("PATCH /v3/entities/project/{id}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")

		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)

		fields, ok := body["fields"].(map[string]any)
		if !ok {
			t.Fatal("Request body missing 'fields' wrapper")
		}
		if fields["summary"] != "Updated project" {
			t.Errorf("Request body fields.summary = %v, want %q", fields["summary"], "Updated project")
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"self":"https://api.tracker.yandex.net/v3/entities/project/1","id":"1","version":2,"entityType":"project","fields":{"summary":"Updated project"}}`)
	})

	entity, _, err := client.Entities.Update(ctx, EntityTypeProject, "1", &EntityUpdateRequest{
		Fields: &EntityCreateFields{
			Summary: Ptr("Updated project"),
		},
	})
	if err != nil {
		t.Fatalf("Entities.Update returned error: %v", err)
	}

	want := &Entity{
		Self:       Ptr("https://api.tracker.yandex.net/v3/entities/project/1"),
		ID:         Ptr("1"),
		Version:    Ptr(2),
		EntityType: Ptr("project"),
		Fields: &EntityFields{
			Summary: Ptr("Updated project"),
		},
	}
	if !reflect.DeepEqual(entity, want) {
		t.Errorf("Entities.Update = %+v, want %+v", entity, want)
	}
}

func TestEntitiesService_Delete(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("DELETE /v3/entities/project/{id}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.Entities.Delete(ctx, EntityTypeProject, "1", nil)
	if err != nil {
		t.Fatalf("Entities.Delete returned error: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusNoContent)
	}
}

func TestEntitiesService_Delete_WithBoard(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("DELETE /v3/entities/project/{id}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")

		withBoard := r.URL.Query().Get("withBoard")
		if withBoard != "true" {
			t.Errorf("withBoard query param = %q, want %q", withBoard, "true")
		}

		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Entities.Delete(ctx, EntityTypeProject, "1", &EntityDeleteOptions{
		WithBoard: Ptr(true),
	})
	if err != nil {
		t.Fatalf("Entities.Delete returned error: %v", err)
	}
}

func TestEntitiesService_Search(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/entities/project/_search", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)

		if body["input"] != "Test" {
			t.Errorf("Request body input = %v, want %q", body["input"], "Test")
		}

		page := r.URL.Query().Get("page")
		if page != "1" {
			t.Errorf("page query param = %q, want %q", page, "1")
		}
		perPage := r.URL.Query().Get("perPage")
		if perPage != "10" {
			t.Errorf("perPage query param = %q, want %q", perPage, "10")
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"hits":1,"pages":1,"values":[{"self":"https://api.tracker.yandex.net/v3/entities/project/1","id":"1","entityType":"project","fields":{"summary":"Test project"}}],"orderBy":"updated"}`)
	})

	result, _, err := client.Entities.Search(ctx, EntityTypeProject, &EntitySearchRequest{
		Input: Ptr("Test"),
	}, &EntitySearchOptions{
		ListOptions: ListOptions{Page: 1, PerPage: 10},
	})
	if err != nil {
		t.Fatalf("Entities.Search returned error: %v", err)
	}
	if got, want := *result.Hits, 1; got != want {
		t.Errorf("Hits = %d, want %d", got, want)
	}
	if got, want := *result.Pages, 1; got != want {
		t.Errorf("Pages = %d, want %d", got, want)
	}
	if got, want := len(result.Values), 1; got != want {
		t.Errorf("len(Values) = %d, want %d", got, want)
	}
	if got, want := *result.OrderBy, "updated"; got != want {
		t.Errorf("OrderBy = %q, want %q", got, want)
	}
}

func TestEntitiesService_Search_WithFilter(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/entities/project/_search", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)

		if body["rootOnly"] != true {
			t.Errorf("Request body rootOnly = %v, want true", body["rootOnly"])
		}
		filter, ok := body["filter"].(map[string]any)
		if !ok {
			t.Fatal("Request body missing 'filter'")
		}
		if filter["entityType"] != "project" {
			t.Errorf("filter.entityType = %v, want %q", filter["entityType"], "project")
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"hits":0,"pages":0,"values":[]}`)
	})

	result, _, err := client.Entities.Search(ctx, EntityTypeProject, &EntitySearchRequest{
		RootOnly: Ptr(true),
		Filter:   map[string]any{"entityType": "project"},
	}, nil)
	if err != nil {
		t.Fatalf("Entities.Search returned error: %v", err)
	}
	if got, want := *result.Hits, 0; got != want {
		t.Errorf("Hits = %d, want %d", got, want)
	}
}

func TestEntitiesService_BulkChange(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/entities/project/bulkchange/_update", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)

		metaEntities, ok := body["metaEntities"].([]any)
		if !ok {
			t.Fatal("Request body missing 'metaEntities'")
		}
		if len(metaEntities) != 2 {
			t.Errorf("len(metaEntities) = %d, want 2", len(metaEntities))
		}
		if metaEntities[0] != "1" || metaEntities[1] != "2" {
			t.Errorf("metaEntities = %v, want [1, 2]", metaEntities)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"self":"https://api.tracker.yandex.net/v3/bulkchange/1","id":"1","status":"CREATED","statusText":"Bulk change created"}`)
	})

	bc, resp, err := client.Entities.BulkChange(ctx, EntityTypeProject, &EntityBulkChangeRequest{
		MetaEntities: []string{"1", "2"},
		Values: &EntityBulkChangeValues{
			Fields: map[string]any{"summary": "Updated"},
		},
	})
	if err != nil {
		t.Fatalf("Entities.BulkChange returned error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusCreated)
	}
	if got, want := *bc.ID, "1"; got != want {
		t.Errorf("ID = %q, want %q", got, want)
	}
	if got, want := *bc.Status, "CREATED"; got != want {
		t.Errorf("Status = %q, want %q", got, want)
	}
}

func TestEntitiesService_GetEvents(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/entities/project/{id}/events/_relative", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		perPage := r.URL.Query().Get("perPage")
		if perPage != "10" {
			t.Errorf("perPage query param = %q, want %q", perPage, "10")
		}
		direction := r.URL.Query().Get("direction")
		if direction != "asc" {
			t.Errorf("direction query param = %q, want %q", direction, "asc")
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"events":[{"id":"event-1","transport":"tracker","display":"Field changed"}],"hasNext":true,"hasPrev":false}`)
	})

	result, _, err := client.Entities.GetEvents(ctx, EntityTypeProject, "1", &EntityEventsOptions{
		PerPage:   10,
		Direction: "asc",
	})
	if err != nil {
		t.Fatalf("Entities.GetEvents returned error: %v", err)
	}
	if got, want := len(result.Events), 1; got != want {
		t.Errorf("len(Events) = %d, want %d", got, want)
	}
	if got, want := *result.Events[0].ID, "event-1"; got != want {
		t.Errorf("Events[0].ID = %q, want %q", got, want)
	}
	if got, want := *result.HasNext, true; got != want {
		t.Errorf("HasNext = %v, want %v", got, want)
	}
	if got, want := *result.HasPrev, false; got != want {
		t.Errorf("HasPrev = %v, want %v", got, want)
	}
}

func TestEntitiesService_GetEvents_CursorPagination(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/entities/project/{id}/events/_relative", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		from := r.URL.Query().Get("from")
		if from != "event-1" {
			t.Errorf("from query param = %q, want %q", from, "event-1")
		}
		direction := r.URL.Query().Get("direction")
		if direction != "next" {
			t.Errorf("direction query param = %q, want %q", direction, "next")
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"events":[{"id":"event-2","transport":"tracker","display":"Status changed"}],"hasNext":false,"hasPrev":true}`)
	})

	result, _, err := client.Entities.GetEvents(ctx, EntityTypeProject, "1", &EntityEventsOptions{
		From:      "event-1",
		Direction: "next",
	})
	if err != nil {
		t.Fatalf("Entities.GetEvents returned error: %v", err)
	}
	if got, want := len(result.Events), 1; got != want {
		t.Errorf("len(Events) = %d, want %d", got, want)
	}
	if got, want := *result.Events[0].ID, "event-2"; got != want {
		t.Errorf("Events[0].ID = %q, want %q", got, want)
	}
	if got, want := *result.HasNext, false; got != want {
		t.Errorf("HasNext = %v, want %v", got, want)
	}
	if got, want := *result.HasPrev, true; got != want {
		t.Errorf("HasPrev = %v, want %v", got, want)
	}
}
