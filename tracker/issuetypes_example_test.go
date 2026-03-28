package tracker_test

import (
	"context"
	"fmt"
	"log"

	"github.com/slavkluev/go-yandex-tracker/tracker"
)

func ExampleIssueTypesService_List() {
	client := tracker.NewClient(
		tracker.WithOAuthToken("your-oauth-token"),
		tracker.WithOrgID("your-org-id"),
	)

	issueTypes, _, err := client.IssueTypes.List(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, t := range issueTypes {
		fmt.Println(*t.Key, *t.Name)
	}
}
