package tracker

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestStatusesService_List(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/statuses", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[
			{
				"self": "https://api.tracker.yandex.net/v3/statuses/1",
				"id": "1",
				"key": "open",
				"name": "Open",
				"version": 1,
				"description": "Issue is open",
				"order": 1,
				"type": "new"
			},
			{
				"self": "https://api.tracker.yandex.net/v3/statuses/2",
				"id": "2",
				"key": "closed",
				"name": "Closed",
				"version": 1,
				"order": 99,
				"type": "done"
			}
		]`)
	})

	ctx := context.Background()
	statuses, _, err := client.Statuses.List(ctx)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	if len(statuses) != 2 {
		t.Fatalf("List returned %d statuses, want 2", len(statuses))
	}

	want := &Status{
		Self:        Ptr("https://api.tracker.yandex.net/v3/statuses/1"),
		ID:          Ptr(FlexString("1")),
		Key:         Ptr("open"),
		Name:        Ptr("Open"),
		Version:     Ptr(FlexString("1")),
		Description: Ptr("Issue is open"),
		Order:       Ptr(1),
		Type:        Ptr("new"),
	}

	if !reflect.DeepEqual(statuses[0], want) {
		t.Errorf("List returned %+v, want %+v", statuses[0], want)
	}
}
