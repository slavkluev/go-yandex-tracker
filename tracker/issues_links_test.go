package tracker

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestIssuesService_GetLinks(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/issues/{key}/links", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{
			"id": "1",
			"type": {"id": "relates", "inward": "relates to", "outward": "relates to"},
			"direction": "outward",
			"object": {"key": "QUEUE-2", "summary": "Related issue"}
		}]`)
	})

	ctx := context.Background()
	links, _, err := client.Issues.GetLinks(ctx, "QUEUE-1")
	if err != nil {
		t.Fatalf("GetLinks returned error: %v", err)
	}

	if len(links) != 1 {
		t.Fatalf("GetLinks returned %d links, want 1", len(links))
	}

	link := links[0]
	if got := *link.ID; got != "1" {
		t.Errorf("ID = %q, want %q", got, "1")
	}
	if got := *link.Type.ID; got != "relates" {
		t.Errorf("Type.ID = %q, want %q", got, "relates")
	}
	if got := *link.Direction; got != "outward" {
		t.Errorf("Direction = %q, want %q", got, "outward")
	}
	if got := *link.Object.Key; got != "QUEUE-2" {
		t.Errorf("Object.Key = %q, want %q", got, "QUEUE-2")
	}
}

func TestIssuesService_GetLinks_Empty(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/issues/{key}/links", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[]`)
	})

	ctx := context.Background()
	links, _, err := client.Issues.GetLinks(ctx, "QUEUE-1")
	if err != nil {
		t.Fatalf("GetLinks returned error: %v", err)
	}
	if len(links) != 0 {
		t.Errorf("GetLinks returned %d links, want 0", len(links))
	}
}

func TestIssuesService_CreateLink(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/issues/{key}/links", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"relationship":"relates","issue":"TREK-2"}`)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"id": "5",
			"type": {"id": "relates", "inward": "relates to", "outward": "relates to"},
			"direction": "outward",
			"object": {"key": "TREK-2", "summary": "Linked issue"}
		}`)
	})

	ctx := context.Background()
	link, resp, err := client.Issues.CreateLink(ctx, "QUEUE-1", &LinkRequest{
		Relationship: Ptr("relates"),
		Issue:        Ptr("TREK-2"),
	})
	if err != nil {
		t.Fatalf("CreateLink returned error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusCreated)
	}
	if got := *link.ID; got != "5" {
		t.Errorf("ID = %q, want %q", got, "5")
	}
	if got := *link.Type.ID; got != "relates" {
		t.Errorf("Type.ID = %q, want %q", got, "relates")
	}
	if got := *link.Object.Key; got != "TREK-2" {
		t.Errorf("Object.Key = %q, want %q", got, "TREK-2")
	}
}

func TestIssuesService_DeleteLink(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("DELETE /v3/issues/{key}/links/{id}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	resp, err := client.Issues.DeleteLink(ctx, "QUEUE-1", "123")
	if err != nil {
		t.Fatalf("DeleteLink returned error: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusNoContent)
	}
}
