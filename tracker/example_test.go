package tracker_test

import (
	"fmt"

	"github.com/slavkluev/go-yandex-tracker/tracker"
)

func ExampleNewClient_oAuth() {
	client := tracker.NewClient(
		tracker.WithOAuthToken("your-oauth-token"),
		tracker.WithOrgID("your-org-id"),
	)

	// Use client to interact with the Yandex Tracker API.
	fmt.Println(client)
}

func ExampleNewClient_iAMToken() {
	client := tracker.NewClient(
		tracker.WithIAMToken("your-iam-token"),
		tracker.WithCloudOrgID("your-cloud-org-id"),
	)

	// Use client to interact with the Yandex Tracker API.
	fmt.Println(client)
}
