package tracker_test

import (
	"context"
	"fmt"
	"log"

	"github.com/slavkluev/go-yandex-tracker/tracker"
)

func ExampleFiltersService_List() {
	client := tracker.NewClient(
		tracker.WithOAuthToken("your-oauth-token"),
		tracker.WithOrgID("your-org-id"),
	)

	filters, _, err := client.Filters.List(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range filters {
		fmt.Println(*f.Name)
	}
}

func ExampleFiltersService_Create() {
	client := tracker.NewClient(
		tracker.WithOAuthToken("your-oauth-token"),
		tracker.WithOrgID("your-org-id"),
	)

	filter, _, err := client.Filters.Create(context.Background(), &tracker.FilterRequest{
		Name:   tracker.Ptr("My Filter"),
		Filter: map[string]any{"queue": "QUEUE", "assignee": "me"},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*filter.Name)
}
