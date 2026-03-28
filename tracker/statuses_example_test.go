package tracker_test

import (
	"context"
	"fmt"
	"log"

	"github.com/slavkluev/go-yandex-tracker/tracker"
)

func ExampleStatusesService_List() {
	client := tracker.NewClient(
		tracker.WithOAuthToken("your-oauth-token"),
		tracker.WithOrgID("your-org-id"),
	)

	statuses, _, err := client.Statuses.List(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range statuses {
		fmt.Println(*s.Key, *s.Name)
	}
}
