package tracker_test

import (
	"context"
	"fmt"
	"log"

	"github.com/slavkluev/go-yandex-tracker/tracker"
)

func ExampleResolutionsService_List() {
	client := tracker.NewClient(
		tracker.WithOAuthToken("your-oauth-token"),
		tracker.WithOrgID("your-org-id"),
	)

	resolutions, _, err := client.Resolutions.List(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range resolutions {
		fmt.Println(*r.Key, *r.Name)
	}
}
