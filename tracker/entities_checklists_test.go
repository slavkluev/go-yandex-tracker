package tracker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestEntitiesService_CreateChecklistItem(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/entities/project/{id}/checklistItems", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		// Verify request body
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["text"] != "New item" {
			t.Errorf("Request body text = %v, want %q", body["text"], "New item")
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/entities/project/1",
			"id": "1",
			"version": 2,
			"entityType": "project",
			"fields": {
				"summary": "Test Project"
			}
		}`)
	})

	entity, resp, err := client.Entities.CreateChecklistItem(ctx, EntityTypeProject, "1", &ChecklistItemRequest{
		Text: Ptr("New item"),
	})
	if err != nil {
		t.Fatalf("Entities.CreateChecklistItem returned error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	want := &Entity{
		Self:       Ptr("https://api.tracker.yandex.net/v3/entities/project/1"),
		ID:         Ptr(FlexString("1")),
		Version:    Ptr(FlexString("2")),
		EntityType: Ptr("project"),
		Fields: &EntityFields{
			Summary: Ptr("Test Project"),
		},
	}
	if !reflect.DeepEqual(entity, want) {
		t.Errorf("Entities.CreateChecklistItem = %+v, want %+v", entity, want)
	}
}

func TestEntitiesService_EditChecklistItem(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("PATCH /v3/entities/project/{id}/checklistItems/{itemId}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")

		// Verify request body
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["text"] != "Updated item" {
			t.Errorf("Request body text = %v, want %q", body["text"], "Updated item")
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/entities/project/1",
			"id": "1",
			"version": 3,
			"entityType": "project",
			"fields": {
				"summary": "Test Project"
			}
		}`)
	})

	entity, resp, err := client.Entities.EditChecklistItem(ctx, EntityTypeProject, "1", "item-1", &ChecklistItemRequest{
		Text: Ptr("Updated item"),
	})
	if err != nil {
		t.Fatalf("Entities.EditChecklistItem returned error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	want := &Entity{
		Self:       Ptr("https://api.tracker.yandex.net/v3/entities/project/1"),
		ID:         Ptr(FlexString("1")),
		Version:    Ptr(FlexString("3")),
		EntityType: Ptr("project"),
		Fields: &EntityFields{
			Summary: Ptr("Test Project"),
		},
	}
	if !reflect.DeepEqual(entity, want) {
		t.Errorf("Entities.EditChecklistItem = %+v, want %+v", entity, want)
	}
}

func TestEntitiesService_MoveChecklistItem(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/entities/project/{id}/checklistItems/{itemId}/_move", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		// Verify request body
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["before"] != "item-2" {
			t.Errorf("Request body before = %v, want %q", body["before"], "item-2")
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/entities/project/1",
			"id": "1",
			"version": 4,
			"entityType": "project",
			"fields": {
				"summary": "Test Project"
			}
		}`)
	})

	entity, resp, err := client.Entities.MoveChecklistItem(ctx, EntityTypeProject, "1", "item-1", &ChecklistMoveRequest{
		Before: Ptr("item-2"),
	})
	if err != nil {
		t.Fatalf("Entities.MoveChecklistItem returned error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	want := &Entity{
		Self:       Ptr("https://api.tracker.yandex.net/v3/entities/project/1"),
		ID:         Ptr(FlexString("1")),
		Version:    Ptr(FlexString("4")),
		EntityType: Ptr("project"),
		Fields: &EntityFields{
			Summary: Ptr("Test Project"),
		},
	}
	if !reflect.DeepEqual(entity, want) {
		t.Errorf("Entities.MoveChecklistItem = %+v, want %+v", entity, want)
	}
}

func TestEntitiesService_DeleteChecklistItem(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("DELETE /v3/entities/project/{id}/checklistItems/{itemId}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")

		// Entity checklist delete returns 200 with body, NOT 204
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/entities/project/1",
			"id": "1",
			"version": 5,
			"entityType": "project",
			"fields": {
				"summary": "Test Project"
			}
		}`)
	})

	entity, resp, err := client.Entities.DeleteChecklistItem(ctx, EntityTypeProject, "1", "item-1")
	if err != nil {
		t.Fatalf("Entities.DeleteChecklistItem returned error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	// Verify the response body was decoded into Entity
	want := &Entity{
		Self:       Ptr("https://api.tracker.yandex.net/v3/entities/project/1"),
		ID:         Ptr(FlexString("1")),
		Version:    Ptr(FlexString("5")),
		EntityType: Ptr("project"),
		Fields: &EntityFields{
			Summary: Ptr("Test Project"),
		},
	}
	if !reflect.DeepEqual(entity, want) {
		t.Errorf("Entities.DeleteChecklistItem = %+v, want %+v", entity, want)
	}
}

func TestEntitiesService_DeleteChecklist(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("DELETE /v3/entities/project/{id}/checklistItems", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")

		// Entity checklist delete-all returns 200 with body, NOT 204
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/entities/project/1",
			"id": "1",
			"version": 6,
			"entityType": "project",
			"fields": {
				"summary": "Test Project"
			}
		}`)
	})

	entity, resp, err := client.Entities.DeleteChecklist(ctx, EntityTypeProject, "1")
	if err != nil {
		t.Fatalf("Entities.DeleteChecklist returned error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	// Verify the response body was decoded into Entity
	want := &Entity{
		Self:       Ptr("https://api.tracker.yandex.net/v3/entities/project/1"),
		ID:         Ptr(FlexString("1")),
		Version:    Ptr(FlexString("6")),
		EntityType: Ptr("project"),
		Fields: &EntityFields{
			Summary: Ptr("Test Project"),
		},
	}
	if !reflect.DeepEqual(entity, want) {
		t.Errorf("Entities.DeleteChecklist = %+v, want %+v", entity, want)
	}
}
