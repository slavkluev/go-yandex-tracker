package tracker_test

import (
	"context"
	"fmt"
	"log"

	"github.com/slavkluev/go-yandex-tracker/tracker"
)

func ExampleUsersService_Myself() {
	client := tracker.NewClient(
		tracker.WithOAuthToken("your-oauth-token"),
		tracker.WithOrgID("your-org-id"),
	)

	user, _, err := client.Users.Myself(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*user.Display)
}

func ExampleUsersService_Get() {
	client := tracker.NewClient(
		tracker.WithOAuthToken("your-oauth-token"),
		tracker.WithOrgID("your-org-id"),
	)

	user, _, err := client.Users.Get(context.Background(), "user123")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*user.Display)
}
