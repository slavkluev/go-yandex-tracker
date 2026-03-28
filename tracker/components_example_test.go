package tracker_test

import (
	"context"
	"fmt"
	"log"

	"github.com/slavkluev/go-yandex-tracker/tracker"
)

func ExampleComponentsService_List() {
	client := tracker.NewClient(
		tracker.WithOAuthToken("your-oauth-token"),
		tracker.WithOrgID("your-org-id"),
	)

	components, _, err := client.Components.List(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range components {
		fmt.Println(*c.Name)
	}
}

func ExampleComponentsService_Create() {
	client := tracker.NewClient(
		tracker.WithOAuthToken("your-oauth-token"),
		tracker.WithOrgID("your-org-id"),
	)

	component, _, err := client.Components.Create(context.Background(), &tracker.ComponentRequest{
		Name:  tracker.Ptr("Backend"),
		Queue: tracker.Ptr("QUEUE"),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*component.Name)
}
