package tracker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestExternalLinksService_ListApplications(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v2/applications", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[
			{
				"self": "https://api.tracker.yandex.net/v2/applications/app1",
				"id": "app1",
				"type": "github",
				"name": "GitHub"
			},
			{
				"self": "https://api.tracker.yandex.net/v2/applications/app2",
				"id": "app2",
				"type": "bitbucket",
				"name": "Bitbucket"
			}
		]`)
	})

	apps, _, err := client.ExternalLinks.ListApplications(ctx)
	if err != nil {
		t.Fatalf("ExternalLinks.ListApplications returned error: %v", err)
	}

	if len(apps) != 2 {
		t.Fatalf("ExternalLinks.ListApplications returned %d apps, want 2", len(apps))
	}

	want := []*ExternalApplication{
		{
			Self: Ptr("https://api.tracker.yandex.net/v2/applications/app1"),
			ID:   Ptr("app1"),
			Type: Ptr("github"),
			Name: Ptr("GitHub"),
		},
		{
			Self: Ptr("https://api.tracker.yandex.net/v2/applications/app2"),
			ID:   Ptr("app2"),
			Type: Ptr("bitbucket"),
			Name: Ptr("Bitbucket"),
		},
	}

	if !reflect.DeepEqual(apps, want) {
		t.Errorf("ExternalLinks.ListApplications returned %+v, want %+v", apps, want)
	}
}

func TestExternalLinksService_ListLinks(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v2/issues/{key}/remotelinks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[
			{
				"self": "https://api.tracker.yandex.net/v2/issues/QUEUE-1/remotelinks/101",
				"id": 101,
				"type": {
					"self": "https://api.tracker.yandex.net/v2/linktypes/relates",
					"id": "relates",
					"key": "relates"
				},
				"direction": "outward",
				"object": {
					"self": "https://api.tracker.yandex.net/v2/applications/app1/objects/ext-1",
					"id": "ext-1",
					"key": "PR-42",
					"application": {
						"self": "https://api.tracker.yandex.net/v2/applications/app1",
						"id": "app1",
						"type": "github",
						"name": "GitHub"
					}
				},
				"createdBy": {
					"self": "https://api.tracker.yandex.net/v2/users/11",
					"id": "11",
					"display": "Ivan Ivanov"
				},
				"updatedBy": {
					"self": "https://api.tracker.yandex.net/v2/users/11",
					"id": "11",
					"display": "Ivan Ivanov"
				}
			}
		]`)
	})

	links, _, err := client.ExternalLinks.ListLinks(ctx, "QUEUE-1")
	if err != nil {
		t.Fatalf("ExternalLinks.ListLinks returned error: %v", err)
	}

	if len(links) != 1 {
		t.Fatalf("ExternalLinks.ListLinks returned %d links, want 1", len(links))
	}

	want := []*ExternalLink{
		{
			Self: Ptr("https://api.tracker.yandex.net/v2/issues/QUEUE-1/remotelinks/101"),
			ID:   Ptr(101),
			Type: &ExternalLinkType{
				Self: Ptr("https://api.tracker.yandex.net/v2/linktypes/relates"),
				ID:   Ptr("relates"),
				Key:  Ptr("relates"),
			},
			Direction: Ptr("outward"),
			Object: &ExternalLinkObject{
				Self: Ptr("https://api.tracker.yandex.net/v2/applications/app1/objects/ext-1"),
				ID:   Ptr("ext-1"),
				Key:  Ptr("PR-42"),
				Application: &ExternalApplication{
					Self: Ptr("https://api.tracker.yandex.net/v2/applications/app1"),
					ID:   Ptr("app1"),
					Type: Ptr("github"),
					Name: Ptr("GitHub"),
				},
			},
			CreatedBy: &User{
				Self:    Ptr("https://api.tracker.yandex.net/v2/users/11"),
				ID:      Ptr("11"),
				Display: Ptr("Ivan Ivanov"),
			},
			UpdatedBy: &User{
				Self:    Ptr("https://api.tracker.yandex.net/v2/users/11"),
				ID:      Ptr("11"),
				Display: Ptr("Ivan Ivanov"),
			},
		},
	}

	if !reflect.DeepEqual(links, want) {
		t.Errorf("ExternalLinks.ListLinks returned %+v, want %+v", links, want)
	}
}

func TestExternalLinksService_CreateLink(t *testing.T) {
	client, mux := setup(t)

	input := &ExternalLinkCreateRequest{
		Relationship: Ptr("relates"),
		Key:          Ptr("PR-42"),
		Origin:       Ptr("app1"),
	}

	mux.HandleFunc("POST /v2/issues/{key}/remotelinks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		// Verify request body.
		var got ExternalLinkCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("Decode request body: %v", err)
		}
		if *got.Relationship != "relates" {
			t.Errorf("Request body relationship = %q, want %q", *got.Relationship, "relates")
		}
		if *got.Key != "PR-42" {
			t.Errorf("Request body key = %q, want %q", *got.Key, "PR-42")
		}
		if *got.Origin != "app1" {
			t.Errorf("Request body origin = %q, want %q", *got.Origin, "app1")
		}

		// Verify backlink query param.
		if bl := r.URL.Query().Get("backlink"); bl != "true" {
			t.Errorf("Query param backlink = %q, want %q", bl, "true")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v2/issues/QUEUE-1/remotelinks/201",
			"id": 201,
			"type": {
				"self": "https://api.tracker.yandex.net/v2/linktypes/relates",
				"id": "relates",
				"key": "relates"
			},
			"direction": "outward",
			"object": {
				"self": "https://api.tracker.yandex.net/v2/applications/app1/objects/ext-1",
				"id": "ext-1",
				"key": "PR-42",
				"application": {
					"self": "https://api.tracker.yandex.net/v2/applications/app1",
					"id": "app1",
					"type": "github",
					"name": "GitHub"
				}
			}
		}`)
	})

	opts := &ExternalLinkCreateOptions{Backlink: Ptr(true)}
	link, _, err := client.ExternalLinks.CreateLink(ctx, "QUEUE-1", input, opts)
	if err != nil {
		t.Fatalf("ExternalLinks.CreateLink returned error: %v", err)
	}

	want := &ExternalLink{
		Self: Ptr("https://api.tracker.yandex.net/v2/issues/QUEUE-1/remotelinks/201"),
		ID:   Ptr(201),
		Type: &ExternalLinkType{
			Self: Ptr("https://api.tracker.yandex.net/v2/linktypes/relates"),
			ID:   Ptr("relates"),
			Key:  Ptr("relates"),
		},
		Direction: Ptr("outward"),
		Object: &ExternalLinkObject{
			Self: Ptr("https://api.tracker.yandex.net/v2/applications/app1/objects/ext-1"),
			ID:   Ptr("ext-1"),
			Key:  Ptr("PR-42"),
			Application: &ExternalApplication{
				Self: Ptr("https://api.tracker.yandex.net/v2/applications/app1"),
				ID:   Ptr("app1"),
				Type: Ptr("github"),
				Name: Ptr("GitHub"),
			},
		},
	}

	if !reflect.DeepEqual(link, want) {
		t.Errorf("ExternalLinks.CreateLink returned %+v, want %+v", link, want)
	}
}

func TestExternalLinksService_DeleteLink(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("DELETE /v2/issues/{key}/remotelinks/{id}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.ExternalLinks.DeleteLink(ctx, "QUEUE-1", 123)
	if err != nil {
		t.Fatalf("ExternalLinks.DeleteLink returned error: %v", err)
	}
	if resp == nil {
		t.Fatal("ExternalLinks.DeleteLink returned nil response")
	}
}
