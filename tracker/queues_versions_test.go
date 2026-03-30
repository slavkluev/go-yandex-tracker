package tracker

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
)

func TestQueuesService_ListVersions(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/queues/{key}/versions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]*QueueVersion{
			{
				Self:      Ptr("https://api.tracker.yandex.net/v3/queues/TEST/versions/1"),
				ID:        Ptr(FlexString("1")),
				Name:      Ptr("v1.0"),
				StartDate: Ptr("2024-01-01"),
				DueDate:   Ptr("2024-06-30"),
				Released:  Ptr(false),
				Archived:  Ptr(false),
			},
			{
				Self:     Ptr("https://api.tracker.yandex.net/v3/queues/TEST/versions/2"),
				ID:       Ptr(FlexString("2")),
				Name:     Ptr("v2.0"),
				Released: Ptr(true),
				Archived: Ptr(false),
			},
		})
	})

	versions, _, err := client.Queues.ListVersions(ctx, "TEST")
	if err != nil {
		t.Fatalf("Queues.ListVersions returned error: %v", err)
	}

	if len(versions) != 2 {
		t.Fatalf("Queues.ListVersions returned %q versions, want 2", len(versions))
	}

	want := []*QueueVersion{
		{
			Self:      Ptr("https://api.tracker.yandex.net/v3/queues/TEST/versions/1"),
			ID:        Ptr(FlexString("1")),
			Name:      Ptr("v1.0"),
			StartDate: Ptr("2024-01-01"),
			DueDate:   Ptr("2024-06-30"),
			Released:  Ptr(false),
			Archived:  Ptr(false),
		},
		{
			Self:     Ptr("https://api.tracker.yandex.net/v3/queues/TEST/versions/2"),
			ID:       Ptr(FlexString("2")),
			Name:     Ptr("v2.0"),
			Released: Ptr(true),
			Archived: Ptr(false),
		},
	}

	if !reflect.DeepEqual(versions, want) {
		t.Errorf("Queues.ListVersions returned %+v, want %+v", versions, want)
	}
}

func TestQueuesService_ListTags(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/queues/{key}/tags", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]string{"tag1", "tag2", "tag3"})
	})

	tags, _, err := client.Queues.ListTags(ctx, "TEST")
	if err != nil {
		t.Fatalf("Queues.ListTags returned error: %v", err)
	}

	want := []string{"tag1", "tag2", "tag3"}
	if !reflect.DeepEqual(tags, want) {
		t.Errorf("Queues.ListTags returned %v, want %v", tags, want)
	}
}
