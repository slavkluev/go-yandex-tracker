package tracker

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestIssueTypesService_List(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/issuetypes", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[
			{
				"self": "https://api.tracker.yandex.net/v3/issuetypes/2",
				"id": "2",
				"key": "task",
				"name": "Task",
				"version": 1,
				"description": "General task",
				"deleted": false
			},
			{
				"self": "https://api.tracker.yandex.net/v3/issuetypes/1",
				"id": "1",
				"key": "bug",
				"name": "Bug",
				"version": 1
			}
		]`)
	})

	ctx := context.Background()
	issueTypes, _, err := client.IssueTypes.List(ctx)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	if len(issueTypes) != 2 {
		t.Fatalf("List returned %d issue types, want 2", len(issueTypes))
	}

	want := &IssueType{
		Self:        Ptr("https://api.tracker.yandex.net/v3/issuetypes/2"),
		ID:          Ptr("2"),
		Key:         Ptr("task"),
		Name:        Ptr("Task"),
		Version:     Ptr(1),
		Description: Ptr("General task"),
		Deleted:     Ptr(false),
	}

	if !reflect.DeepEqual(issueTypes[0], want) {
		t.Errorf("List returned %+v, want %+v", issueTypes[0], want)
	}
}
