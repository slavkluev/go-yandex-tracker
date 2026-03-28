package tracker_test

import (
	"context"
	"fmt"
	"log"

	"github.com/slavkluev/go-yandex-tracker/tracker"
)

func ExampleBulkChangeService_Move() {
	client := tracker.NewClient(
		tracker.WithOAuthToken("your-oauth-token"),
		tracker.WithOrgID("your-org-id"),
	)

	// Initiate bulk move -- returns immediately with operation ID.
	op, _, err := client.BulkChange.Move(context.Background(), &tracker.BulkMoveRequest{
		Queue:  tracker.Ptr("TARGET"),
		Issues: []string{"SOURCE-1", "SOURCE-2", "SOURCE-3"},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Poll for completion using the operation ID.
	status, _, err := client.BulkChange.GetStatus(context.Background(), *op.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*status.Status)
}
