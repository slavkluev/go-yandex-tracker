package tracker

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
)

func TestQueuesService_UpdatePermissions(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("PATCH /v3/queues/{key}/permissions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testHeader(t, r, "Content-Type", "application/json")

		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if _, ok := body["create"]; !ok {
			t.Error("Request body missing 'create' field")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(QueuePermissions{
			Self:    Ptr("https://api.tracker.yandex.net/v3/queues/TEST/permissions"),
			Version: Ptr(2),
			Create: &PermissionGroup{
				Users: []*PermissionUser{
					{Self: Ptr("https://api.tracker.yandex.net/v3/users/123"), ID: Ptr("123"), Display: Ptr("John Doe")},
					{Self: Ptr("https://api.tracker.yandex.net/v3/users/789"), ID: Ptr("789"), Display: Ptr("New User")},
				},
			},
		})
	})

	perms, _, err := client.Queues.UpdatePermissions(ctx, "TEST", &QueuePermissionsUpdateRequest{
		Create: &PermissionUpdate{
			Add: &PermissionSubjects{
				Users: []string{"789"},
			},
		},
	})
	if err != nil {
		t.Fatalf("Queues.UpdatePermissions returned error: %v", err)
	}

	want := &QueuePermissions{
		Self:    Ptr("https://api.tracker.yandex.net/v3/queues/TEST/permissions"),
		Version: Ptr(2),
		Create: &PermissionGroup{
			Users: []*PermissionUser{
				{Self: Ptr("https://api.tracker.yandex.net/v3/users/123"), ID: Ptr("123"), Display: Ptr("John Doe")},
				{Self: Ptr("https://api.tracker.yandex.net/v3/users/789"), ID: Ptr("789"), Display: Ptr("New User")},
			},
		},
	}

	if !reflect.DeepEqual(perms, want) {
		t.Errorf("Queues.UpdatePermissions returned %+v, want %+v", perms, want)
	}
}
