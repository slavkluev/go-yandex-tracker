package tracker

import (
	"fmt"
	"net/http"
	"testing"
)

func TestEntitiesService_ListLinks(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/entities/project/{id}/links", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{
			"type": "depends on",
			"linkFieldValues": {
				"summary": "Linked entity",
				"id": "entity-2"
			}
		}]`)
	})

	links, _, err := client.Entities.ListLinks(ctx, EntityTypeProject, "1")
	if err != nil {
		t.Fatalf("ListLinks returned error: %v", err)
	}

	if len(links) != 1 {
		t.Fatalf("ListLinks returned %d links, want 1", len(links))
	}

	link := links[0]
	if got := *link.Type; got != "depends on" {
		t.Errorf("Type = %q, want %q", got, "depends on")
	}
	if got, ok := link.LinkFieldValues["summary"]; !ok || got != "Linked entity" {
		t.Errorf("LinkFieldValues[summary] = %v, want %q", got, "Linked entity")
	}
	if got, ok := link.LinkFieldValues["id"]; !ok || got != "entity-2" {
		t.Errorf("LinkFieldValues[id] = %v, want %q", got, "entity-2")
	}
}

func TestEntitiesService_CreateLink(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/entities/project/{id}/links", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		// Verify request body is a JSON array
		testBody(t, r, `[{"relationship":"depends on","entity":"entity-2"}]`)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/entities/project/1",
			"id": "1",
			"entityType": "project",
			"fields": {"summary": "Test Project"}
		}`)
	})

	links := []*EntityLink{
		{
			Relationship: Ptr("depends on"),
			Entity:       Ptr("entity-2"),
		},
	}

	entity, _, err := client.Entities.CreateLink(ctx, EntityTypeProject, "1", links)
	if err != nil {
		t.Fatalf("CreateLink returned error: %v", err)
	}

	if got := *entity.ID; got != "1" {
		t.Errorf("ID = %q, want %q", got, "1")
	}
	if got := *entity.EntityType; got != "project" {
		t.Errorf("EntityType = %q, want %q", got, "project")
	}
	if entity.Fields == nil {
		t.Fatal("Fields is nil, want non-nil")
	}
	if got := *entity.Fields.Summary; got != "Test Project" {
		t.Errorf("Fields.Summary = %q, want %q", got, "Test Project")
	}
}

func TestEntitiesService_DeleteLink(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("DELETE /v3/entities/project/{id}/links", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")

		// Verify query parameter (not path parameter)
		if got := r.URL.Query().Get("right"); got != "entity-2" {
			t.Errorf("Query param right = %q, want %q", got, "entity-2")
		}

		// Entity link delete returns 200 OK (not 204)
		w.WriteHeader(http.StatusOK)
	})

	resp, err := client.Entities.DeleteLink(ctx, EntityTypeProject, "1", &EntityLinkDeleteOptions{
		Right: "entity-2",
	})
	if err != nil {
		t.Fatalf("DeleteLink returned error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}
