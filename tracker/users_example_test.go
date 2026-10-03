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

func ExampleUser_DisplayOr() {
	lead := &tracker.User{ID: tracker.Ptr(tracker.FlexString("1234567890"))}
	var nobody *tracker.User

	fmt.Println(lead.DisplayOr("-"))
	fmt.Println(nobody.DisplayOr("-"))
	// Output:
	// 1234567890
	// -
}
