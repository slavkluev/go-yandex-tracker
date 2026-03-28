package tracker

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestResolutionsService_List(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/resolutions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[
			{
				"self": "https://api.tracker.yandex.net/v3/resolutions/1",
				"id": "1",
				"key": "fixed",
				"name": "Fixed",
				"version": 1,
				"description": "Issue resolved",
				"order": 1
			},
			{
				"self": "https://api.tracker.yandex.net/v3/resolutions/2",
				"id": "2",
				"key": "wontFix",
				"name": "Won't fix",
				"version": 1,
				"order": 2
			}
		]`)
	})

	ctx := context.Background()
	resolutions, _, err := client.Resolutions.List(ctx)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	if len(resolutions) != 2 {
		t.Fatalf("List returned %d resolutions, want 2", len(resolutions))
	}

	want := &Resolution{
		Self:        Ptr("https://api.tracker.yandex.net/v3/resolutions/1"),
		ID:          Ptr("1"),
		Key:         Ptr("fixed"),
		Name:        Ptr("Fixed"),
		Version:     Ptr(1),
		Description: Ptr("Issue resolved"),
		Order:       Ptr(1),
	}

	if !reflect.DeepEqual(resolutions[0], want) {
		t.Errorf("List returned %+v, want %+v", resolutions[0], want)
	}
}
