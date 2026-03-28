package tracker_test

import (
	"context"
	"fmt"
	"log"

	"github.com/slavkluev/go-yandex-tracker/tracker"
)

func ExampleDashboardsService_Create() {
	client := tracker.NewClient(
		tracker.WithOAuthToken("your-oauth-token"),
		tracker.WithOrgID("your-org-id"),
	)

	dashboard, _, err := client.Dashboards.Create(context.Background(), &tracker.DashboardCreateRequest{
		Name: tracker.Ptr("My Dashboard"),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*dashboard.Name)
}
