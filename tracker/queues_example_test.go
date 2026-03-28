package tracker_test

import (
	"context"
	"fmt"
	"log"

	"github.com/slavkluev/go-yandex-tracker/tracker"
)

func ExampleQueuesService_Get() {
	client := tracker.NewClient(
		tracker.WithOAuthToken("your-oauth-token"),
		tracker.WithOrgID("your-org-id"),
	)

	queue, _, err := client.Queues.Get(context.Background(), "QUEUE", nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*queue.Key, *queue.Name)
}

func ExampleQueuesService_Create() {
	client := tracker.NewClient(
		tracker.WithOAuthToken("your-oauth-token"),
		tracker.WithOrgID("your-org-id"),
	)

	queue, _, err := client.Queues.Create(context.Background(), &tracker.QueueCreateRequest{
		Key:  tracker.Ptr("NEWQUEUE"),
		Name: tracker.Ptr("New Queue"),
		Lead: tracker.Ptr("admin"),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*queue.Key)
}
