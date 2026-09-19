package tracker

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

// workflowW21Steps is the steps array of the W21 reference example.
const workflowW21Steps = `[
	{
		"status": {"self": "https://api.tracker.yandex.net/v3/statuses/1", "id": "1", "key": "open", "display": "Open"},
		"actions": [
			{"id": "inProgress", "name": "Start progress",
			 "target": {"self": "https://api.tracker.yandex.net/v3/statuses/3", "id": "3", "key": "inProgress", "display": "In progress"}}
		]
	},
	{
		"status": {"self": "https://api.tracker.yandex.net/v3/statuses/3", "id": "3", "key": "inProgress", "display": "In progress"},
		"actions": [
			{"id": "close", "name": "Close",
			 "target": {"self": "https://api.tracker.yandex.net/v3/statuses/8", "id": "8", "key": "closed", "display": "Closed"}}
		]
	},
	{
		"status": {"self": "https://api.tracker.yandex.net/v3/statuses/8", "id": "8", "key": "closed", "display": "Closed"}
	}
]`

// workflowW21InitialAction is the initialAction of the W21 reference example.
const workflowW21InitialAction = `{
	"id": "open", "name": "Open",
	"target": {"self": "https://api.tracker.yandex.net/v3/statuses/1", "id": "1", "key": "open", "display": "Open"}
}`

func testStatus(id, key, display string) *Status {
	return &Status{
		Self:    Ptr("https://api.tracker.yandex.net/v3/statuses/" + id),
		ID:      Ptr(FlexString(id)),
		Key:     Ptr(key),
		Display: Ptr(display),
	}
}

func wantWorkflowW21Steps() []*WorkflowStep {
	return []*WorkflowStep{
		{
			Status: testStatus("1", "open", "Open"),
			Actions: []*WorkflowAction{
				{ID: Ptr(FlexString("inProgress")), Name: Ptr("Start progress"), Target: testStatus("3", "inProgress", "In progress")},
			},
		},
		{
			Status: testStatus("3", "inProgress", "In progress"),
			Actions: []*WorkflowAction{
				{ID: Ptr(FlexString("close")), Name: Ptr("Close"), Target: testStatus("8", "closed", "Closed")},
			},
		},
		{
			Status: testStatus("8", "closed", "Closed"),
		},
	}
}

func wantWorkflowW21InitialAction() *WorkflowAction {
	return &WorkflowAction{ID: Ptr(FlexString("open")), Name: Ptr("Open"), Target: testStatus("1", "open", "Open")}
}

func TestWorkflowsService_Get(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/workflows/{id}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		if got := r.PathValue("id"); got != "W21" {
			t.Errorf("workflow id = %q, want %q", got, "W21")
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"self": "https://api.tracker.yandex.net/v3/workflows/W21",
			"id": "W21",
			"name": "Design",
			"version": 1,
			"steps": %s,
			"initialAction": %s,
			"queue": {"self": "https://api.tracker.yandex.net/v3/queues/DESIGN", "id": "4", "key": "DESIGN", "display": "DESIGN"},
			"created": "2026-08-11T14:37:06.356+0000",
			"updated": "2026-08-11T14:37:06.356+0000",
			"createdBy": {"self": "https://api.tracker.yandex.net/v3/users/11********", "id": "11********", "display": "User Name", "cloudUid": "ajeppa7dgp53********", "passportUid": 1100000000},
			"updatedBy": {"self": "https://api.tracker.yandex.net/v3/users/11********", "id": "11********", "display": "User Name", "cloudUid": "ajeppa7dgp53********", "passportUid": 1100000000},
			"deleted": false,
			"type": "visual"
		}`, workflowW21Steps, workflowW21InitialAction)
	})

	workflow, _, err := client.Workflows.Get(ctx, "W21")
	if err != nil {
		t.Fatalf("Workflows.Get returned error: %v", err)
	}

	user := &User{
		Self:        Ptr("https://api.tracker.yandex.net/v3/users/11********"),
		ID:          Ptr(FlexString("11********")),
		Display:     Ptr("User Name"),
		CloudUID:    Ptr("ajeppa7dgp53********"),
		PassportUID: Ptr(1100000000),
	}
	want := &Workflow{
		Self:          Ptr("https://api.tracker.yandex.net/v3/workflows/W21"),
		ID:            Ptr(FlexString("W21")),
		Name:          Ptr("Design"),
		Version:       Ptr(FlexString("1")),
		Steps:         wantWorkflowW21Steps(),
		InitialAction: wantWorkflowW21InitialAction(),
		Queue: &Queue{
			Self:    Ptr("https://api.tracker.yandex.net/v3/queues/DESIGN"),
			ID:      Ptr(FlexString("4")),
			Key:     Ptr("DESIGN"),
			Display: Ptr("DESIGN"),
		},
		Created:   newTestTimestamp(t, "2026-08-11T14:37:06.356+0000"),
		Updated:   newTestTimestamp(t, "2026-08-11T14:37:06.356+0000"),
		CreatedBy: user,
		UpdatedBy: user,
		Deleted:   Ptr(false),
		Type:      Ptr("visual"),
	}

	if !reflect.DeepEqual(workflow, want) {
		t.Fatalf("Workflows.Get returned %+v, want %+v", workflow, want)
	}
	if got := workflow.Steps[2].Actions; got != nil {
		t.Errorf("terminal step Actions = %+v, want nil", got)
	}
}

func TestWorkflowsService_Get_WithoutOptionalFields(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/workflows/{id}", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"self": "https://api.tracker.yandex.net/v3/workflows/W21",
			"id": "W21",
			"name": "Design",
			"version": 1,
			"steps": %s,
			"initialAction": %s,
			"created": "2026-08-11T14:37:06.356+0000",
			"updated": "2026-08-11T14:37:06.356+0000",
			"deleted": false
		}`, workflowW21Steps, workflowW21InitialAction)
	})

	workflow, _, err := client.Workflows.Get(ctx, "W21")
	if err != nil {
		t.Fatalf("Workflows.Get returned error: %v", err)
	}

	want := &Workflow{
		Self:          Ptr("https://api.tracker.yandex.net/v3/workflows/W21"),
		ID:            Ptr(FlexString("W21")),
		Name:          Ptr("Design"),
		Version:       Ptr(FlexString("1")),
		Steps:         wantWorkflowW21Steps(),
		InitialAction: wantWorkflowW21InitialAction(),
		Created:       newTestTimestamp(t, "2026-08-11T14:37:06.356+0000"),
		Updated:       newTestTimestamp(t, "2026-08-11T14:37:06.356+0000"),
		Deleted:       Ptr(false),
	}

	if !reflect.DeepEqual(workflow, want) {
		t.Errorf("Workflows.Get returned %+v, want %+v", workflow, want)
	}
}

func TestWorkflowsService_Get_NotFound(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/workflows/{id}", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"errorMessages":["Object not found"],"errors":{}}`)
	})

	_, _, err := client.Workflows.Get(ctx, "W0")
	if err == nil {
		t.Fatal("Workflows.Get expected error for 404, got nil")
	}
	if !IsNotFound(err) {
		t.Errorf("IsNotFound(err) = false, want true")
	}
}
