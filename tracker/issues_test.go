package tracker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestIssuesService_Create(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/issues/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["summary"] != "Test issue" {
			t.Errorf("Request body summary = %v, want %q", body["summary"], "Test issue")
		}
		if body["queue"] != "TEST" {
			t.Errorf("Request body queue = %v, want %q", body["queue"], "TEST")
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"self":"https://api.tracker.yandex.net/v3/issues/TEST-1","id":"1","key":"TEST-1","summary":"Test issue"}`)
	})

	issue, resp, err := client.Issues.Create(ctx, &IssueRequest{
		Summary: Ptr("Test issue"),
		Queue:   Ptr("TEST"),
	})
	if err != nil {
		t.Fatalf("Issues.Create returned error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	want := &Issue{
		Self:    Ptr("https://api.tracker.yandex.net/v3/issues/TEST-1"),
		ID:      Ptr(FlexString("1")),
		Key:     Ptr("TEST-1"),
		Summary: Ptr("Test issue"),
	}
	if !reflect.DeepEqual(issue, want) {
		t.Errorf("Issues.Create = %+v, want %+v", issue, want)
	}
}

func TestIssuesService_Create_InvalidJSON(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/issues/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"errorMessages":["Invalid request"],"errors":{}}`)
	})

	_, _, err := client.Issues.Create(ctx, &IssueRequest{
		Summary: Ptr("Test"),
	})
	if err == nil {
		t.Fatal("Issues.Create expected error, got nil")
	}
}

func TestIssuesService_Get(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/issues/{key}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"self":"https://api.tracker.yandex.net/v3/issues/QUEUE-1","id":"1","key":"QUEUE-1","summary":"Test issue","description":"A test issue"}`)
	})

	issue, _, err := client.Issues.Get(ctx, "QUEUE-1", nil)
	if err != nil {
		t.Fatalf("Issues.Get returned error: %v", err)
	}

	want := &Issue{
		Self:        Ptr("https://api.tracker.yandex.net/v3/issues/QUEUE-1"),
		ID:          Ptr(FlexString("1")),
		Key:         Ptr("QUEUE-1"),
		Summary:     Ptr("Test issue"),
		Description: Ptr("A test issue"),
	}
	if !reflect.DeepEqual(issue, want) {
		t.Errorf("Issues.Get = %+v, want %+v", issue, want)
	}
}

func TestIssuesService_Get_WithExpand(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/issues/{key}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		expand := r.URL.Query().Get("expand")
		if expand != "transitions" {
			t.Errorf("expand query param = %q, want %q", expand, "transitions")
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"key":"QUEUE-1"}`)
	})

	_, _, err := client.Issues.Get(ctx, "QUEUE-1", &IssueGetOptions{
		Expand: "transitions",
	})
	if err != nil {
		t.Fatalf("Issues.Get returned error: %v", err)
	}
}

func TestIssuesService_Get_NotFound(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/issues/{key}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"errorMessages":["Object not found"],"errors":{}}`)
	})

	_, _, err := client.Issues.Get(ctx, "NONEXIST-1", nil)
	if err == nil {
		t.Fatal("Issues.Get expected error for 404, got nil")
	}
	if !IsNotFound(err) {
		t.Errorf("IsNotFound(err) = false, want true")
	}
}

func TestIssuesService_Edit(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("PATCH /v3/issues/{key}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")

		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["summary"] != "Updated title" {
			t.Errorf("Request body summary = %v, want %q", body["summary"], "Updated title")
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"key":"QUEUE-1","summary":"Updated title"}`)
	})

	issue, _, err := client.Issues.Edit(ctx, "QUEUE-1", &IssueRequest{
		Summary: Ptr("Updated title"),
	}, nil)
	if err != nil {
		t.Fatalf("Issues.Edit returned error: %v", err)
	}

	want := &Issue{
		Key:     Ptr("QUEUE-1"),
		Summary: Ptr("Updated title"),
	}
	if !reflect.DeepEqual(issue, want) {
		t.Errorf("Issues.Edit = %+v, want %+v", issue, want)
	}
}

func TestIssuesService_Edit_WithVersion(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("PATCH /v3/issues/{key}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")

		version := r.URL.Query().Get("version")
		if version != "2" {
			t.Errorf("version query param = %q, want %q", version, "2")
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"key":"QUEUE-1"}`)
	})

	_, _, err := client.Issues.Edit(ctx, "QUEUE-1", &IssueRequest{
		Summary: Ptr("Updated"),
	}, &IssueEditOptions{Version: 2})
	if err != nil {
		t.Fatalf("Issues.Edit returned error: %v", err)
	}
}

func TestIssuesService_Search(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/issues/_search", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["query"] != "queue: TEST" {
			t.Errorf("Request body query = %v, want %q", body["query"], "queue: TEST")
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Total-Count", "2")
		fmt.Fprint(w, `[{"key":"QUEUE-1"},{"key":"QUEUE-2"}]`)
	})

	issues, resp, err := client.Issues.Search(ctx, &IssueSearchRequest{
		Query: Ptr("queue: TEST"),
	}, nil)
	if err != nil {
		t.Fatalf("Issues.Search returned error: %v", err)
	}
	if len(issues) != 2 {
		t.Fatalf("Issues.Search returned %d issues, want 2", len(issues))
	}
	if got := *issues[0].Key; got != "QUEUE-1" {
		t.Errorf("issues[0].Key = %q, want %q", got, "QUEUE-1")
	}
	if got := *issues[1].Key; got != "QUEUE-2" {
		t.Errorf("issues[1].Key = %q, want %q", got, "QUEUE-2")
	}
	if resp.TotalCount != 2 {
		t.Errorf("resp.TotalCount = %d, want 2", resp.TotalCount)
	}
}

func TestIssuesService_Search_WithPagination(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/issues/_search", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		page := r.URL.Query().Get("page")
		if page != "2" {
			t.Errorf("page query param = %q, want %q", page, "2")
		}
		perPage := r.URL.Query().Get("perPage")
		if perPage != "10" {
			t.Errorf("perPage query param = %q, want %q", perPage, "10")
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[]`)
	})

	_, _, err := client.Issues.Search(ctx, &IssueSearchRequest{
		Query: Ptr("queue: TEST"),
	}, &IssueSearchOptions{
		ListOptions: ListOptions{Page: 2, PerPage: 10},
	})
	if err != nil {
		t.Fatalf("Issues.Search returned error: %v", err)
	}
}

func TestIssuesService_ScrollSearch(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/issues/_search", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		scrollType := r.URL.Query().Get("scrollType")
		if scrollType != "sorted" {
			t.Errorf("scrollType query param = %q, want %q", scrollType, "sorted")
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Scroll-Id", "scroll123")
		w.Header().Set("X-Scroll-Token", "token456")
		fmt.Fprint(w, `[{"key":"QUEUE-1"}]`)
	})

	issues, resp, err := client.Issues.ScrollSearch(ctx, &IssueSearchRequest{
		Query: Ptr("queue: TEST"),
	}, &ScrollSearchOptions{
		ScrollType: "sorted",
	})
	if err != nil {
		t.Fatalf("Issues.ScrollSearch returned error: %v", err)
	}
	if len(issues) != 1 {
		t.Fatalf("Issues.ScrollSearch returned %d issues, want 1", len(issues))
	}
	if resp.ScrollID != "scroll123" {
		t.Errorf("resp.ScrollID = %q, want %q", resp.ScrollID, "scroll123")
	}
	if resp.ScrollToken != "token456" {
		t.Errorf("resp.ScrollToken = %q, want %q", resp.ScrollToken, "token456")
	}
}

func TestIssuesService_ScrollNext(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/issues/_search", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		scrollID := r.URL.Query().Get("scrollId")
		if scrollID != "scroll123" {
			t.Errorf("scrollId query param = %q, want %q", scrollID, "scroll123")
		}
		scrollToken := r.URL.Query().Get("scrollToken")
		if scrollToken != "token456" {
			t.Errorf("scrollToken query param = %q, want %q", scrollToken, "token456")
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Scroll-Id", "scroll123")
		w.Header().Set("X-Scroll-Token", "token789")
		fmt.Fprint(w, `[{"key":"QUEUE-2"}]`)
	})

	issues, resp, err := client.Issues.ScrollNext(ctx, &IssueSearchRequest{
		Query: Ptr("queue: TEST"),
	}, "scroll123", "token456")
	if err != nil {
		t.Fatalf("Issues.ScrollNext returned error: %v", err)
	}
	if len(issues) != 1 {
		t.Fatalf("Issues.ScrollNext returned %d issues, want 1", len(issues))
	}
	if got := *issues[0].Key; got != "QUEUE-2" {
		t.Errorf("issues[0].Key = %q, want %q", got, "QUEUE-2")
	}
	if resp.ScrollToken != "token789" {
		t.Errorf("resp.ScrollToken = %q, want %q", resp.ScrollToken, "token789")
	}
}

func TestIssuesService_Count(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/issues/_count", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["query"] != "queue: TEST" {
			t.Errorf("Request body query = %v, want %q", body["query"], "queue: TEST")
		}

		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, `42`)
	})

	count, _, err := client.Issues.Count(ctx, &IssueSearchRequest{
		Query: Ptr("queue: TEST"),
	})
	if err != nil {
		t.Fatalf("Issues.Count returned error: %v", err)
	}
	if count != 42 {
		t.Errorf("Issues.Count = %d, want 42", count)
	}
}

func TestIssuesService_Count_Zero(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/issues/_count", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, `0`)
	})

	count, _, err := client.Issues.Count(ctx, &IssueSearchRequest{
		Query: Ptr("queue: TEST"),
	})
	if err != nil {
		t.Fatalf("Issues.Count returned error: %v", err)
	}
	if count != 0 {
		t.Errorf("Issues.Count = %d, want 0", count)
	}
}

func TestIssuesService_Move(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/issues/{key}/_move", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		queue := r.URL.Query().Get("queue")
		if queue != "NEW" {
			t.Errorf("queue query param = %q, want %q", queue, "NEW")
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"key":"NEW-1","summary":"Moved issue"}`)
	})

	issue, _, err := client.Issues.Move(ctx, "QUEUE-1", &IssueMoveOptions{
		Queue: "NEW",
	}, nil)
	if err != nil {
		t.Fatalf("Issues.Move returned error: %v", err)
	}
	if got := *issue.Key; got != "NEW-1" {
		t.Errorf("issue.Key = %q, want %q", got, "NEW-1")
	}
	if got := *issue.Summary; got != "Moved issue" {
		t.Errorf("issue.Summary = %q, want %q", got, "Moved issue")
	}
}

// ctx is a shared context for test convenience.
var ctx = context.Background()
