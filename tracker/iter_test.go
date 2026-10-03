package tracker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"iter"
	"net/http"
	"reflect"
	"slices"
	"strings"
	"testing"
)

// fakePage is one response of a paged endpoint.
type fakePage struct {
	body       string
	totalPages string // X-Total-Pages; empty sends no header
	link       string // Link; empty sends no header
	status     int    // 0 sends 200
}

type pageLog struct {
	queries []string
	bodies  []string
}

// servePages answers the n-th request to pattern with the n-th page and logs the
// query and body of every request.
func servePages(t *testing.T, mux *http.ServeMux, pattern string, pages ...fakePage) *pageLog {
	t.Helper()
	log := &pageLog{}
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		log.queries = append(log.queries, r.URL.RawQuery)
		log.bodies = append(log.bodies, strings.TrimSpace(string(body)))

		n := len(log.queries)
		if n > len(pages) {
			t.Errorf("request %d (%s) after the last page", n, r.URL)
			http.Error(w, "no such page", http.StatusInternalServerError)
			return
		}

		p := pages[n-1]
		if p.totalPages != "" {
			w.Header().Set("X-Total-Pages", p.totalPages)
		}
		if p.link != "" {
			w.Header().Set("Link", p.link)
		}
		w.Header().Set("Content-Type", "application/json")
		if p.status != 0 {
			w.WriteHeader(p.status)
		}
		fmt.Fprint(w, p.body)
	})
	return log
}

// drain ranges over seq and returns the ID of every item and every error it
// yields. An error must come with a nil item.
func drain[T any](t *testing.T, seq iter.Seq2[*T, error]) ([]string, []error) {
	t.Helper()
	var ids []string
	var errs []error
	for item, err := range seq {
		if err != nil {
			if item != nil {
				t.Errorf("error %v came with item %+v, want nil", err, item)
			}
			errs = append(errs, err)
			continue
		}
		ids = append(ids, itemID(item))
	}
	return ids, errs
}

func itemID(item any) string {
	v := reflect.ValueOf(item)
	if v.IsNil() {
		return "<nil>"
	}
	id, _ := v.Elem().FieldByName("ID").Interface().(*FlexString)
	if id == nil {
		return ""
	}
	return string(*id)
}

func TestPageIter_StopsAfterTotalPages(t *testing.T) {
	client, mux := setup(t)
	log := servePages(t, mux, "GET /v3/statuses",
		fakePage{body: `[{"id":"a"},{"id":"b"}]`, totalPages: "3", link: `<https://api.tracker.yandex.net/v3/statuses?page=2&perPage=2>; rel="next"`},
		fakePage{body: `[{"id":"c"},{"id":"d"}]`, totalPages: "3", link: `<https://api.tracker.yandex.net/v3/statuses?page=3&perPage=2>; rel="next"`},
		fakePage{body: `[{"id":"e"},{"id":"f"}]`, totalPages: "3"},
	)

	opts := &ListOptions{PerPage: 2}
	ids, errs := drain(t, client.Statuses.ListIter(context.Background(), opts))

	if len(errs) != 0 {
		t.Fatalf("ListIter yielded errors %v", errs)
	}
	if want := []string{"a", "b", "c", "d", "e", "f"}; !slices.Equal(ids, want) {
		t.Errorf("ListIter yielded %v, want %v", ids, want)
	}
	if want := []string{"perPage=2", "page=2&perPage=2", "page=3&perPage=2"}; !slices.Equal(log.queries, want) {
		t.Errorf("requests %q, want %q", log.queries, want)
	}
	if *opts != (ListOptions{PerPage: 2}) {
		t.Errorf("ListIter changed opts to %+v", *opts)
	}
}

func TestPageIter_WithoutTotalPagesStopsOnEmptyPage(t *testing.T) {
	client, mux := setup(t)
	log := servePages(t, mux, "GET /v3/statuses",
		fakePage{body: `[{"id":"a"},{"id":"b"}]`},
		fakePage{body: `[]`},
	)

	ids, errs := drain(t, client.Statuses.ListIter(context.Background(), nil))

	if len(errs) != 0 {
		t.Fatalf("ListIter yielded errors %v", errs)
	}
	if want := []string{"a", "b"}; !slices.Equal(ids, want) {
		t.Errorf("ListIter yielded %v, want %v", ids, want)
	}
	if want := []string{"", "page=2"}; !slices.Equal(log.queries, want) {
		t.Errorf("requests %q, want %q", log.queries, want)
	}
}

func TestPageIter_StartsAtCallerPage(t *testing.T) {
	client, mux := setup(t)
	log := servePages(t, mux, "GET /v3/statuses",
		fakePage{body: `[{"id":"c"}]`, totalPages: "3"},
		fakePage{body: `[{"id":"e"}]`, totalPages: "3"},
	)

	ids, errs := drain(t, client.Statuses.ListIter(context.Background(), &ListOptions{Page: 2}))

	if len(errs) != 0 {
		t.Fatalf("ListIter yielded errors %v", errs)
	}
	if want := []string{"c", "e"}; !slices.Equal(ids, want) {
		t.Errorf("ListIter yielded %v, want %v", ids, want)
	}
	if want := []string{"page=2", "page=3"}; !slices.Equal(log.queries, want) {
		t.Errorf("requests %q, want %q", log.queries, want)
	}
}

func TestCursorIter_StopsOnEmptyPage(t *testing.T) {
	client, mux := setup(t)
	// Tracker sends rel="next" and X-Total-Pages: 1 even when more comments follow.
	log := servePages(t, mux, "GET /v3/issues/TEST-1/comments",
		fakePage{body: `[{"id":"a"},{"id":"b"}]`, totalPages: "1", link: `<https://api.tracker.yandex.net/v3/issues/TEST-1/comments?id=b&perPage=2>; rel="next"`},
		fakePage{body: `[{"id":"c"}]`, totalPages: "1", link: `<https://api.tracker.yandex.net/v3/issues/TEST-1/comments?id=c&perPage=2>; rel="next"`},
		fakePage{body: `[]`},
	)

	opts := &CommentListOptions{PerPage: 2}
	ids, errs := drain(t, client.Issues.ListCommentsIter(context.Background(), "TEST-1", opts))

	if len(errs) != 0 {
		t.Fatalf("ListCommentsIter yielded errors %v", errs)
	}
	if want := []string{"a", "b", "c"}; !slices.Equal(ids, want) {
		t.Errorf("ListCommentsIter yielded %v, want %v", ids, want)
	}
	if want := []string{"perPage=2", "id=b&perPage=2", "id=c&perPage=2"}; !slices.Equal(log.queries, want) {
		t.Errorf("requests %q, want %q", log.queries, want)
	}
	if *opts != (CommentListOptions{PerPage: 2}) {
		t.Errorf("ListCommentsIter changed opts to %+v", *opts)
	}
}

func TestCursorIter_StuckCursorYieldsOneError(t *testing.T) {
	tests := []struct {
		name      string
		cursor    string
		body      string
		wantIDs   []string
		wantQuery string
	}{
		{name: "last item has no ID", body: `[{"id":"a"},{}]`, wantIDs: []string{"a", ""}, wantQuery: ""},
		{name: "last item repeats the cursor", cursor: "b", body: `[{"id":"a"},{"id":"b"}]`, wantIDs: []string{"a", "b"}, wantQuery: "id=b"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, mux := setup(t)
			log := servePages(t, mux, "GET /v3/issues/TEST-1/comments", fakePage{body: tt.body})

			ids, errs := drain(t, client.Issues.ListCommentsIter(context.Background(), "TEST-1", &CommentListOptions{ID: tt.cursor}))

			if !slices.Equal(ids, tt.wantIDs) {
				t.Errorf("ListCommentsIter yielded %v, want %v", ids, tt.wantIDs)
			}
			if len(errs) != 1 {
				t.Errorf("ListCommentsIter yielded errors %v, want one", errs)
			}
			if want := []string{tt.wantQuery}; !slices.Equal(log.queries, want) {
				t.Errorf("requests %q, want %q", log.queries, want)
			}
		})
	}
}

func TestIter_BreakSendsNoMoreRequests(t *testing.T) {
	client, mux := setup(t)
	statuses := servePages(t, mux, "GET /v3/statuses",
		fakePage{body: `[{"id":"a"},{"id":"b"}]`, totalPages: "2"},
	)
	comments := servePages(t, mux, "GET /v3/issues/TEST-1/comments",
		fakePage{body: `[{"id":"a"},{"id":"b"}]`},
	)

	var got []string
	for status, err := range client.Statuses.ListIter(context.Background(), nil) {
		if err != nil {
			t.Fatalf("ListIter yielded error %v", err)
		}
		got = append(got, itemID(status))
		if len(got) == 1 {
			break
		}
	}
	for comment, err := range client.Issues.ListCommentsIter(context.Background(), "TEST-1", nil) {
		if err != nil {
			t.Fatalf("ListCommentsIter yielded error %v", err)
		}
		got = append(got, itemID(comment))
		if len(got) == 2 {
			break
		}
	}

	if want := []string{"a", "a"}; !slices.Equal(got, want) {
		t.Errorf("iterators yielded %v, want %v", got, want)
	}
	if len(statuses.queries) != 1 || len(comments.queries) != 1 {
		t.Errorf("requests: statuses %q, comments %q, want one each", statuses.queries, comments.queries)
	}
}

func TestIter_ErrorEndsIteration(t *testing.T) {
	failure := fakePage{body: `{"errorMessages":["boom"]}`, status: http.StatusInternalServerError}
	tests := []struct {
		name    string
		pattern string
		first   fakePage
		run     func(t *testing.T, c *Client) ([]string, []error)
	}{
		{
			name:    "page-numbered",
			pattern: "GET /v3/statuses",
			first:   fakePage{body: `[{"id":"a"},{"id":"b"}]`, totalPages: "3"},
			run: func(t *testing.T, c *Client) ([]string, []error) {
				t.Helper()
				return drain(t, c.Statuses.ListIter(context.Background(), nil))
			},
		},
		{
			name:    "cursor",
			pattern: "GET /v3/issues/TEST-1/comments",
			first:   fakePage{body: `[{"id":"a"},{"id":"b"}]`},
			run: func(t *testing.T, c *Client) ([]string, []error) {
				t.Helper()
				return drain(t, c.Issues.ListCommentsIter(context.Background(), "TEST-1", nil))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, mux := setup(t)
			log := servePages(t, mux, tt.pattern, tt.first, failure)

			ids, errs := tt.run(t, client)

			if want := []string{"a", "b"}; !slices.Equal(ids, want) {
				t.Errorf("iterator yielded %v, want %v", ids, want)
			}
			var errResp *ErrorResponse
			if len(errs) != 1 || !errors.As(errs[0], &errResp) || errResp.Response.StatusCode != http.StatusInternalServerError {
				t.Errorf("iterator yielded errors %v, want one 500", errs)
			}
			if len(log.queries) != 2 {
				t.Errorf("requests %q, want two", log.queries)
			}
		})
	}
}

// drainWith ranges over list(opts) and fails the test if that changed opts.
func drainWith[O, T any](t *testing.T, opts *O, list func(*O) iter.Seq2[*T, error]) ([]string, []error) {
	t.Helper()
	before := *opts
	ids, errs := drain(t, list(opts))
	if !reflect.DeepEqual(*opts, before) {
		t.Errorf("iterator changed opts from %+v to %+v", before, *opts)
	}
	return ids, errs
}

func TestIterators_PageWithCallerOptions(t *testing.T) {
	ctx := context.Background()
	numbered := []fakePage{
		{body: `[{"id":"a"},{"id":"b"}]`, totalPages: "2"},
		{body: `[{"id":"c"}]`, totalPages: "2"},
	}
	cursor := []fakePage{
		{body: `[{"id":"a"},{"id":"b"}]`},
		{body: `[{"id":"c"}]`},
		{body: `[]`},
	}

	tests := []struct {
		name     string
		pattern  string
		pages    []fakePage
		run      func(t *testing.T, c *Client) ([]string, []error)
		want     []string
		wantBody string
	}{
		{
			name:    "Statuses.ListIter",
			pattern: "GET /v3/statuses",
			pages:   numbered,
			run: func(t *testing.T, c *Client) ([]string, []error) {
				t.Helper()
				return drainWith(t, &ListOptions{PerPage: 2}, func(o *ListOptions) iter.Seq2[*Status, error] {
					return c.Statuses.ListIter(ctx, o)
				})
			},
			want: []string{"perPage=2", "page=2&perPage=2"},
		},
		{
			name:    "Resolutions.ListIter",
			pattern: "GET /v3/resolutions",
			pages:   numbered,
			run: func(t *testing.T, c *Client) ([]string, []error) {
				t.Helper()
				return drainWith(t, &ListOptions{PerPage: 2}, func(o *ListOptions) iter.Seq2[*Resolution, error] {
					return c.Resolutions.ListIter(ctx, o)
				})
			},
			want: []string{"perPage=2", "page=2&perPage=2"},
		},
		{
			name:    "IssueTypes.ListIter",
			pattern: "GET /v3/issuetypes",
			pages:   numbered,
			run: func(t *testing.T, c *Client) ([]string, []error) {
				t.Helper()
				return drainWith(t, &ListOptions{PerPage: 2}, func(o *ListOptions) iter.Seq2[*IssueType, error] {
					return c.IssueTypes.ListIter(ctx, o)
				})
			},
			want: []string{"perPage=2", "page=2&perPage=2"},
		},
		{
			name:    "Priorities.ListIter",
			pattern: "GET /v3/priorities",
			pages:   numbered,
			run: func(t *testing.T, c *Client) ([]string, []error) {
				t.Helper()
				opts := &PriorityListOptions{ListOptions: ListOptions{PerPage: 2}, Localized: Ptr(true)}
				return drainWith(t, opts, func(o *PriorityListOptions) iter.Seq2[*Priority, error] {
					return c.Priorities.ListIter(ctx, o)
				})
			},
			want: []string{"localized=true&perPage=2", "localized=true&page=2&perPage=2"},
		},
		{
			name:    "Queues.ListIter",
			pattern: "GET /v3/queues",
			pages:   numbered,
			run: func(t *testing.T, c *Client) ([]string, []error) {
				t.Helper()
				opts := &QueueListOptions{ListOptions: ListOptions{PerPage: 2}, Expand: "projects"}
				return drainWith(t, opts, func(o *QueueListOptions) iter.Seq2[*Queue, error] {
					return c.Queues.ListIter(ctx, o)
				})
			},
			want: []string{"expand=projects&perPage=2", "expand=projects&page=2&perPage=2"},
		},
		{
			name:    "Users.ListIter",
			pattern: "GET /v3/users",
			pages:   numbered,
			run: func(t *testing.T, c *Client) ([]string, []error) {
				t.Helper()
				return drainWith(t, &UserListOptions{ListOptions: ListOptions{PerPage: 2}}, func(o *UserListOptions) iter.Seq2[*User, error] {
					return c.Users.ListIter(ctx, o)
				})
			},
			want: []string{"perPage=2", "page=2&perPage=2"},
		},
		{
			name:    "Issues.SearchIter",
			pattern: "POST /v3/issues/_search",
			pages:   numbered,
			run: func(t *testing.T, c *Client) ([]string, []error) {
				t.Helper()
				search := &IssueSearchRequest{Queue: Ptr("TEST")}
				opts := &IssueSearchOptions{ListOptions: ListOptions{PerPage: 2}, Expand: "transitions"}
				return drainWith(t, opts, func(o *IssueSearchOptions) iter.Seq2[*Issue, error] {
					return c.Issues.SearchIter(ctx, search, o)
				})
			},
			want:     []string{"expand=transitions&perPage=2", "expand=transitions&page=2&perPage=2"},
			wantBody: `{"queue":"TEST"}`,
		},
		{
			name:    "Issues.GetChangelogIter",
			pattern: "GET /v3/issues/TEST-1/changelog",
			pages:   cursor,
			run: func(t *testing.T, c *Client) ([]string, []error) {
				t.Helper()
				opts := &ChangelogOptions{PerPage: 2, Field: "status", Type: "IssueUpdated"}
				return drainWith(t, opts, func(o *ChangelogOptions) iter.Seq2[*Changelog, error] {
					return c.Issues.GetChangelogIter(ctx, "TEST-1", o)
				})
			},
			want: []string{
				"field=status&perPage=2&type=IssueUpdated",
				"field=status&id=b&perPage=2&type=IssueUpdated",
				"field=status&id=c&perPage=2&type=IssueUpdated",
			},
		},
		{
			name:    "Issues.ListCommentsIter",
			pattern: "GET /v3/issues/TEST-1/comments",
			pages:   cursor,
			run: func(t *testing.T, c *Client) ([]string, []error) {
				t.Helper()
				return drainWith(t, &CommentListOptions{PerPage: 2}, func(o *CommentListOptions) iter.Seq2[*Comment, error] {
					return c.Issues.ListCommentsIter(ctx, "TEST-1", o)
				})
			},
			want: []string{"perPage=2", "id=b&perPage=2", "id=c&perPage=2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, mux := setup(t)
			log := servePages(t, mux, tt.pattern, tt.pages...)

			ids, errs := tt.run(t, client)

			if len(errs) != 0 {
				t.Fatalf("iterator yielded errors %v", errs)
			}
			if want := []string{"a", "b", "c"}; !slices.Equal(ids, want) {
				t.Errorf("iterator yielded %v, want %v", ids, want)
			}
			if !slices.Equal(log.queries, tt.want) {
				t.Errorf("requests %q, want %q", log.queries, tt.want)
			}
			for i, body := range log.bodies {
				if body != tt.wantBody {
					t.Errorf("request %d body %s, want %s", i+1, body, tt.wantBody)
				}
			}
		})
	}
}
