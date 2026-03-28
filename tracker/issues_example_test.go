package tracker_test

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/slavkluev/go-yandex-tracker/tracker"
)

func ExampleIssuesService_Get() {
	client := tracker.NewClient(
		tracker.WithOAuthToken("your-oauth-token"),
		tracker.WithOrgID("your-org-id"),
	)

	issue, _, err := client.Issues.Get(context.Background(), "QUEUE-1", nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*issue.Key, *issue.Summary)
}

func ExampleIssuesService_Create() {
	client := tracker.NewClient(
		tracker.WithOAuthToken("your-oauth-token"),
		tracker.WithOrgID("your-org-id"),
	)

	issue, _, err := client.Issues.Create(context.Background(), &tracker.IssueRequest{
		Summary: tracker.Ptr("Bug: login fails"),
		Queue:   tracker.Ptr("QUEUE"),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*issue.Key)
}

func ExampleIssuesService_Search() {
	client := tracker.NewClient(
		tracker.WithOAuthToken("your-oauth-token"),
		tracker.WithOrgID("your-org-id"),
	)

	issues, _, err := client.Issues.Search(context.Background(), &tracker.IssueSearchRequest{
		Filter: map[string]any{
			"queue":    "QUEUE",
			"assignee": "me",
		},
	}, nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, issue := range issues {
		fmt.Println(*issue.Key)
	}
}

func ExampleIssuesService_Search_pagination() {
	client := tracker.NewClient(
		tracker.WithOAuthToken("your-oauth-token"),
		tracker.WithOrgID("your-org-id"),
	)

	ctx := context.Background()
	req := &tracker.IssueSearchRequest{
		Filter: map[string]any{"queue": "QUEUE"},
	}
	opts := &tracker.IssueSearchOptions{
		ListOptions: tracker.ListOptions{Page: 1, PerPage: 50},
	}

	for {
		issues, _, err := client.Issues.Search(ctx, req, opts)
		if err != nil {
			log.Fatal(err)
		}

		for _, issue := range issues {
			fmt.Println(*issue.Key)
		}

		if len(issues) < opts.PerPage {
			break
		}
		opts.Page++
	}
}

func ExampleIssuesService_ListComments() {
	client := tracker.NewClient(
		tracker.WithOAuthToken("your-oauth-token"),
		tracker.WithOrgID("your-org-id"),
	)

	comments, _, err := client.Issues.ListComments(context.Background(), "QUEUE-1", nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, comment := range comments {
		fmt.Println(*comment.Text)
	}
}

func ExampleIssuesService_UploadAttachment() {
	client := tracker.NewClient(
		tracker.WithOAuthToken("your-oauth-token"),
		tracker.WithOrgID("your-org-id"),
	)

	reader := strings.NewReader("file content")

	attachment, _, err := client.Issues.UploadAttachment(context.Background(), "QUEUE-1", "report.txt", reader)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*attachment.Name)
}
