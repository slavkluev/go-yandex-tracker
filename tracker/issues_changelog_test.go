package tracker

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestIssuesService_GetChangelog(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/issues/{key}/changelog", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{
			"id": "1",
			"type": "IssueUpdated",
			"transport": "front",
			"updatedAt": "2017-06-11T05:16:01.339+0000",
			"updatedBy": {"id": "user1"},
			"fields": [{
				"field": {
					"self": "https://api.tracker.yandex.net/v3/fields/status",
					"id": "status",
					"display": "Status"
				},
				"from": "open",
				"to": "closed"
			}],
			"comments": {
				"added": [{
					"self": "https://api.tracker.yandex.net/v3/issues/QUEUE-1/comments/10",
					"id": "10",
					"display": "Comment text"
				}]
			},
			"executedTriggers": [{
				"trigger": {
					"self": "https://api.tracker.yandex.net/v3/queues/QUEUE/triggers/29",
					"id": "29",
					"display": "Trigger-42"
				},
				"success": true,
				"message": "Success"
			}]
		}]`)
	})

	ctx := context.Background()
	changelog, _, err := client.Issues.GetChangelog(ctx, "QUEUE-1", nil)
	if err != nil {
		t.Fatalf("GetChangelog returned error: %v", err)
	}

	if len(changelog) != 1 {
		t.Fatalf("GetChangelog returned %d entries, want 1", len(changelog))
	}

	entry := changelog[0]
	if got := *entry.ID; got != "1" {
		t.Errorf("ID = %q, want %q", got, "1")
	}
	if got := *entry.Type; got != "IssueUpdated" {
		t.Errorf("Type = %q, want %q", got, "IssueUpdated")
	}
	if got := *entry.Transport; got != "front" {
		t.Errorf("Transport = %q, want %q", got, "front")
	}

	if len(entry.Fields) != 1 {
		t.Fatalf("Fields length = %d, want 1", len(entry.Fields))
	}
	if got := string(*entry.Fields[0].Field.ID); got != "status" {
		t.Errorf("Fields[0].Field.ID = %q, want %q", got, "status")
	}
	if got := *entry.Fields[0].Field.Display; got != "Status" {
		t.Errorf("Fields[0].Field.Display = %q, want %q", got, "Status")
	}

	if entry.Comments == nil {
		t.Fatal("Comments is nil")
	}
	if len(entry.Comments.Added) != 1 {
		t.Fatalf("Comments.Added length = %d, want 1", len(entry.Comments.Added))
	}
	if got := string(*entry.Comments.Added[0].ID); got != "10" {
		t.Errorf("Comments.Added[0].ID = %q, want %q", got, "10")
	}
	if got := *entry.Comments.Added[0].Display; got != "Comment text" {
		t.Errorf("Comments.Added[0].Display = %q, want %q", got, "Comment text")
	}

	if len(entry.ExecutedTriggers) != 1 {
		t.Fatalf("ExecutedTriggers length = %d, want 1", len(entry.ExecutedTriggers))
	}
	if got := string(*entry.ExecutedTriggers[0].Trigger.ID); got != "29" {
		t.Errorf("ExecutedTriggers[0].Trigger.ID = %q, want %q", got, "29")
	}
	if got := *entry.ExecutedTriggers[0].Success; got != true {
		t.Errorf("ExecutedTriggers[0].Success = %v, want true", got)
	}
	if got := *entry.ExecutedTriggers[0].Message; got != "Success" {
		t.Errorf("ExecutedTriggers[0].Message = %q, want %q", got, "Success")
	}
}

func TestIssuesService_GetChangelog_WithOptions(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/issues/{key}/changelog", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		gotID := r.URL.Query().Get("id")
		if gotID != "last123" {
			t.Errorf("query param id = %q, want %q", gotID, "last123")
		}
		gotPerPage := r.URL.Query().Get("perPage")
		if gotPerPage != "50" {
			t.Errorf("query param perPage = %q, want %q", gotPerPage, "50")
		}
		gotField := r.URL.Query().Get("field")
		if gotField != "status" {
			t.Errorf("query param field = %q, want %q", gotField, "status")
		}
		gotType := r.URL.Query().Get("type")
		if gotType != "IssueUpdated" {
			t.Errorf("query param type = %q, want %q", gotType, "IssueUpdated")
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[]`)
	})

	ctx := context.Background()
	opts := &ChangelogOptions{ID: "last123", PerPage: 50, Field: "status", Type: "IssueUpdated"}
	_, _, err := client.Issues.GetChangelog(ctx, "QUEUE-1", opts)
	if err != nil {
		t.Fatalf("GetChangelog returned error: %v", err)
	}
}

func TestIssuesService_GetChangelog_Links(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/issues/{key}/changelog", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{
			"id": "1",
			"type": "IssueLinked",
			"links": [{
				"from": null,
				"to": {
					"direction": "inward",
					"object": {
						"self": "https://api.tracker.yandex.net/v3/issues/SIG-2",
						"id": "123",
						"key": "SIG-2",
						"display": "Some issue"
					},
					"type": {
						"self": "https://api.tracker.yandex.net/v3/linktypes/relates",
						"id": "relates",
						"inward": "Linked",
						"outward": "Linked"
					}
				}
			}]
		}]`)
	})

	ctx := context.Background()
	changelog, _, err := client.Issues.GetChangelog(ctx, "QUEUE-1", nil)
	if err != nil {
		t.Fatalf("GetChangelog returned error: %v", err)
	}

	entry := changelog[0]
	if got := *entry.Type; got != "IssueLinked" {
		t.Errorf("Type = %q, want %q", got, "IssueLinked")
	}
	if len(entry.Links) != 1 {
		t.Fatalf("Links length = %d, want 1", len(entry.Links))
	}
	if entry.Links[0].From != nil {
		t.Error("Links[0].From should be nil")
	}
	if got := *entry.Links[0].To.Direction; got != "inward" {
		t.Errorf("Links[0].To.Direction = %q, want %q", got, "inward")
	}
	if got := *entry.Links[0].To.Object.Key; got != "SIG-2" {
		t.Errorf("Links[0].To.Object.Key = %q, want %q", got, "SIG-2")
	}
	if got := string(*entry.Links[0].To.Type.ID); got != "relates" {
		t.Errorf("Links[0].To.Type.ID = %q, want %q", got, "relates")
	}
	if got := *entry.Links[0].To.Type.Inward; got != "Linked" {
		t.Errorf("Links[0].To.Type.Inward = %q, want %q", got, "Linked")
	}
}

func TestIssuesService_GetChangelog_Unlinked(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/issues/{key}/changelog", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{
			"id": "1",
			"type": "IssueUnlinked",
			"links": [{
				"from": {
					"direction": "outward",
					"object": {
						"self": "https://api.tracker.yandex.net/v3/issues/SIG-2",
						"id": "123",
						"key": "SIG-2",
						"display": "Some issue"
					},
					"type": {
						"self": "https://api.tracker.yandex.net/v3/linktypes/relates",
						"id": "relates",
						"inward": "Linked",
						"outward": "Linked"
					}
				},
				"to": null
			}]
		}]`)
	})

	ctx := context.Background()
	changelog, _, err := client.Issues.GetChangelog(ctx, "QUEUE-1", nil)
	if err != nil {
		t.Fatalf("GetChangelog returned error: %v", err)
	}

	entry := changelog[0]
	if got := *entry.Links[0].From.Direction; got != "outward" {
		t.Errorf("Links[0].From.Direction = %q, want %q", got, "outward")
	}
	if entry.Links[0].To != nil {
		t.Error("Links[0].To should be nil")
	}
}

func TestIssuesService_GetChangelog_LinkChanged(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/issues/{key}/changelog", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{
			"id": "1",
			"type": "IssueLinkChanged",
			"links": [{
				"from": {
					"direction": "outward",
					"object": {"key": "SIG-2"},
					"type": {"id": "relates"}
				},
				"to": {
					"direction": "outward",
					"object": {"key": "SIG-2"},
					"type": {"id": "depends"}
				}
			}]
		}]`)
	})

	ctx := context.Background()
	changelog, _, err := client.Issues.GetChangelog(ctx, "QUEUE-1", nil)
	if err != nil {
		t.Fatalf("GetChangelog returned error: %v", err)
	}

	entry := changelog[0]
	if got := string(*entry.Links[0].From.Type.ID); got != "relates" {
		t.Errorf("Links[0].From.Type.ID = %q, want %q", got, "relates")
	}
	if got := string(*entry.Links[0].To.Type.ID); got != "depends" {
		t.Errorf("Links[0].To.Type.ID = %q, want %q", got, "depends")
	}
}

func TestIssuesService_GetChangelog_AttachmentAdded(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/issues/{key}/changelog", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{
			"id": "1",
			"type": "IssueAttachmentAdded",
			"attachments": {
				"added": [{
					"self": "https://api.tracker.yandex.net/v3/issues/SIG-1/attachments/4",
					"id": "4",
					"display": "test.txt"
				}]
			}
		}]`)
	})

	ctx := context.Background()
	changelog, _, err := client.Issues.GetChangelog(ctx, "QUEUE-1", nil)
	if err != nil {
		t.Fatalf("GetChangelog returned error: %v", err)
	}

	entry := changelog[0]
	if entry.Attachments == nil {
		t.Fatal("Attachments is nil")
	}
	if len(entry.Attachments.Added) != 1 {
		t.Fatalf("Attachments.Added length = %d, want 1", len(entry.Attachments.Added))
	}
	if got := string(*entry.Attachments.Added[0].ID); got != "4" {
		t.Errorf("Attachments.Added[0].ID = %q, want %q", got, "4")
	}
	if got := *entry.Attachments.Added[0].Display; got != "test.txt" {
		t.Errorf("Attachments.Added[0].Display = %q, want %q", got, "test.txt")
	}
}

func TestIssuesService_GetChangelog_AttachmentRemoved(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/issues/{key}/changelog", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{
			"id": "1",
			"type": "IssueAttachmentRemoved",
			"attachments": {
				"removed": [{
					"self": "https://api.tracker.yandex.net/v3/issues/SIG-1/attachments/4",
					"id": "4",
					"display": "test.txt"
				}]
			}
		}]`)
	})

	ctx := context.Background()
	changelog, _, err := client.Issues.GetChangelog(ctx, "QUEUE-1", nil)
	if err != nil {
		t.Fatalf("GetChangelog returned error: %v", err)
	}

	entry := changelog[0]
	if entry.Attachments == nil {
		t.Fatal("Attachments is nil")
	}
	if len(entry.Attachments.Removed) != 1 {
		t.Fatalf("Attachments.Removed length = %d, want 1", len(entry.Attachments.Removed))
	}
	if got := *entry.Attachments.Removed[0].Display; got != "test.txt" {
		t.Errorf("Attachments.Removed[0].Display = %q, want %q", got, "test.txt")
	}
}

func TestIssuesService_GetChangelog_Worklog(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/issues/{key}/changelog", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{
			"id": "1",
			"type": "IssueUpdated",
			"worklog": [{
				"record": {
					"self": "https://api.tracker.yandex.net/v3/issues/SIG-1/worklog/7",
					"id": "7",
					"display": ""
				},
				"from": null,
				"to": {
					"duration": "PT30M",
					"start": "2026-04-07T16:52:44.904+0000"
				}
			}]
		}]`)
	})

	ctx := context.Background()
	changelog, _, err := client.Issues.GetChangelog(ctx, "QUEUE-1", nil)
	if err != nil {
		t.Fatalf("GetChangelog returned error: %v", err)
	}

	entry := changelog[0]
	if len(entry.Worklog) != 1 {
		t.Fatalf("Worklog length = %d, want 1", len(entry.Worklog))
	}
	wl := entry.Worklog[0]
	if got := string(*wl.Record.ID); got != "7" {
		t.Errorf("Worklog[0].Record.ID = %q, want %q", got, "7")
	}
	if wl.From != nil {
		t.Error("Worklog[0].From should be nil")
	}
	if wl.To == nil {
		t.Fatal("Worklog[0].To is nil")
	}
	if got := wl.To.Duration.String(); got != "30m0s" {
		t.Errorf("Worklog[0].To.Duration = %q, want %q", got, "30m0s")
	}
}

func TestIssuesService_GetChangelog_WorklogUpdated(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/issues/{key}/changelog", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{
			"id": "1",
			"type": "IssueUpdated",
			"worklog": [{
				"record": {"id": "7"},
				"from": {"duration": "PT30M", "start": "2026-04-07T16:00:00.000+0000"},
				"to": {"duration": "PT1H", "start": "2026-04-07T16:00:00.000+0000"}
			}]
		}]`)
	})

	ctx := context.Background()
	changelog, _, err := client.Issues.GetChangelog(ctx, "QUEUE-1", nil)
	if err != nil {
		t.Fatalf("GetChangelog returned error: %v", err)
	}

	wl := changelog[0].Worklog[0]
	if got := wl.From.Duration.String(); got != "30m0s" {
		t.Errorf("From.Duration = %q, want %q", got, "30m0s")
	}
	if got := wl.To.Duration.String(); got != "1h0m0s" {
		t.Errorf("To.Duration = %q, want %q", got, "1h0m0s")
	}
}

func TestIssuesService_GetChangelog_WorklogRemoved(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/issues/{key}/changelog", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{
			"id": "1",
			"type": "IssueUpdated",
			"worklog": [{
				"record": {"id": "7"},
				"from": {"duration": "PT1H", "start": "2026-04-07T16:00:00.000+0000"},
				"to": null
			}]
		}]`)
	})

	ctx := context.Background()
	changelog, _, err := client.Issues.GetChangelog(ctx, "QUEUE-1", nil)
	if err != nil {
		t.Fatalf("GetChangelog returned error: %v", err)
	}

	wl := changelog[0].Worklog[0]
	if wl.From == nil {
		t.Fatal("From is nil")
	}
	if wl.To != nil {
		t.Error("To should be nil")
	}
}

func TestIssuesService_GetChangelog_RelatedResolution(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/issues/{key}/changelog", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{
			"id": "1",
			"type": "RelatedIssueResolutionChanged",
			"relatedResolutions": [{
				"direction": "outward",
				"issue": {
					"self": "https://api.tracker.yandex.net/v3/issues/SIG-10",
					"id": "123",
					"key": "SIG-10",
					"display": "Dep issue"
				},
				"linkType": {
					"self": "https://api.tracker.yandex.net/v3/linktypes/depends",
					"id": "depends",
					"inward": "Blocker",
					"outward": "Depends on"
				},
				"newResolution": {
					"self": "https://api.tracker.yandex.net/v3/resolutions/1",
					"id": "1",
					"key": "fixed",
					"display": "Fixed"
				}
			}]
		}]`)
	})

	ctx := context.Background()
	changelog, _, err := client.Issues.GetChangelog(ctx, "QUEUE-1", nil)
	if err != nil {
		t.Fatalf("GetChangelog returned error: %v", err)
	}

	entry := changelog[0]
	if len(entry.RelatedResolutions) != 1 {
		t.Fatalf("RelatedResolutions length = %d, want 1", len(entry.RelatedResolutions))
	}
	rr := entry.RelatedResolutions[0]
	if got := *rr.Direction; got != "outward" {
		t.Errorf("Direction = %q, want %q", got, "outward")
	}
	if got := *rr.Issue.Key; got != "SIG-10" {
		t.Errorf("Issue.Key = %q, want %q", got, "SIG-10")
	}
	if got := string(*rr.LinkType.ID); got != "depends" {
		t.Errorf("LinkType.ID = %q, want %q", got, "depends")
	}
	if got := *rr.NewResolution.Key; got != "fixed" {
		t.Errorf("NewResolution.Key = %q, want %q", got, "fixed")
	}
}

func TestIssuesService_GetChangelog_CommentRemoved(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/issues/{key}/changelog", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{
			"id": "1",
			"type": "IssueCommentRemoved",
			"comments": {
				"removed": [{
					"self": "https://api.tracker.yandex.net/v3/issues/SIG-2/comments/7",
					"id": "7",
					"display": "Comment text",
					"objectId": "69d53570cd87d256994e5d00"
				}]
			}
		}]`)
	})

	ctx := context.Background()
	changelog, _, err := client.Issues.GetChangelog(ctx, "QUEUE-1", nil)
	if err != nil {
		t.Fatalf("GetChangelog returned error: %v", err)
	}

	entry := changelog[0]
	if len(entry.Comments.Removed) != 1 {
		t.Fatalf("Comments.Removed length = %d, want 1", len(entry.Comments.Removed))
	}
	cr := entry.Comments.Removed[0]
	if got := string(*cr.ID); got != "7" {
		t.Errorf("ID = %q, want %q", got, "7")
	}
	if got := *cr.Display; got != "Comment text" {
		t.Errorf("Display = %q, want %q", got, "Comment text")
	}
	if got := *cr.ObjectID; got != "69d53570cd87d256994e5d00" {
		t.Errorf("ObjectID = %q, want %q", got, "69d53570cd87d256994e5d00")
	}
}

func TestIssuesService_GetChangelog_CommentUpdated(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/issues/{key}/changelog", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{
			"id": "1",
			"type": "IssueCommentUpdated",
			"comments": {
				"updated": [{
					"comment": {
						"self": "https://api.tracker.yandex.net/v3/issues/SIG-2/comments/7",
						"id": "7",
						"display": "New text",
						"objectId": "abc123"
					},
					"from": "Old text",
					"to": "New text"
				}]
			}
		}]`)
	})

	ctx := context.Background()
	changelog, _, err := client.Issues.GetChangelog(ctx, "QUEUE-1", nil)
	if err != nil {
		t.Fatalf("GetChangelog returned error: %v", err)
	}

	entry := changelog[0]
	if len(entry.Comments.Updated) != 1 {
		t.Fatalf("Comments.Updated length = %d, want 1", len(entry.Comments.Updated))
	}
	cu := entry.Comments.Updated[0]
	if got := string(*cu.Comment.ID); got != "7" {
		t.Errorf("Comment.ID = %q, want %q", got, "7")
	}
	if got := *cu.Comment.ObjectID; got != "abc123" {
		t.Errorf("Comment.ObjectID = %q, want %q", got, "abc123")
	}
	if got, ok := cu.From.(string); !ok || got != "Old text" {
		t.Errorf("From = %v, want %q", cu.From, "Old text")
	}
	if got, ok := cu.To.(string); !ok || got != "New text" {
		t.Errorf("To = %v, want %q", cu.To, "New text")
	}
}

func TestIssuesService_GetChangelog_ReactionAdded(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/issues/{key}/changelog", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{
			"id": "1",
			"type": "IssueCommentReactionAdded",
			"comments": {
				"updated": [{
					"addedReaction": "like",
					"comment": {
						"id": "4",
						"objectId": "abc123",
						"self": "https://api.tracker.yandex.net/v3/issues/SIG-1/comments/4"
					}
				}]
			}
		}]`)
	})

	ctx := context.Background()
	changelog, _, err := client.Issues.GetChangelog(ctx, "QUEUE-1", nil)
	if err != nil {
		t.Fatalf("GetChangelog returned error: %v", err)
	}

	cu := changelog[0].Comments.Updated[0]
	if got := *cu.AddedReaction; got != "like" {
		t.Errorf("AddedReaction = %q, want %q", got, "like")
	}
	if got := string(*cu.Comment.ID); got != "4" {
		t.Errorf("Comment.ID = %q, want %q", got, "4")
	}
}

func TestIssuesService_GetChangelog_ReactionRemoved(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/issues/{key}/changelog", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{
			"id": "1",
			"type": "IssueCommentReactionRemoved",
			"comments": {
				"updated": [{
					"removedReaction": "heart",
					"comment": {
						"id": "5",
						"objectId": "def456",
						"self": "https://api.tracker.yandex.net/v3/issues/SIG-1/comments/5"
					}
				}]
			}
		}]`)
	})

	ctx := context.Background()
	changelog, _, err := client.Issues.GetChangelog(ctx, "QUEUE-1", nil)
	if err != nil {
		t.Fatalf("GetChangelog returned error: %v", err)
	}

	cu := changelog[0].Comments.Updated[0]
	if got := *cu.RemovedReaction; got != "heart" {
		t.Errorf("RemovedReaction = %q, want %q", got, "heart")
	}
}

func TestIssuesService_GetChangelog_Empty(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/issues/{key}/changelog", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[]`)
	})

	ctx := context.Background()
	changelog, _, err := client.Issues.GetChangelog(ctx, "QUEUE-1", nil)
	if err != nil {
		t.Fatalf("GetChangelog returned error: %v", err)
	}
	if len(changelog) != 0 {
		t.Errorf("GetChangelog returned %d entries, want 0", len(changelog))
	}
}
