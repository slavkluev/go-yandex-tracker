package tracker_test

import (
	"context"
	"fmt"
	"log"

	"github.com/slavkluev/go-yandex-tracker/tracker"
)

func ExampleEntitiesService_Create() {
	client := tracker.NewClient(
		tracker.WithOAuthToken("your-oauth-token"),
		tracker.WithOrgID("your-org-id"),
	)

	entity, _, err := client.Entities.Create(context.Background(), tracker.EntityTypeProject, &tracker.EntityCreateRequest{
		Fields: &tracker.EntityCreateFields{
			Summary:    tracker.Ptr("New Project"),
			TeamAccess: tracker.Ptr(true),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*entity.ID)
}

func ExampleEntitiesService_Search() {
	client := tracker.NewClient(
		tracker.WithOAuthToken("your-oauth-token"),
		tracker.WithOrgID("your-org-id"),
	)

	result, _, err := client.Entities.Search(context.Background(), tracker.EntityTypeProject, &tracker.EntitySearchRequest{
		Input: tracker.Ptr("project"),
	}, nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range result.Values {
		fmt.Println(*e.ID)
	}
}
