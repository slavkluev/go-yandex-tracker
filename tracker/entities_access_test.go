package tracker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestEntitiesService_GetAccess(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/entities/project/{id}/extendedPermissions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"acl": {
				"READ": {
					"users": [{"self":"https://api.tracker.yandex.net/v2/users/111","id":"111"}],
					"groups": [{"id":"1","display":"Group 1"}],
					"roles": ["author"]
				},
				"WRITE": {
					"users": [],
					"groups": [],
					"roles": ["assignee"]
				},
				"GRANT": {
					"users": [{"self":"https://api.tracker.yandex.net/v2/users/111","id":"111"}],
					"groups": [],
					"roles": []
				}
			},
			"permissionSources": [{"self":"https://api.tracker.yandex.net/v3/entities/portfolio/parent-1","id":"parent-1","display":"Parent Portfolio"}],
			"parentEntities": {
				"primary": {"self":"https://api.tracker.yandex.net/v3/entities/portfolio/parent-1","id":"parent-1","display":"Parent Portfolio"}
			}
		}`)
	})

	access, resp, err := client.Entities.GetAccess(ctx, EntityTypeProject, "1")
	if err != nil {
		t.Fatalf("Entities.GetAccess returned error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	// Verify ACL is not nil
	if access.ACL == nil {
		t.Fatal("access.ACL is nil, want non-nil")
	}

	// Verify READ entry
	if access.ACL.Read == nil {
		t.Fatal("access.ACL.Read is nil, want non-nil")
	}
	if got := len(access.ACL.Read.Users); got != 1 {
		t.Errorf("len(access.ACL.Read.Users) = %d, want 1", got)
	}
	if got := len(access.ACL.Read.Groups); got != 1 {
		t.Errorf("len(access.ACL.Read.Groups) = %d, want 1", got)
	}
	if got := len(access.ACL.Read.Roles); got != 1 {
		t.Errorf("len(access.ACL.Read.Roles) = %d, want 1", got)
	}
	if access.ACL.Read.Roles[0] != "author" {
		t.Errorf("access.ACL.Read.Roles[0] = %q, want %q", access.ACL.Read.Roles[0], "author")
	}

	// Verify WRITE entry
	if access.ACL.Write == nil {
		t.Fatal("access.ACL.Write is nil, want non-nil")
	}
	if got := len(access.ACL.Write.Roles); got != 1 {
		t.Errorf("len(access.ACL.Write.Roles) = %d, want 1", got)
	}
	if access.ACL.Write.Roles[0] != "assignee" {
		t.Errorf("access.ACL.Write.Roles[0] = %q, want %q", access.ACL.Write.Roles[0], "assignee")
	}

	// Verify GRANT entry
	if access.ACL.Grant == nil {
		t.Fatal("access.ACL.Grant is nil, want non-nil")
	}
	if got := len(access.ACL.Grant.Users); got != 1 {
		t.Errorf("len(access.ACL.Grant.Users) = %d, want 1", got)
	}

	// Verify permissionSources
	if got := len(access.PermissionSources); got != 1 {
		t.Errorf("len(access.PermissionSources) = %d, want 1", got)
	}

	// Verify parentEntities
	if access.ParentEntities == nil {
		t.Fatal("access.ParentEntities is nil, want non-nil")
	}
	if access.ParentEntities.Primary == nil {
		t.Fatal("access.ParentEntities.Primary is nil, want non-nil")
	}
	if got := *access.ParentEntities.Primary.ID; got != "parent-1" {
		t.Errorf("access.ParentEntities.Primary.ID = %q, want %q", got, "parent-1")
	}
}

func TestEntitiesService_UpdateAccess(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("PATCH /v3/entities/project/{id}/extendedPermissions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")

		// Decode and verify request body has grant.READ.users
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)

		acl, ok := body["acl"].(map[string]any)
		if !ok {
			t.Fatal("Request body missing 'acl' field")
		}
		grant, ok := acl["grant"].(map[string]any)
		if !ok {
			t.Fatal("Request body missing 'acl.grant' field")
		}
		read, ok := grant["READ"].(map[string]any)
		if !ok {
			t.Fatal("Request body missing 'acl.grant.READ' field")
		}
		users, ok := read["users"].([]any)
		if !ok {
			t.Fatal("Request body missing 'acl.grant.READ.users' field")
		}
		if len(users) != 1 || users[0] != "user-1" {
			t.Errorf("Request body acl.grant.READ.users = %v, want [\"user-1\"]", users)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"acl": {
				"READ": {
					"users": [{"self":"https://api.tracker.yandex.net/v2/users/user-1","id":"user-1"}],
					"groups": [],
					"roles": []
				},
				"WRITE": {
					"users": [],
					"groups": [],
					"roles": []
				},
				"GRANT": {
					"users": [],
					"groups": [],
					"roles": []
				}
			}
		}`)
	})

	access, resp, err := client.Entities.UpdateAccess(ctx, EntityTypeProject, "1", &EntityAccessUpdateRequest{
		ACL: &EntityACLUpdate{
			Grant: &EntityACLPermissions{
				Read: &EntityACLSubjects{
					Users: []string{"user-1"},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("Entities.UpdateAccess returned error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	// Verify response parsed correctly
	if access.ACL == nil {
		t.Fatal("access.ACL is nil, want non-nil")
	}
	if access.ACL.Read == nil {
		t.Fatal("access.ACL.Read is nil, want non-nil")
	}
	if got := len(access.ACL.Read.Users); got != 1 {
		t.Errorf("len(access.ACL.Read.Users) = %d, want 1", got)
	}
	if got := *access.ACL.Read.Users[0].ID; got != "user-1" {
		t.Errorf("access.ACL.Read.Users[0].ID = %q, want %q", got, "user-1")
	}
}

func TestEntitiesService_UpdateAccess_Revoke(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("PATCH /v3/entities/project/{id}/extendedPermissions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")

		// Decode and verify request body has revoke.WRITE.users
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)

		acl, ok := body["acl"].(map[string]any)
		if !ok {
			t.Fatal("Request body missing 'acl' field")
		}
		revoke, ok := acl["revoke"].(map[string]any)
		if !ok {
			t.Fatal("Request body missing 'acl.revoke' field")
		}
		write, ok := revoke["WRITE"].(map[string]any)
		if !ok {
			t.Fatal("Request body missing 'acl.revoke.WRITE' field")
		}
		users, ok := write["users"].([]any)
		if !ok {
			t.Fatal("Request body missing 'acl.revoke.WRITE.users' field")
		}
		if len(users) != 1 || users[0] != "user-2" {
			t.Errorf("Request body acl.revoke.WRITE.users = %v, want [\"user-2\"]", users)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"acl": {
				"READ": {
					"users": [],
					"groups": [],
					"roles": []
				},
				"WRITE": {
					"users": [],
					"groups": [],
					"roles": []
				},
				"GRANT": {
					"users": [],
					"groups": [],
					"roles": []
				}
			}
		}`)
	})

	access, resp, err := client.Entities.UpdateAccess(ctx, EntityTypeProject, "1", &EntityAccessUpdateRequest{
		ACL: &EntityACLUpdate{
			Revoke: &EntityACLPermissions{
				Write: &EntityACLSubjects{
					Users: []string{"user-2"},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("Entities.UpdateAccess returned error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	// Verify response parsed correctly
	if access.ACL == nil {
		t.Fatal("access.ACL is nil, want non-nil")
	}

	want := &EntityAccess{
		ACL: &EntityACL{
			Read: &EntityACLEntry{
				Users:  []*User{},
				Groups: []*PermissionGroupRef{},
				Roles:  []string{},
			},
			Write: &EntityACLEntry{
				Users:  []*User{},
				Groups: []*PermissionGroupRef{},
				Roles:  []string{},
			},
			Grant: &EntityACLEntry{
				Users:  []*User{},
				Groups: []*PermissionGroupRef{},
				Roles:  []string{},
			},
		},
	}
	if !reflect.DeepEqual(access, want) {
		t.Errorf("Entities.UpdateAccess (revoke) = %+v, want %+v", access, want)
	}
}
