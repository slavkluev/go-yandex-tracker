package tracker_test

import (
	"context"
	"fmt"
	"log"

	"github.com/slavkluev/go-yandex-tracker/tracker"
)

func ExamplePrioritiesService_List() {
	client := tracker.NewClient(
		tracker.WithOAuthToken("your-oauth-token"),
		tracker.WithOrgID("your-org-id"),
	)

	priorities, _, err := client.Priorities.List(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range priorities {
		fmt.Println(*p.Key, *p.Name)
	}
}
