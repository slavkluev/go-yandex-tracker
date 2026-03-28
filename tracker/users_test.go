package tracker

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestUsersService_Myself(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/myself", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/myself",
			"uid": 12345678,
			"login": "user1",
			"trackerUid": 12345678,
			"passportUid": 12345678,
			"firstName": "John",
			"lastName": "Doe",
			"display": "John Doe",
			"email": "john@example.com",
			"external": false,
			"hasLicense": true,
			"dismissed": false,
			"useNewFilters": true,
			"disableNotifications": false,
			"welcomeMailSent": true
		}`)
	})

	ctx := context.Background()
	user, _, err := client.Users.Myself(ctx)
	if err != nil {
		t.Fatalf("Myself returned error: %v", err)
	}

	want := &User{
		Self:                 Ptr("https://api.tracker.yandex.net/v3/myself"),
		UID:                  Ptr(12345678),
		Login:                Ptr("user1"),
		TrackerUID:           Ptr(12345678),
		PassportUID:          Ptr(12345678),
		FirstName:            Ptr("John"),
		LastName:             Ptr("Doe"),
		Display:              Ptr("John Doe"),
		Email:                Ptr("john@example.com"),
		External:             Ptr(false),
		HasLicense:           Ptr(true),
		Dismissed:            Ptr(false),
		UseNewFilters:        Ptr(true),
		DisableNotifications: Ptr(false),
		WelcomeMailSent:      Ptr(true),
	}

	if !reflect.DeepEqual(user, want) {
		t.Errorf("Myself returned %+v, want %+v", user, want)
	}
}

func TestUsersService_Get(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		id := r.PathValue("id")
		if id != "user123" {
			t.Errorf("URL path id = %q, want %q", id, "user123")
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/users/user123",
			"uid": 87654321,
			"login": "user123",
			"trackerUid": 87654321,
			"passportUid": 87654321,
			"firstName": "Jane",
			"lastName": "Smith",
			"display": "Jane Smith",
			"email": "jane@example.com",
			"external": false,
			"hasLicense": true,
			"dismissed": false
		}`)
	})

	ctx := context.Background()
	user, _, err := client.Users.Get(ctx, "user123")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}

	want := &User{
		Self:        Ptr("https://api.tracker.yandex.net/v3/users/user123"),
		UID:         Ptr(87654321),
		Login:       Ptr("user123"),
		TrackerUID:  Ptr(87654321),
		PassportUID: Ptr(87654321),
		FirstName:   Ptr("Jane"),
		LastName:    Ptr("Smith"),
		Display:     Ptr("Jane Smith"),
		Email:       Ptr("jane@example.com"),
		External:    Ptr(false),
		HasLicense:  Ptr(true),
		Dismissed:   Ptr(false),
	}

	if !reflect.DeepEqual(user, want) {
		t.Errorf("Get returned %+v, want %+v", user, want)
	}
}

func TestUsersService_List(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[
			{
				"self": "https://api.tracker.yandex.net/v3/users/user1",
				"uid": 11111111,
				"login": "user1",
				"display": "User One"
			},
			{
				"self": "https://api.tracker.yandex.net/v3/users/user2",
				"uid": 22222222,
				"login": "user2",
				"display": "User Two"
			}
		]`)
	})

	ctx := context.Background()
	users, _, err := client.Users.List(ctx, nil)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	if len(users) != 2 {
		t.Fatalf("List returned %d users, want 2", len(users))
	}

	want := []*User{
		{
			Self:    Ptr("https://api.tracker.yandex.net/v3/users/user1"),
			UID:     Ptr(11111111),
			Login:   Ptr("user1"),
			Display: Ptr("User One"),
		},
		{
			Self:    Ptr("https://api.tracker.yandex.net/v3/users/user2"),
			UID:     Ptr(22222222),
			Login:   Ptr("user2"),
			Display: Ptr("User Two"),
		},
	}

	if !reflect.DeepEqual(users, want) {
		t.Errorf("List returned %+v, want %+v", users, want)
	}
}

func TestUsersService_List_Pagination(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		if got := r.URL.Query().Get("page"); got != "2" {
			t.Errorf("query param page = %q, want %q", got, "2")
		}
		if got := r.URL.Query().Get("perPage"); got != "10" {
			t.Errorf("query param perPage = %q, want %q", got, "10")
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[]`)
	})

	ctx := context.Background()
	_, _, err := client.Users.List(ctx, &UserListOptions{
		ListOptions: ListOptions{Page: 2, PerPage: 10},
	})
	if err != nil {
		t.Fatalf("List with pagination returned error: %v", err)
	}
}
