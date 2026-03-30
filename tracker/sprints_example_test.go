package tracker_test

import (
	"context"
	"fmt"
	"log"

	"github.com/slavkluev/go-yandex-tracker/tracker"
)

func ExampleSprintsService_List() {
	client := tracker.NewClient(
		tracker.WithOAuthToken("your-oauth-token"),
		tracker.WithOrgID("your-org-id"),
	)

	sprints, _, err := client.Sprints.List(context.Background(), "1")
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range sprints {
		fmt.Println(*s.Name)
	}
}

func ExampleSprintsService_Create() {
	client := tracker.NewClient(
		tracker.WithOAuthToken("your-oauth-token"),
		tracker.WithOrgID("your-org-id"),
	)

	sprint, _, err := client.Sprints.Create(context.Background(), &tracker.SprintCreateRequest{
		Name:  tracker.Ptr("Sprint 1"),
		Board: &tracker.BoardRef{ID: tracker.Ptr(tracker.FlexString("1"))},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*sprint.Name)
}
