package tracker_test

import (
	"context"
	"fmt"
	"log"

	"github.com/slavkluev/go-yandex-tracker/tracker"
)

func ExampleExternalLinksService_ListLinks() {
	client := tracker.NewClient(
		tracker.WithOAuthToken("your-oauth-token"),
		tracker.WithOrgID("your-org-id"),
	)

	links, _, err := client.ExternalLinks.ListLinks(context.Background(), "QUEUE-1")
	if err != nil {
		log.Fatal(err)
	}

	for _, link := range links {
		fmt.Println(*link.ID)
	}
}
