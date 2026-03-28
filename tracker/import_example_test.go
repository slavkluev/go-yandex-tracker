package tracker_test

import (
	"context"
	"fmt"
	"log"

	"github.com/slavkluev/go-yandex-tracker/tracker"
)

func ExampleImportService_ImportIssue() {
	client := tracker.NewClient(
		tracker.WithOAuthToken("your-oauth-token"),
		tracker.WithOrgID("your-org-id"),
	)

	issue, _, err := client.Import.ImportIssue(context.Background(), &tracker.ImportIssueRequest{
		Queue:     tracker.Ptr("QUEUE"),
		Summary:   tracker.Ptr("Imported issue"),
		CreatedAt: tracker.Ptr("2024-01-15T10:00:00.000+0000"),
		CreatedBy: tracker.Ptr("admin"),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*issue.Key)
}
